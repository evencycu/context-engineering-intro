# \InternalAPIDestinationsAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InternalV1DestinationsBotBotIdGet**](InternalAPIDestinationsAPI.md#InternalV1DestinationsBotBotIdGet) | **Get** /internal/v1/destinations/bot/{botId} | Get destinations by bot
[**InternalV1DestinationsGet**](InternalAPIDestinationsAPI.md#InternalV1DestinationsGet) | **Get** /internal/v1/destinations | List destinations
[**InternalV1DestinationsIdDelete**](InternalAPIDestinationsAPI.md#InternalV1DestinationsIdDelete) | **Delete** /internal/v1/destinations/{id} | Delete destination
[**InternalV1DestinationsIdGet**](InternalAPIDestinationsAPI.md#InternalV1DestinationsIdGet) | **Get** /internal/v1/destinations/{id} | Get destination
[**InternalV1DestinationsIdPut**](InternalAPIDestinationsAPI.md#InternalV1DestinationsIdPut) | **Put** /internal/v1/destinations/{id} | Update destination
[**InternalV1DestinationsIdTargetsPatch**](InternalAPIDestinationsAPI.md#InternalV1DestinationsIdTargetsPatch) | **Patch** /internal/v1/destinations/{id}/targets | Update destination targets
[**InternalV1DestinationsIdValidatePost**](InternalAPIDestinationsAPI.md#InternalV1DestinationsIdValidatePost) | **Post** /internal/v1/destinations/{id}/validate | Validate destination targets
[**InternalV1DestinationsPost**](InternalAPIDestinationsAPI.md#InternalV1DestinationsPost) | **Post** /internal/v1/destinations | Create destination
[**InternalV1DestinationsProjectProjectIdGet**](InternalAPIDestinationsAPI.md#InternalV1DestinationsProjectProjectIdGet) | **Get** /internal/v1/destinations/project/{projectId} | Get destinations by project
[**InternalV1DestinationsSearchGet**](InternalAPIDestinationsAPI.md#InternalV1DestinationsSearchGet) | **Get** /internal/v1/destinations/search | Search destinations



## InternalV1DestinationsBotBotIdGet

> DestinationListResponse InternalV1DestinationsBotBotIdGet(ctx, botId).Execute()

Get destinations by bot



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
	botId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsBotBotIdGet(context.Background(), botId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsBotBotIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1DestinationsBotBotIdGet`: DestinationListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIDestinationsAPI.InternalV1DestinationsBotBotIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**botId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1DestinationsBotBotIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DestinationListResponse**](DestinationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1DestinationsGet

> DestinationListResponse InternalV1DestinationsGet(ctx).Limit(limit).Offset(offset).Search(search).Execute()

List destinations



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
	resp, r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsGet(context.Background()).Limit(limit).Offset(offset).Search(search).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1DestinationsGet`: DestinationListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIDestinationsAPI.InternalV1DestinationsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1DestinationsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | Number of items to return | [default to 10]
 **offset** | **int32** | Number of items to skip | [default to 0]
 **search** | **string** | Search term | 

### Return type

[**DestinationListResponse**](DestinationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1DestinationsIdDelete

> InternalV1DestinationsIdDelete(ctx, id).Execute()

Delete destination



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
	r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsIdDelete``: %v\n", err)
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

Other parameters are passed through a pointer to a apiInternalV1DestinationsIdDeleteRequest struct via the builder pattern


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


## InternalV1DestinationsIdGet

> DestinationResponse InternalV1DestinationsIdGet(ctx, id).Execute()

Get destination



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
	resp, r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1DestinationsIdGet`: DestinationResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIDestinationsAPI.InternalV1DestinationsIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1DestinationsIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DestinationResponse**](DestinationResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1DestinationsIdPut

> DestinationResponse InternalV1DestinationsIdPut(ctx, id).DestinationUpdateRequest(destinationUpdateRequest).Execute()

Update destination



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
	destinationUpdateRequest := *openapiclient.NewDestinationUpdateRequest() // DestinationUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsIdPut(context.Background(), id).DestinationUpdateRequest(destinationUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1DestinationsIdPut`: DestinationResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIDestinationsAPI.InternalV1DestinationsIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1DestinationsIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **destinationUpdateRequest** | [**DestinationUpdateRequest**](DestinationUpdateRequest.md) |  | 

### Return type

[**DestinationResponse**](DestinationResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1DestinationsIdTargetsPatch

> DestinationResponse InternalV1DestinationsIdTargetsPatch(ctx, id).DestinationTargetsUpdateRequest(destinationTargetsUpdateRequest).Execute()

Update destination targets



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
	destinationTargetsUpdateRequest := *openapiclient.NewDestinationTargetsUpdateRequest([]openapiclient.DestinationTarget{*openapiclient.NewDestinationTarget("Type_example", "ConversationId_example", "DisplayName_example", "TenantId_example")}) // DestinationTargetsUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsIdTargetsPatch(context.Background(), id).DestinationTargetsUpdateRequest(destinationTargetsUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsIdTargetsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1DestinationsIdTargetsPatch`: DestinationResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIDestinationsAPI.InternalV1DestinationsIdTargetsPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1DestinationsIdTargetsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **destinationTargetsUpdateRequest** | [**DestinationTargetsUpdateRequest**](DestinationTargetsUpdateRequest.md) |  | 

### Return type

[**DestinationResponse**](DestinationResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1DestinationsIdValidatePost

> DestinationValidationResponse InternalV1DestinationsIdValidatePost(ctx, id).Execute()

Validate destination targets



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
	resp, r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsIdValidatePost(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsIdValidatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1DestinationsIdValidatePost`: DestinationValidationResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIDestinationsAPI.InternalV1DestinationsIdValidatePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1DestinationsIdValidatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DestinationValidationResponse**](DestinationValidationResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1DestinationsPost

> DestinationResponse InternalV1DestinationsPost(ctx).DestinationCreateRequest(destinationCreateRequest).Execute()

Create destination



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
	destinationCreateRequest := *openapiclient.NewDestinationCreateRequest("ProjectId_example", "BotId_example", "Type_example", []openapiclient.DestinationTarget{*openapiclient.NewDestinationTarget("Type_example", "ConversationId_example", "DisplayName_example", "TenantId_example")}) // DestinationCreateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsPost(context.Background()).DestinationCreateRequest(destinationCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1DestinationsPost`: DestinationResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIDestinationsAPI.InternalV1DestinationsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1DestinationsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **destinationCreateRequest** | [**DestinationCreateRequest**](DestinationCreateRequest.md) |  | 

### Return type

[**DestinationResponse**](DestinationResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1DestinationsProjectProjectIdGet

> DestinationListResponse InternalV1DestinationsProjectProjectIdGet(ctx, projectId).Execute()

Get destinations by project



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
	resp, r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsProjectProjectIdGet(context.Background(), projectId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsProjectProjectIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1DestinationsProjectProjectIdGet`: DestinationListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIDestinationsAPI.InternalV1DestinationsProjectProjectIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1DestinationsProjectProjectIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DestinationListResponse**](DestinationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1DestinationsSearchGet

> DestinationListResponse InternalV1DestinationsSearchGet(ctx).Search(search).Type_(type_).Status(status).Execute()

Search destinations



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
	search := "search_example" // string | Search term (optional)
	type_ := "type__example" // string |  (optional)
	status := "status_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIDestinationsAPI.InternalV1DestinationsSearchGet(context.Background()).Search(search).Type_(type_).Status(status).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIDestinationsAPI.InternalV1DestinationsSearchGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1DestinationsSearchGet`: DestinationListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIDestinationsAPI.InternalV1DestinationsSearchGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1DestinationsSearchGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **search** | **string** | Search term | 
 **type_** | **string** |  | 
 **status** | **string** |  | 

### Return type

[**DestinationListResponse**](DestinationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

