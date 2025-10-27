# \InternalAPIProjectsAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InternalV1ProjectsCompanyCompanyIdGet**](InternalAPIProjectsAPI.md#InternalV1ProjectsCompanyCompanyIdGet) | **Get** /internal/v1/projects/company/{companyId} | Get projects by company
[**InternalV1ProjectsGet**](InternalAPIProjectsAPI.md#InternalV1ProjectsGet) | **Get** /internal/v1/projects | List projects
[**InternalV1ProjectsIdDelete**](InternalAPIProjectsAPI.md#InternalV1ProjectsIdDelete) | **Delete** /internal/v1/projects/{id} | Delete project
[**InternalV1ProjectsIdGet**](InternalAPIProjectsAPI.md#InternalV1ProjectsIdGet) | **Get** /internal/v1/projects/{id} | Get project
[**InternalV1ProjectsIdLimitsPatch**](InternalAPIProjectsAPI.md#InternalV1ProjectsIdLimitsPatch) | **Patch** /internal/v1/projects/{id}/limits | Update project limits
[**InternalV1ProjectsIdPut**](InternalAPIProjectsAPI.md#InternalV1ProjectsIdPut) | **Put** /internal/v1/projects/{id} | Update project
[**InternalV1ProjectsKeyKeyNameGet**](InternalAPIProjectsAPI.md#InternalV1ProjectsKeyKeyNameGet) | **Get** /internal/v1/projects/key/{keyName} | Get project by key name
[**InternalV1ProjectsPost**](InternalAPIProjectsAPI.md#InternalV1ProjectsPost) | **Post** /internal/v1/projects | Create project



## InternalV1ProjectsCompanyCompanyIdGet

> ProjectListResponse InternalV1ProjectsCompanyCompanyIdGet(ctx, companyId).Execute()

Get projects by company



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
	companyId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIProjectsAPI.InternalV1ProjectsCompanyCompanyIdGet(context.Background(), companyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIProjectsAPI.InternalV1ProjectsCompanyCompanyIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1ProjectsCompanyCompanyIdGet`: ProjectListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIProjectsAPI.InternalV1ProjectsCompanyCompanyIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**companyId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1ProjectsCompanyCompanyIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ProjectListResponse**](ProjectListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1ProjectsGet

> ProjectListResponse InternalV1ProjectsGet(ctx).Limit(limit).Offset(offset).Search(search).Execute()

List projects



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
	resp, r, err := apiClient.InternalAPIProjectsAPI.InternalV1ProjectsGet(context.Background()).Limit(limit).Offset(offset).Search(search).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIProjectsAPI.InternalV1ProjectsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1ProjectsGet`: ProjectListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIProjectsAPI.InternalV1ProjectsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1ProjectsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | Number of items to return | [default to 10]
 **offset** | **int32** | Number of items to skip | [default to 0]
 **search** | **string** | Search term | 

### Return type

[**ProjectListResponse**](ProjectListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1ProjectsIdDelete

> InternalV1ProjectsIdDelete(ctx, id).Execute()

Delete project



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
	r, err := apiClient.InternalAPIProjectsAPI.InternalV1ProjectsIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIProjectsAPI.InternalV1ProjectsIdDelete``: %v\n", err)
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

Other parameters are passed through a pointer to a apiInternalV1ProjectsIdDeleteRequest struct via the builder pattern


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


## InternalV1ProjectsIdGet

> ProjectResponse InternalV1ProjectsIdGet(ctx, id).Execute()

Get project



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
	resp, r, err := apiClient.InternalAPIProjectsAPI.InternalV1ProjectsIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIProjectsAPI.InternalV1ProjectsIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1ProjectsIdGet`: ProjectResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIProjectsAPI.InternalV1ProjectsIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1ProjectsIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1ProjectsIdLimitsPatch

> ProjectResponse InternalV1ProjectsIdLimitsPatch(ctx, id).ProjectLimitsUpdateRequest(projectLimitsUpdateRequest).Execute()

Update project limits



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
	projectLimitsUpdateRequest := *openapiclient.NewProjectLimitsUpdateRequest() // ProjectLimitsUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIProjectsAPI.InternalV1ProjectsIdLimitsPatch(context.Background(), id).ProjectLimitsUpdateRequest(projectLimitsUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIProjectsAPI.InternalV1ProjectsIdLimitsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1ProjectsIdLimitsPatch`: ProjectResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIProjectsAPI.InternalV1ProjectsIdLimitsPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1ProjectsIdLimitsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **projectLimitsUpdateRequest** | [**ProjectLimitsUpdateRequest**](ProjectLimitsUpdateRequest.md) |  | 

### Return type

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1ProjectsIdPut

> ProjectResponse InternalV1ProjectsIdPut(ctx, id).ProjectUpdateRequest(projectUpdateRequest).Execute()

Update project



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
	projectUpdateRequest := *openapiclient.NewProjectUpdateRequest() // ProjectUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIProjectsAPI.InternalV1ProjectsIdPut(context.Background(), id).ProjectUpdateRequest(projectUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIProjectsAPI.InternalV1ProjectsIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1ProjectsIdPut`: ProjectResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIProjectsAPI.InternalV1ProjectsIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1ProjectsIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **projectUpdateRequest** | [**ProjectUpdateRequest**](ProjectUpdateRequest.md) |  | 

### Return type

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1ProjectsKeyKeyNameGet

> ProjectResponse InternalV1ProjectsKeyKeyNameGet(ctx, keyName).Execute()

Get project by key name



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
	keyName := "keyName_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIProjectsAPI.InternalV1ProjectsKeyKeyNameGet(context.Background(), keyName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIProjectsAPI.InternalV1ProjectsKeyKeyNameGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1ProjectsKeyKeyNameGet`: ProjectResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIProjectsAPI.InternalV1ProjectsKeyKeyNameGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**keyName** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1ProjectsKeyKeyNameGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1ProjectsPost

> ProjectResponse InternalV1ProjectsPost(ctx).ProjectCreateRequest(projectCreateRequest).Execute()

Create project



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
	projectCreateRequest := *openapiclient.NewProjectCreateRequest("CompanyId_example", "NotifyKey_example", "Description_example", "CreatedBy_example") // ProjectCreateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIProjectsAPI.InternalV1ProjectsPost(context.Background()).ProjectCreateRequest(projectCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIProjectsAPI.InternalV1ProjectsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1ProjectsPost`: ProjectResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIProjectsAPI.InternalV1ProjectsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1ProjectsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **projectCreateRequest** | [**ProjectCreateRequest**](ProjectCreateRequest.md) |  | 

### Return type

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

