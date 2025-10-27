# \InternalAPIQueueAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InternalV1QueueClearPost**](InternalAPIQueueAPI.md#InternalV1QueueClearPost) | **Post** /internal/v1/queue/clear | Clear queue
[**InternalV1QueueStatsGet**](InternalAPIQueueAPI.md#InternalV1QueueStatsGet) | **Get** /internal/v1/queue/stats | Get queue statistics
[**InternalV1QueueStatusGet**](InternalAPIQueueAPI.md#InternalV1QueueStatusGet) | **Get** /internal/v1/queue/status | Get queue status



## InternalV1QueueClearPost

> QueueClearResponse InternalV1QueueClearPost(ctx).Execute()

Clear queue



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIQueueAPI.InternalV1QueueClearPost(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIQueueAPI.InternalV1QueueClearPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1QueueClearPost`: QueueClearResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIQueueAPI.InternalV1QueueClearPost`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1QueueClearPostRequest struct via the builder pattern


### Return type

[**QueueClearResponse**](QueueClearResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1QueueStatsGet

> QueueStatsResponse InternalV1QueueStatsGet(ctx).Execute()

Get queue statistics



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIQueueAPI.InternalV1QueueStatsGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIQueueAPI.InternalV1QueueStatsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1QueueStatsGet`: QueueStatsResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIQueueAPI.InternalV1QueueStatsGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1QueueStatsGetRequest struct via the builder pattern


### Return type

[**QueueStatsResponse**](QueueStatsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1QueueStatusGet

> QueueStatusResponse InternalV1QueueStatusGet(ctx).Execute()

Get queue status



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIQueueAPI.InternalV1QueueStatusGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIQueueAPI.InternalV1QueueStatusGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1QueueStatusGet`: QueueStatusResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIQueueAPI.InternalV1QueueStatusGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1QueueStatusGetRequest struct via the builder pattern


### Return type

[**QueueStatusResponse**](QueueStatusResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

