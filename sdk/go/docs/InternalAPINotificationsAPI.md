# \InternalAPINotificationsAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InternalV1NotificationsDateRangeGet**](InternalAPINotificationsAPI.md#InternalV1NotificationsDateRangeGet) | **Get** /internal/v1/notifications/date-range | Get notifications by date range
[**InternalV1NotificationsGet**](InternalAPINotificationsAPI.md#InternalV1NotificationsGet) | **Get** /internal/v1/notifications | List notifications
[**InternalV1NotificationsIdDelete**](InternalAPINotificationsAPI.md#InternalV1NotificationsIdDelete) | **Delete** /internal/v1/notifications/{id} | Cancel notification
[**InternalV1NotificationsIdGet**](InternalAPINotificationsAPI.md#InternalV1NotificationsIdGet) | **Get** /internal/v1/notifications/{id} | Get notification
[**InternalV1NotificationsIdRetryPost**](InternalAPINotificationsAPI.md#InternalV1NotificationsIdRetryPost) | **Post** /internal/v1/notifications/{id}/retry | Retry notification
[**InternalV1NotificationsPost**](InternalAPINotificationsAPI.md#InternalV1NotificationsPost) | **Post** /internal/v1/notifications | Send notification
[**InternalV1NotificationsProjectProjectIdGet**](InternalAPINotificationsAPI.md#InternalV1NotificationsProjectProjectIdGet) | **Get** /internal/v1/notifications/project/{projectId} | Get notifications by project
[**InternalV1NotificationsSenderSenderIdGet**](InternalAPINotificationsAPI.md#InternalV1NotificationsSenderSenderIdGet) | **Get** /internal/v1/notifications/sender/{senderId} | Get notifications by sender
[**InternalV1NotificationsStatusStatusGet**](InternalAPINotificationsAPI.md#InternalV1NotificationsStatusStatusGet) | **Get** /internal/v1/notifications/status/{status} | Get notifications by status



## InternalV1NotificationsDateRangeGet

> NotificationDateRangeResponse InternalV1NotificationsDateRangeGet(ctx).StartDate(startDate).EndDate(endDate).Execute()

Get notifications by date range



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	startDate := time.Now() // time.Time | 
	endDate := time.Now() // time.Time | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPINotificationsAPI.InternalV1NotificationsDateRangeGet(context.Background()).StartDate(startDate).EndDate(endDate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPINotificationsAPI.InternalV1NotificationsDateRangeGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1NotificationsDateRangeGet`: NotificationDateRangeResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPINotificationsAPI.InternalV1NotificationsDateRangeGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1NotificationsDateRangeGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **startDate** | **time.Time** |  | 
 **endDate** | **time.Time** |  | 

### Return type

[**NotificationDateRangeResponse**](NotificationDateRangeResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1NotificationsGet

> NotificationListResponse InternalV1NotificationsGet(ctx).Limit(limit).Offset(offset).Search(search).Status(status).Execute()

List notifications



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
	status := "status_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPINotificationsAPI.InternalV1NotificationsGet(context.Background()).Limit(limit).Offset(offset).Search(search).Status(status).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPINotificationsAPI.InternalV1NotificationsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1NotificationsGet`: NotificationListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPINotificationsAPI.InternalV1NotificationsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1NotificationsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | Number of items to return | [default to 10]
 **offset** | **int32** | Number of items to skip | [default to 0]
 **search** | **string** | Search term | 
 **status** | **string** |  | 

### Return type

[**NotificationListResponse**](NotificationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1NotificationsIdDelete

> NotificationCancelResponse InternalV1NotificationsIdDelete(ctx, id).Execute()

Cancel notification



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
	resp, r, err := apiClient.InternalAPINotificationsAPI.InternalV1NotificationsIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPINotificationsAPI.InternalV1NotificationsIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1NotificationsIdDelete`: NotificationCancelResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPINotificationsAPI.InternalV1NotificationsIdDelete`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1NotificationsIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**NotificationCancelResponse**](NotificationCancelResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1NotificationsIdGet

> NotificationResponse InternalV1NotificationsIdGet(ctx, id).Execute()

Get notification



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
	resp, r, err := apiClient.InternalAPINotificationsAPI.InternalV1NotificationsIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPINotificationsAPI.InternalV1NotificationsIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1NotificationsIdGet`: NotificationResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPINotificationsAPI.InternalV1NotificationsIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1NotificationsIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**NotificationResponse**](NotificationResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1NotificationsIdRetryPost

> NotificationRetryResponse InternalV1NotificationsIdRetryPost(ctx, id).Execute()

Retry notification



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
	resp, r, err := apiClient.InternalAPINotificationsAPI.InternalV1NotificationsIdRetryPost(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPINotificationsAPI.InternalV1NotificationsIdRetryPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1NotificationsIdRetryPost`: NotificationRetryResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPINotificationsAPI.InternalV1NotificationsIdRetryPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1NotificationsIdRetryPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**NotificationRetryResponse**](NotificationRetryResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1NotificationsPost

> NotificationSendResponse InternalV1NotificationsPost(ctx).NotificationSendRequest(notificationSendRequest).Execute()

Send notification



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
	notificationSendRequest := *openapiclient.NewNotificationSendRequest("ProjectId_example", "MessageType_example", "Content_example", []string{"Targets_example"}) // NotificationSendRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPINotificationsAPI.InternalV1NotificationsPost(context.Background()).NotificationSendRequest(notificationSendRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPINotificationsAPI.InternalV1NotificationsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1NotificationsPost`: NotificationSendResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPINotificationsAPI.InternalV1NotificationsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1NotificationsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **notificationSendRequest** | [**NotificationSendRequest**](NotificationSendRequest.md) |  | 

### Return type

[**NotificationSendResponse**](NotificationSendResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1NotificationsProjectProjectIdGet

> NotificationListResponse InternalV1NotificationsProjectProjectIdGet(ctx, projectId).Execute()

Get notifications by project



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
	projectId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPINotificationsAPI.InternalV1NotificationsProjectProjectIdGet(context.Background(), projectId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPINotificationsAPI.InternalV1NotificationsProjectProjectIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1NotificationsProjectProjectIdGet`: NotificationListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPINotificationsAPI.InternalV1NotificationsProjectProjectIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1NotificationsProjectProjectIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**NotificationListResponse**](NotificationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1NotificationsSenderSenderIdGet

> NotificationListResponse InternalV1NotificationsSenderSenderIdGet(ctx, senderId).Execute()

Get notifications by sender



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
	senderId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPINotificationsAPI.InternalV1NotificationsSenderSenderIdGet(context.Background(), senderId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPINotificationsAPI.InternalV1NotificationsSenderSenderIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1NotificationsSenderSenderIdGet`: NotificationListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPINotificationsAPI.InternalV1NotificationsSenderSenderIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**senderId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1NotificationsSenderSenderIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**NotificationListResponse**](NotificationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1NotificationsStatusStatusGet

> NotificationListResponse InternalV1NotificationsStatusStatusGet(ctx, status).Execute()

Get notifications by status



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
	resp, r, err := apiClient.InternalAPINotificationsAPI.InternalV1NotificationsStatusStatusGet(context.Background(), status).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPINotificationsAPI.InternalV1NotificationsStatusStatusGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1NotificationsStatusStatusGet`: NotificationListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPINotificationsAPI.InternalV1NotificationsStatusStatusGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**status** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1NotificationsStatusStatusGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**NotificationListResponse**](NotificationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

