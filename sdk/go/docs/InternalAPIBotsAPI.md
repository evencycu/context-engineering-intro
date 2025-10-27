# \InternalAPIBotsAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InternalV1BotsPlatformGet**](InternalAPIBotsAPI.md#InternalV1BotsPlatformGet) | **Get** /internal/v1/bots/platform | List Teams bots
[**InternalV1BotsPlatformIdCapabilitiesPatch**](InternalAPIBotsAPI.md#InternalV1BotsPlatformIdCapabilitiesPatch) | **Patch** /internal/v1/bots/platform/{id}/capabilities | Update bot capabilities
[**InternalV1BotsPlatformIdDelete**](InternalAPIBotsAPI.md#InternalV1BotsPlatformIdDelete) | **Delete** /internal/v1/bots/platform/{id} | Delete Teams bot
[**InternalV1BotsPlatformIdGet**](InternalAPIBotsAPI.md#InternalV1BotsPlatformIdGet) | **Get** /internal/v1/bots/platform/{id} | Get Teams bot
[**InternalV1BotsPlatformIdPut**](InternalAPIBotsAPI.md#InternalV1BotsPlatformIdPut) | **Put** /internal/v1/bots/platform/{id} | Update Teams bot
[**InternalV1BotsPlatformIdStatusPatch**](InternalAPIBotsAPI.md#InternalV1BotsPlatformIdStatusPatch) | **Patch** /internal/v1/bots/platform/{id}/status | Update bot status
[**InternalV1BotsPlatformIdTestPost**](InternalAPIBotsAPI.md#InternalV1BotsPlatformIdTestPost) | **Post** /internal/v1/bots/platform/{id}/test | Test bot connection
[**InternalV1BotsPlatformPost**](InternalAPIBotsAPI.md#InternalV1BotsPlatformPost) | **Post** /internal/v1/bots/platform | Create Teams bot
[**InternalV1BotsStatusStatusGet**](InternalAPIBotsAPI.md#InternalV1BotsStatusStatusGet) | **Get** /internal/v1/bots/status/{status} | Get bots by status



## InternalV1BotsPlatformGet

> BotListResponse InternalV1BotsPlatformGet(ctx).Limit(limit).Offset(offset).Search(search).Execute()

List Teams bots



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	limit := int32(56) // int32 | Number of items to return (optional) (default to 10)
	offset := int32(56) // int32 | Number of items to skip (optional) (default to 0)
	search := "search_example" // string | Search term (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBotsAPI.InternalV1BotsPlatformGet(context.Background()).Limit(limit).Offset(offset).Search(search).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBotsAPI.InternalV1BotsPlatformGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BotsPlatformGet`: BotListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBotsAPI.InternalV1BotsPlatformGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BotsPlatformGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | Number of items to return | [default to 10]
 **offset** | **int32** | Number of items to skip | [default to 0]
 **search** | **string** | Search term | 

### Return type

[**BotListResponse**](BotListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BotsPlatformIdCapabilitiesPatch

> BotResponse InternalV1BotsPlatformIdCapabilitiesPatch(ctx, id).BotCapabilitiesUpdateRequest(botCapabilitiesUpdateRequest).Execute()

Update bot capabilities



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	botCapabilitiesUpdateRequest := *openapiclient.NewBotCapabilitiesUpdateRequest([]string{"Capabilities_example"}) // BotCapabilitiesUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBotsAPI.InternalV1BotsPlatformIdCapabilitiesPatch(context.Background(), id).BotCapabilitiesUpdateRequest(botCapabilitiesUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBotsAPI.InternalV1BotsPlatformIdCapabilitiesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BotsPlatformIdCapabilitiesPatch`: BotResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBotsAPI.InternalV1BotsPlatformIdCapabilitiesPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BotsPlatformIdCapabilitiesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **botCapabilitiesUpdateRequest** | [**BotCapabilitiesUpdateRequest**](BotCapabilitiesUpdateRequest.md) |  | 

### Return type

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BotsPlatformIdDelete

> InternalV1BotsPlatformIdDelete(ctx, id).Execute()

Delete Teams bot



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.InternalAPIBotsAPI.InternalV1BotsPlatformIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBotsAPI.InternalV1BotsPlatformIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BotsPlatformIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BotsPlatformIdGet

> BotResponse InternalV1BotsPlatformIdGet(ctx, id).Execute()

Get Teams bot



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBotsAPI.InternalV1BotsPlatformIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBotsAPI.InternalV1BotsPlatformIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BotsPlatformIdGet`: BotResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBotsAPI.InternalV1BotsPlatformIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BotsPlatformIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BotsPlatformIdPut

> BotResponse InternalV1BotsPlatformIdPut(ctx, id).BotUpdateRequest(botUpdateRequest).Execute()

Update Teams bot



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	botUpdateRequest := *openapiclient.NewBotUpdateRequest() // BotUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBotsAPI.InternalV1BotsPlatformIdPut(context.Background(), id).BotUpdateRequest(botUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBotsAPI.InternalV1BotsPlatformIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BotsPlatformIdPut`: BotResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBotsAPI.InternalV1BotsPlatformIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BotsPlatformIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **botUpdateRequest** | [**BotUpdateRequest**](BotUpdateRequest.md) |  | 

### Return type

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BotsPlatformIdStatusPatch

> BotResponse InternalV1BotsPlatformIdStatusPatch(ctx, id).BotStatusUpdateRequest(botStatusUpdateRequest).Execute()

Update bot status



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	botStatusUpdateRequest := *openapiclient.NewBotStatusUpdateRequest("Status_example") // BotStatusUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBotsAPI.InternalV1BotsPlatformIdStatusPatch(context.Background(), id).BotStatusUpdateRequest(botStatusUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBotsAPI.InternalV1BotsPlatformIdStatusPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BotsPlatformIdStatusPatch`: BotResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBotsAPI.InternalV1BotsPlatformIdStatusPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BotsPlatformIdStatusPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **botStatusUpdateRequest** | [**BotStatusUpdateRequest**](BotStatusUpdateRequest.md) |  | 

### Return type

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BotsPlatformIdTestPost

> BotTestResponse InternalV1BotsPlatformIdTestPost(ctx, id).Execute()

Test bot connection



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBotsAPI.InternalV1BotsPlatformIdTestPost(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBotsAPI.InternalV1BotsPlatformIdTestPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BotsPlatformIdTestPost`: BotTestResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBotsAPI.InternalV1BotsPlatformIdTestPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BotsPlatformIdTestPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BotTestResponse**](BotTestResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BotsPlatformPost

> BotResponse InternalV1BotsPlatformPost(ctx).BotCreateRequest(botCreateRequest).Execute()

Create Teams bot



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	botCreateRequest := *openapiclient.NewBotCreateRequest("Name_example", "AppId_example", "AppPassword_example", []string{"Capabilities_example"}) // BotCreateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBotsAPI.InternalV1BotsPlatformPost(context.Background()).BotCreateRequest(botCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBotsAPI.InternalV1BotsPlatformPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BotsPlatformPost`: BotResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBotsAPI.InternalV1BotsPlatformPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BotsPlatformPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **botCreateRequest** | [**BotCreateRequest**](BotCreateRequest.md) |  | 

### Return type

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BotsStatusStatusGet

> BotListResponse InternalV1BotsStatusStatusGet(ctx, status).Execute()

Get bots by status



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	status := "status_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBotsAPI.InternalV1BotsStatusStatusGet(context.Background(), status).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBotsAPI.InternalV1BotsStatusStatusGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BotsStatusStatusGet`: BotListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBotsAPI.InternalV1BotsStatusStatusGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**status** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BotsStatusStatusGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BotListResponse**](BotListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

