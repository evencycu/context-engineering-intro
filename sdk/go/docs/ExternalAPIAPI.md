# \ExternalAPIAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1DestinationsNotifyKeyGet**](ExternalAPIAPI.md#ApiV1DestinationsNotifyKeyGet) | **Get** /api/v1/destinations/{notifyKey} | Get project destinations (External API)
[**ApiV1MessagesPost**](ExternalAPIAPI.md#ApiV1MessagesPost) | **Post** /api/v1/messages | Handle Bot Framework messages
[**ApiV1MessagesTestPost**](ExternalAPIAPI.md#ApiV1MessagesTestPost) | **Post** /api/v1/messages/test | Test proactive message
[**ApiV1NotifyPost**](ExternalAPIAPI.md#ApiV1NotifyPost) | **Post** /api/v1/notify | Send notification (External API)
[**ApiV1ProvisionNotifyKeyDisablePost**](ExternalAPIAPI.md#ApiV1ProvisionNotifyKeyDisablePost) | **Post** /api/v1/provision/{notifyKey}/disable | Disable project provision
[**ApiV1ProvisionNotifyKeyEnablePost**](ExternalAPIAPI.md#ApiV1ProvisionNotifyKeyEnablePost) | **Post** /api/v1/provision/{notifyKey}/enable | Enable project provision
[**ApiV1ProvisionNotifyKeyGet**](ExternalAPIAPI.md#ApiV1ProvisionNotifyKeyGet) | **Get** /api/v1/provision/{notifyKey} | Get project provision
[**ApiV1ProvisionNotifyKeyPut**](ExternalAPIAPI.md#ApiV1ProvisionNotifyKeyPut) | **Put** /api/v1/provision/{notifyKey} | Update project provision
[**ApiV1ProvisionPost**](ExternalAPIAPI.md#ApiV1ProvisionPost) | **Post** /api/v1/provision | Create project provision



## ApiV1DestinationsNotifyKeyGet

> ProjectDestinationsResponse ApiV1DestinationsNotifyKeyGet(ctx, notifyKey).Execute()

Get project destinations (External API)



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
	notifyKey := "notifyKey_example" // string | Project notify key

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExternalAPIAPI.ApiV1DestinationsNotifyKeyGet(context.Background(), notifyKey).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExternalAPIAPI.ApiV1DestinationsNotifyKeyGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1DestinationsNotifyKeyGet`: ProjectDestinationsResponse
	fmt.Fprintf(os.Stdout, "Response from `ExternalAPIAPI.ApiV1DestinationsNotifyKeyGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**notifyKey** | **string** | Project notify key | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1DestinationsNotifyKeyGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ProjectDestinationsResponse**](ProjectDestinationsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1MessagesPost

> MessageResponse ApiV1MessagesPost(ctx).BotFrameworkActivity(botFrameworkActivity).Execute()

Handle Bot Framework messages



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
	botFrameworkActivity := *openapiclient.NewBotFrameworkActivity() // BotFrameworkActivity | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExternalAPIAPI.ApiV1MessagesPost(context.Background()).BotFrameworkActivity(botFrameworkActivity).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExternalAPIAPI.ApiV1MessagesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1MessagesPost`: MessageResponse
	fmt.Fprintf(os.Stdout, "Response from `ExternalAPIAPI.ApiV1MessagesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1MessagesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **botFrameworkActivity** | [**BotFrameworkActivity**](BotFrameworkActivity.md) |  | 

### Return type

[**MessageResponse**](MessageResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1MessagesTestPost

> MessageResponse ApiV1MessagesTestPost(ctx).ProactiveTestRequest(proactiveTestRequest).Execute()

Test proactive message



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
	proactiveTestRequest := *openapiclient.NewProactiveTestRequest("ConversationId_example", "Message_example") // ProactiveTestRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExternalAPIAPI.ApiV1MessagesTestPost(context.Background()).ProactiveTestRequest(proactiveTestRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExternalAPIAPI.ApiV1MessagesTestPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1MessagesTestPost`: MessageResponse
	fmt.Fprintf(os.Stdout, "Response from `ExternalAPIAPI.ApiV1MessagesTestPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1MessagesTestPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **proactiveTestRequest** | [**ProactiveTestRequest**](ProactiveTestRequest.md) |  | 

### Return type

[**MessageResponse**](MessageResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1NotifyPost

> ExternalNotifyResponse ApiV1NotifyPost(ctx).ExternalNotifyRequest(externalNotifyRequest).Execute()

Send notification (External API)



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
	externalNotifyRequest := *openapiclient.NewExternalNotifyRequest("my-project-key", "Hello from external API") // ExternalNotifyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExternalAPIAPI.ApiV1NotifyPost(context.Background()).ExternalNotifyRequest(externalNotifyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExternalAPIAPI.ApiV1NotifyPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1NotifyPost`: ExternalNotifyResponse
	fmt.Fprintf(os.Stdout, "Response from `ExternalAPIAPI.ApiV1NotifyPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1NotifyPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **externalNotifyRequest** | [**ExternalNotifyRequest**](ExternalNotifyRequest.md) |  | 

### Return type

[**ExternalNotifyResponse**](ExternalNotifyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1ProvisionNotifyKeyDisablePost

> ProvisionResponse ApiV1ProvisionNotifyKeyDisablePost(ctx, notifyKey).Execute()

Disable project provision



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
	notifyKey := "notifyKey_example" // string | Project notify key

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExternalAPIAPI.ApiV1ProvisionNotifyKeyDisablePost(context.Background(), notifyKey).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExternalAPIAPI.ApiV1ProvisionNotifyKeyDisablePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1ProvisionNotifyKeyDisablePost`: ProvisionResponse
	fmt.Fprintf(os.Stdout, "Response from `ExternalAPIAPI.ApiV1ProvisionNotifyKeyDisablePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**notifyKey** | **string** | Project notify key | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1ProvisionNotifyKeyDisablePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1ProvisionNotifyKeyEnablePost

> ProvisionResponse ApiV1ProvisionNotifyKeyEnablePost(ctx, notifyKey).Execute()

Enable project provision



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
	notifyKey := "notifyKey_example" // string | Project notify key

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExternalAPIAPI.ApiV1ProvisionNotifyKeyEnablePost(context.Background(), notifyKey).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExternalAPIAPI.ApiV1ProvisionNotifyKeyEnablePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1ProvisionNotifyKeyEnablePost`: ProvisionResponse
	fmt.Fprintf(os.Stdout, "Response from `ExternalAPIAPI.ApiV1ProvisionNotifyKeyEnablePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**notifyKey** | **string** | Project notify key | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1ProvisionNotifyKeyEnablePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1ProvisionNotifyKeyGet

> ProvisionResponse ApiV1ProvisionNotifyKeyGet(ctx, notifyKey).Execute()

Get project provision



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
	notifyKey := "notifyKey_example" // string | Project notify key

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExternalAPIAPI.ApiV1ProvisionNotifyKeyGet(context.Background(), notifyKey).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExternalAPIAPI.ApiV1ProvisionNotifyKeyGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1ProvisionNotifyKeyGet`: ProvisionResponse
	fmt.Fprintf(os.Stdout, "Response from `ExternalAPIAPI.ApiV1ProvisionNotifyKeyGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**notifyKey** | **string** | Project notify key | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1ProvisionNotifyKeyGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1ProvisionNotifyKeyPut

> ProvisionResponse ApiV1ProvisionNotifyKeyPut(ctx, notifyKey).ProvisionUpdateRequest(provisionUpdateRequest).Execute()

Update project provision



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
	notifyKey := "notifyKey_example" // string | Project notify key
	provisionUpdateRequest := *openapiclient.NewProvisionUpdateRequest() // ProvisionUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExternalAPIAPI.ApiV1ProvisionNotifyKeyPut(context.Background(), notifyKey).ProvisionUpdateRequest(provisionUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExternalAPIAPI.ApiV1ProvisionNotifyKeyPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1ProvisionNotifyKeyPut`: ProvisionResponse
	fmt.Fprintf(os.Stdout, "Response from `ExternalAPIAPI.ApiV1ProvisionNotifyKeyPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**notifyKey** | **string** | Project notify key | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiV1ProvisionNotifyKeyPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **provisionUpdateRequest** | [**ProvisionUpdateRequest**](ProvisionUpdateRequest.md) |  | 

### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiV1ProvisionPost

> ProvisionResponse ApiV1ProvisionPost(ctx).ProvisionCreateRequest(provisionCreateRequest).Execute()

Create project provision



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
	provisionCreateRequest := *openapiclient.NewProvisionCreateRequest("NotifyKey_example", "CompanyId_example", "ProjectName_example") // ProvisionCreateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExternalAPIAPI.ApiV1ProvisionPost(context.Background()).ProvisionCreateRequest(provisionCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExternalAPIAPI.ApiV1ProvisionPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiV1ProvisionPost`: ProvisionResponse
	fmt.Fprintf(os.Stdout, "Response from `ExternalAPIAPI.ApiV1ProvisionPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiV1ProvisionPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **provisionCreateRequest** | [**ProvisionCreateRequest**](ProvisionCreateRequest.md) |  | 

### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

