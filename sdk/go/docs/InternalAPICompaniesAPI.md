# \InternalAPICompaniesAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InternalV1CompaniesGet**](InternalAPICompaniesAPI.md#InternalV1CompaniesGet) | **Get** /internal/v1/companies | List companies
[**InternalV1CompaniesIdBillingPatch**](InternalAPICompaniesAPI.md#InternalV1CompaniesIdBillingPatch) | **Patch** /internal/v1/companies/{id}/billing | Update company billing status
[**InternalV1CompaniesIdDelete**](InternalAPICompaniesAPI.md#InternalV1CompaniesIdDelete) | **Delete** /internal/v1/companies/{id} | Delete company
[**InternalV1CompaniesIdGet**](InternalAPICompaniesAPI.md#InternalV1CompaniesIdGet) | **Get** /internal/v1/companies/{id} | Get company
[**InternalV1CompaniesIdPut**](InternalAPICompaniesAPI.md#InternalV1CompaniesIdPut) | **Put** /internal/v1/companies/{id} | Update company
[**InternalV1CompaniesIdStatusPatch**](InternalAPICompaniesAPI.md#InternalV1CompaniesIdStatusPatch) | **Patch** /internal/v1/companies/{id}/status | Update company status
[**InternalV1CompaniesPost**](InternalAPICompaniesAPI.md#InternalV1CompaniesPost) | **Post** /internal/v1/companies | Create company



## InternalV1CompaniesGet

> CompanyListResponse InternalV1CompaniesGet(ctx).Limit(limit).Offset(offset).Search(search).Execute()

List companies



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
	resp, r, err := apiClient.InternalAPICompaniesAPI.InternalV1CompaniesGet(context.Background()).Limit(limit).Offset(offset).Search(search).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPICompaniesAPI.InternalV1CompaniesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1CompaniesGet`: CompanyListResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPICompaniesAPI.InternalV1CompaniesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1CompaniesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | Number of items to return | [default to 10]
 **offset** | **int32** | Number of items to skip | [default to 0]
 **search** | **string** | Search term | 

### Return type

[**CompanyListResponse**](CompanyListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1CompaniesIdBillingPatch

> CompanyResponse InternalV1CompaniesIdBillingPatch(ctx, id).ProjectBillingUpdateRequest(projectBillingUpdateRequest).Execute()

Update company billing status



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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Company ID
	projectBillingUpdateRequest := *openapiclient.NewProjectBillingUpdateRequest() // ProjectBillingUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPICompaniesAPI.InternalV1CompaniesIdBillingPatch(context.Background(), id).ProjectBillingUpdateRequest(projectBillingUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPICompaniesAPI.InternalV1CompaniesIdBillingPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1CompaniesIdBillingPatch`: CompanyResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPICompaniesAPI.InternalV1CompaniesIdBillingPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Company ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1CompaniesIdBillingPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **projectBillingUpdateRequest** | [**ProjectBillingUpdateRequest**](ProjectBillingUpdateRequest.md) |  | 

### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1CompaniesIdDelete

> InternalV1CompaniesIdDelete(ctx, id).Execute()

Delete company



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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Company ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.InternalAPICompaniesAPI.InternalV1CompaniesIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPICompaniesAPI.InternalV1CompaniesIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Company ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1CompaniesIdDeleteRequest struct via the builder pattern


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


## InternalV1CompaniesIdGet

> CompanyResponse InternalV1CompaniesIdGet(ctx, id).Execute()

Get company



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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Company ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPICompaniesAPI.InternalV1CompaniesIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPICompaniesAPI.InternalV1CompaniesIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1CompaniesIdGet`: CompanyResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPICompaniesAPI.InternalV1CompaniesIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Company ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1CompaniesIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1CompaniesIdPut

> CompanyResponse InternalV1CompaniesIdPut(ctx, id).CompanyUpdateRequest(companyUpdateRequest).Execute()

Update company



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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Company ID
	companyUpdateRequest := *openapiclient.NewCompanyUpdateRequest() // CompanyUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPICompaniesAPI.InternalV1CompaniesIdPut(context.Background(), id).CompanyUpdateRequest(companyUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPICompaniesAPI.InternalV1CompaniesIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1CompaniesIdPut`: CompanyResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPICompaniesAPI.InternalV1CompaniesIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Company ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1CompaniesIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **companyUpdateRequest** | [**CompanyUpdateRequest**](CompanyUpdateRequest.md) |  | 

### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1CompaniesIdStatusPatch

> CompanyResponse InternalV1CompaniesIdStatusPatch(ctx, id).CompanyStatusUpdateRequest(companyStatusUpdateRequest).Execute()

Update company status



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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | Company ID
	companyStatusUpdateRequest := *openapiclient.NewCompanyStatusUpdateRequest("Status_example") // CompanyStatusUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPICompaniesAPI.InternalV1CompaniesIdStatusPatch(context.Background(), id).CompanyStatusUpdateRequest(companyStatusUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPICompaniesAPI.InternalV1CompaniesIdStatusPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1CompaniesIdStatusPatch`: CompanyResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPICompaniesAPI.InternalV1CompaniesIdStatusPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Company ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1CompaniesIdStatusPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **companyStatusUpdateRequest** | [**CompanyStatusUpdateRequest**](CompanyStatusUpdateRequest.md) |  | 

### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1CompaniesPost

> CompanyResponse InternalV1CompaniesPost(ctx).CompanyCreateRequest(companyCreateRequest).Execute()

Create company



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
	companyCreateRequest := *openapiclient.NewCompanyCreateRequest("Name_example", "ContactEmail_example") // CompanyCreateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPICompaniesAPI.InternalV1CompaniesPost(context.Background()).CompanyCreateRequest(companyCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPICompaniesAPI.InternalV1CompaniesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1CompaniesPost`: CompanyResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPICompaniesAPI.InternalV1CompaniesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1CompaniesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **companyCreateRequest** | [**CompanyCreateRequest**](CompanyCreateRequest.md) |  | 

### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

