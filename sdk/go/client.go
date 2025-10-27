package teamsnotify

// Client provides access to all Teams Notification API services
type Client struct {
    baseURL    string
    apiKey     string
    httpClient *HTTPClient
}

// NewClient creates a new Client
func NewClient(baseURL string) *Client {
    return &Client{
        baseURL:    baseURL,
        httpClient: NewHTTPClient(),
    }
}

// NewClientWithAuth creates a new Client with authentication
func NewClientWithAuth(baseURL, apiKey string) *Client {
    return &Client{
        baseURL:    baseURL,
        apiKey:     apiKey,
        httpClient: NewHTTPClient().WithAuth(apiKey),
    }
}

// Notifications returns notifications service
func (c *Client) Notifications() *NotificationService {
    return NewNotificationService(c)
}

// Users returns users service
func (c *Client) Users() *UserService {
    return NewUserService(c)
}

// Projects returns projects service
func (c *Client) Projects() *ProjectService {
    return NewProjectService(c)
}

// Destinations returns destinations service
func (c *Client) Destinations() *DestinationService {
    return NewDestinationService(c)
}

// Bots returns bots service
func (c *Client) Bots() *BotService {
    return NewBotService(c)
}
