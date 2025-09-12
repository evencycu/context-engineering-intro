package main

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ====== 基本設定 ======
var (
	appID     = "844146d7-4ac9-4e4d-a463-d6e027714e81"     // 你的 Bot Microsoft App ID
	appSecret = "HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy" // 你的 Bot client secret
)

// ====== Activity 結構（節錄常用欄位）======
type ChannelAccount struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ConversationRef struct {
	ID string `json:"id,omitempty"`
}

type Activity struct {
	Type         string          `json:"type"`
	ID           string          `json:"id,omitempty"`
	Timestamp    *time.Time      `json:"timestamp,omitempty"`
	ServiceURL   string          `json:"serviceUrl"`
	From         ChannelAccount  `json:"from"`
	Recipient    ChannelAccount  `json:"recipient"`
	Conversation ConversationRef `json:"conversation"`
	ReplyToID    string          `json:"replyToId,omitempty"`
	Text         string          `json:"text,omitempty"`
	ChannelID    string          `json:"channelId,omitempty"`
	Locale       string          `json:"locale,omitempty"`
}

// ====== OIDC / JWKS 型別 ======
type openIDConfig struct {
	Issuer  string `json:"issuer"`
	JwksURI string `json:"jwks_uri"`
}

type jwkSet struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kty string   `json:"kty"`
	E   string   `json:"e"`
	N   string   `json:"n"`
	Alg string   `json:"alg"`
	Kid string   `json:"kid"`
	X5t string   `json:"x5t,omitempty"`
	X5c []string `json:"x5c,omitempty"` // 也可能提供 cert 鏈
}

// ====== 入口 ======
func main() {
	if appID == "" || appSecret == "" {
		log.Fatal("Please set MICROSOFT_APP_ID and MICROSOFT_APP_PASSWORD")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/messages", handleMessages)
	mux.HandleFunc("/api/proactive", handleProactive)

	addr := ":3978"
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// ====== Handler：驗證簽章 + 解析 + 回覆 ======
func handleMessages(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %+v \n", r)

	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}

	// 讀 body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body error", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 解析 activity（先做，等下驗 serviceUrl）
	var act Activity
	log.Printf("body: %s \n", string(body))
	if err := json.Unmarshal(body, &act); err != nil {
		log.Printf("json.Unmarshal error: %v", err)
		http.Error(w, "invalid activity", http.StatusBadRequest)
		return
	}
	log.Printf("act: %+v \n", act)

	if act.ID == "CopilotAgent" {
		log.Printf("from CopilotAgent: %+v \n", act)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"reply\":\"everyday is a good day\"}"))
		return
	}
	log.Printf("convRef serviceUrl=%s conversationId=%s userId=%s botId=%s",
		act.ServiceURL, act.Conversation.ID, act.From.ID, act.Recipient.ID)
	// 允許 Emulator 無簽章（開發模式可選）。Teams/正式一定要驗。
	authz := r.Header.Get("Authorization")
	log.Printf("Authorization: %s \n", authz)
	if authz == "" {
		// 開發時：允許跳過；正式請直接拒絕。
		log.Println("no Authorization header (emulator?) - skipping token validation in dev mode")
	} else {
		if err := validateBFToken(authz, act.ServiceURL); err != nil {
			log.Printf("validateBFToken error: %v", err)
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
	}

	// // 簡單 echo：把使用者訊息回回去
	// replyText := fmt.Sprintf("You said: %s", act.Text)
	// log.Printf("ReplyText: %s \n", replyText)
	// if err := sendReply(r.Context(), &act, replyText); err != nil {
	// 	log.Printf("sendReply error: %v", err)
	// 	// 不阻擋 200，避免 Teams 重送（你也可回 502 讓它重試）
	// }
	// log.Printf("ReplyText: end")
	singleConversation := "a:1TRu8CjpyEoKXjA2TUhxbLN-5TvJKwzuzwcyTvxkOmBC7A3Ty2GbGiO14h_ZpT4OUaExpjc1xnwNEJkJ9UsSD3pIrLHIZP9LB1KMTuvbfrQtW829RkSazvaIztlqg8P3L"
	// groupConversation := "19:7f469298f07f4abca81d21ad980095fb@thread.v2"
	groupConversation := "19:0a608973d8484d30981e68f1ed45a5c8@thread.v2" // test group
	conversationID := ""
	userID := ""
	if act.Conversation.ID == singleConversation {
		conversationID = groupConversation
	} else {
		conversationID = singleConversation
		userID = "29:1FrpqFNZ44LuLurxUt5tpSed4vXWXTdi428VAbwv9FSK-cI4UWRtIKdwEwFKCehU48w1hrQvIc5AXAoUtNKKBYg"
	}
	botID := "28:844146d7-4ac9-4e4d-a463-d6e027714e81"
	if conversationID == singleConversation {
		for i := 0; i < 105; i++ {
			if err := proactiveMessage(r.Context(), act.ServiceURL, conversationID, userID, botID, fmt.Sprintf("test %d", i)); err != nil {
				log.Printf("proactiveMessage error: %v", err)
				//return
			}

		}
	} else {
		if err := proactiveMessage(r.Context(), act.ServiceURL, conversationID, userID, botID, act.Text); err != nil {
			log.Printf("proactiveMessage error: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func handleProactive(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %+v \n", r)

	body, err := io.ReadAll(r.Body)
	log.Printf("body: %s \n", string(body))
	if err != nil {
		http.Error(w, "read body error", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	serviceUrl := "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/"
	conversationID := "a:1TRu8CjpyEoKXjA2TUhxbLN-5TvJKwzuzwcyTvxkOmBC7A3Ty2GbGiO14h_ZpT4OUaExpjc1xnwNEJkJ9UsSD3pIrLHIZP9LB1KMTuvbfrQtW829RkSazvaIztlqg8P3L"
	userID := "29:1FrpqFNZ44LuLurxUt5tpSed4vXWXTdi428VAbwv9FSK-cI4UWRtIKdwEwFKCehU48w1hrQvIc5AXAoUtNKKBYg "
	botID := "28:844146d7-4ac9-4e4d-a463-d6e027714e81"

	for i := 0; i < 1000; i++ {
		if err := proactiveMessage(r.Context(), serviceUrl, conversationID, userID, botID, fmt.Sprintf("test %d", i)); err != nil {
			log.Printf("proactiveMessage error: %v", err)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func proactiveMessage(ctx context.Context, serviceURL, conversationID, userID, botID, text string) error {
	// 1. 取 Connector access token（跟前面 echo 版本相同）
	accessToken, err := getConnectorToken(ctx)
	if err != nil {
		return err
	}

	// 2. 組新的 Activity
	act := Activity{
		Type:         "message",
		From:         ChannelAccount{ID: botID},
		Recipient:    ChannelAccount{ID: userID},
		Conversation: ConversationRef{ID: conversationID},
		Text:         text,
	}

	b, _ := json.Marshal(act)
	url := fmt.Sprintf("%s/v3/conversations/%s/activities", strings.TrimRight(serviceURL, "/"), conversationID)

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("send proactive failed: %s", string(body))
	}
	return nil
}

// ====== 驗證 Bot Framework/Teams JWT 簽章 ======
func validateBFToken(authorization string, serviceURL string) error {
	const (
		// 常見的 OpenID metadata（全球雲）
		openID = "https://login.botframework.com/v1/.well-known/openidconfiguration"
		// 另一個常見 issuer（Teams/全球）：https://api.botframework.com
		expectedIssuer1 = "https://api.botframework.com"
		expectedIssuer2 = "https://login.botframework.com"
	)

	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return errors.New("invalid authorization header")
	}
	raw := parts[1]

	// 下載 OpenID config & JWKS（建議加快取與快取）
	oidc, err := fetchOpenID(openID)
	if err != nil {
		return fmt.Errorf("fetch openid: %w", err)
	}
	keys, err := fetchJWKS(oidc.JwksURI)
	if err != nil {
		return fmt.Errorf("fetch jwks: %w", err)
	}

	// 自訂 keyfunc：根據 kid 尋找對應 RSA 公鑰
	keyFunc := func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		for _, k := range keys.Keys {
			if k.Kid == kid && k.Kty == "RSA" {
				// 解析 modulus/exponent
				pub, perr := jwkToRSAPublicKey(k)
				return pub, perr
			}
		}
		return nil, errors.New("kid not found")
	}

	// 解析 + 基本驗證
	token, err := jwt.Parse(raw, keyFunc,
		jwt.WithAudience(appID),
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithLeeway(2*time.Minute),
	)
	if err != nil || !token.Valid {
		return fmt.Errorf("jwt parse/invalid: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("claims cast failed")
	}

	// iss 檢查
	iss, _ := claims["iss"].(string)
	if iss != expectedIssuer1 && iss != expectedIssuer2 {
		return fmt.Errorf("unexpected iss: %s", iss)
	}

	// serviceUrl 檢查（避免 token 重放到別的 serviceUrl）
	if su, _ := claims["serviceurl"].(string); su != "" {
		if !strings.EqualFold(su, serviceURL) {
			return fmt.Errorf("serviceurl mismatch: token=%s, activity=%s", su, serviceURL)
		}
	}

	return nil
}

// ====== 呼叫 Connector v3 發回覆 ======
func sendReply(ctx context.Context, act *Activity, text string) error {
	// 1) 取 Connector access token
	accessToken, err := getConnectorToken(ctx)
	if err != nil {
		return fmt.Errorf("get connector token: %w", err)
	}

	// 2) 組 reply activity
	reply := Activity{
		Type:         "message",
		From:         act.Recipient, // 由 bot 發出（通常就是 activity.recipient）
		Recipient:    act.From,      // 回給使用者
		Conversation: act.Conversation,
		ReplyToID:    act.ID,
		Text:         text,
	}

	b, _ := json.Marshal(reply)
	url := fmt.Sprintf("%s/v3/conversations/%s/activities", strings.TrimRight(act.ServiceURL, "/"), act.Conversation.ID)

	req, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		rb, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("connector post %s: %s", resp.Status, string(rb))
	}
	return nil
}

// ====== 取得 Connector Access Token（Client Credentials）======
func getConnectorToken(ctx context.Context) (string, error) {
	// Bot Framework 的特殊租戶：botframework.com
	tenantID := "051cece0-e4dc-4aed-b471-bf29824e1ee6"
	if tenantID == "" {
		return "", errors.New("MICROSOFT_APP_TENANT_ID is empty")
	}
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", appID)
	data.Set("client_secret", appSecret)
	data.Set("scope", "https://api.botframework.com/.default")

	req, _ := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token failed: %s", string(b))
	}
	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.AccessToken, nil
}

// ====== OpenID / JWKS 下載（建議加上快取與容錯）======
func fetchOpenID(url string) (*openIDConfig, error) {
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var cfg openIDConfig
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func fetchJWKS(url string) (*jwkSet, error) {
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var set jwkSet
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return nil, err
	}
	return &set, nil
}

// ====== JWK 轉 RSA Public Key & 小工具 ======
func jwkToRSAPublicKey(k jwk) (*rsa.PublicKey, error) {
	// 這裡為了簡潔，直接用 x5c（若提供）轉公鑰；或自行 decode N/E
	// 生產版本建議同時支援 N/E 與 x5c。
	if len(k.X5c) > 0 {
		return certToRSAPublicKey(k.X5c[0])
	}
	return neToRSAPublicKey(k.N, k.E)
}

func certToRSAPublicKey(b64cert string) (*rsa.PublicKey, error) {
	// x5c 是 base64(der) 的憑證，不含 PEM header/footer
	der, err := base64.StdEncoding.DecodeString(b64cert)
	if err != nil {
		return nil, fmt.Errorf("x5c base64 decode: %w", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, fmt.Errorf("x5c parse cert: %w", err)
	}
	pk, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("x5c is not RSA public key")
	}
	return pk, nil
}
func neToRSAPublicKey(nB64url, eB64url string) (*rsa.PublicKey, error) {
	// JWK 的 N/E 採 base64url (無 '=' padding)
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64url)
	if err != nil {
		return nil, fmt.Errorf("decode N: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eB64url)
	if err != nil {
		return nil, fmt.Errorf("decode E: %w", err)
	}
	// E 是 big-endian 的無號整數，通常是 65537 (0x010001)
	eInt := 0
	for _, b := range eBytes {
		eInt = (eInt << 8) | int(b)
	}
	if eInt <= 0 {
		return nil, errors.New("invalid exponent")
	}
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}, nil
}

func urlEncode(s string) string {
	r := strings.NewReplacer(" ", "%20", ":", "%3A", "/", "%2F")
	return r.Replace(s)
}
