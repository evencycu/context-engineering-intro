# \InternalAPIBillingAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InternalV1BillingAnalyticsOverviewGet**](InternalAPIBillingAPI.md#InternalV1BillingAnalyticsOverviewGet) | **Get** /internal/v1/billing/analytics/overview | Get billing analytics overview
[**InternalV1BillingPlansGet**](InternalAPIBillingAPI.md#InternalV1BillingPlansGet) | **Get** /internal/v1/billing/plans | List billing plans
[**InternalV1BillingPlansIdGet**](InternalAPIBillingAPI.md#InternalV1BillingPlansIdGet) | **Get** /internal/v1/billing/plans/{id} | Get billing plan
[**InternalV1BillingPlansIdPut**](InternalAPIBillingAPI.md#InternalV1BillingPlansIdPut) | **Put** /internal/v1/billing/plans/{id} | Update billing plan
[**InternalV1BillingPlansPost**](InternalAPIBillingAPI.md#InternalV1BillingPlansPost) | **Post** /internal/v1/billing/plans | Create billing plan
[**InternalV1BillingProjectProjectIdGet**](InternalAPIBillingAPI.md#InternalV1BillingProjectProjectIdGet) | **Get** /internal/v1/billing/project/{projectId} | Get project billing
[**InternalV1BillingProjectProjectIdPut**](InternalAPIBillingAPI.md#InternalV1BillingProjectProjectIdPut) | **Put** /internal/v1/billing/project/{projectId} | Update project billing
[**InternalV1BillingUsageCompanyCompanyIdGet**](InternalAPIBillingAPI.md#InternalV1BillingUsageCompanyCompanyIdGet) | **Get** /internal/v1/billing/usage/company/{companyId} | Get company usage
[**InternalV1BillingUsageGet**](InternalAPIBillingAPI.md#InternalV1BillingUsageGet) | **Get** /internal/v1/billing/usage | Get usage records
[**InternalV1BillingUsageProjectProjectIdGet**](InternalAPIBillingAPI.md#InternalV1BillingUsageProjectProjectIdGet) | **Get** /internal/v1/billing/usage/project/{projectId} | Get project usage
[**InternalV1BillingUsageSummaryGet**](InternalAPIBillingAPI.md#InternalV1BillingUsageSummaryGet) | **Get** /internal/v1/billing/usage/summary | Get usage summary



## InternalV1BillingAnalyticsOverviewGet

> BillingAnalyticsResponse InternalV1BillingAnalyticsOverviewGet(ctx).StartDate(startDate).EndDate(endDate).Execute()

Get billing analytics overview



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
	startDate := time.Now() // string |  (optional)
	endDate := time.Now() // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingAnalyticsOverviewGet(context.Background()).StartDate(startDate).EndDate(endDate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingAnalyticsOverviewGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingAnalyticsOverviewGet`: BillingAnalyticsResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingAnalyticsOverviewGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingAnalyticsOverviewGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **startDate** | **string** |  | 
 **endDate** | **string** |  | 

### Return type

[**BillingAnalyticsResponse**](BillingAnalyticsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingPlansGet

> BillingPlansResponse InternalV1BillingPlansGet(ctx).Execute()

List billing plans



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
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingPlansGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingPlansGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingPlansGet`: BillingPlansResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingPlansGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingPlansGetRequest struct via the builder pattern


### Return type

[**BillingPlansResponse**](BillingPlansResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingPlansIdGet

> BillingPlanResponse InternalV1BillingPlansIdGet(ctx, id).Execute()

Get billing plan



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
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingPlansIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingPlansIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingPlansIdGet`: BillingPlanResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingPlansIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingPlansIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**BillingPlanResponse**](BillingPlanResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingPlansIdPut

> BillingPlanResponse InternalV1BillingPlansIdPut(ctx, id).BillingPlanUpdateRequest(billingPlanUpdateRequest).Execute()

Update billing plan



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
	billingPlanUpdateRequest := *openapiclient.NewBillingPlanUpdateRequest() // BillingPlanUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingPlansIdPut(context.Background(), id).BillingPlanUpdateRequest(billingPlanUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingPlansIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingPlansIdPut`: BillingPlanResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingPlansIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingPlansIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **billingPlanUpdateRequest** | [**BillingPlanUpdateRequest**](BillingPlanUpdateRequest.md) |  | 

### Return type

[**BillingPlanResponse**](BillingPlanResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingPlansPost

> BillingPlanResponse InternalV1BillingPlansPost(ctx).BillingPlanCreateRequest(billingPlanCreateRequest).Execute()

Create billing plan



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
	billingPlanCreateRequest := *openapiclient.NewBillingPlanCreateRequest("Name_example", float32(123), int32(123)) // BillingPlanCreateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingPlansPost(context.Background()).BillingPlanCreateRequest(billingPlanCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingPlansPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingPlansPost`: BillingPlanResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingPlansPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingPlansPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **billingPlanCreateRequest** | [**BillingPlanCreateRequest**](BillingPlanCreateRequest.md) |  | 

### Return type

[**BillingPlanResponse**](BillingPlanResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingProjectProjectIdGet

> ProjectBillingResponse InternalV1BillingProjectProjectIdGet(ctx, projectId).Execute()

Get project billing



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
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingProjectProjectIdGet(context.Background(), projectId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingProjectProjectIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingProjectProjectIdGet`: ProjectBillingResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingProjectProjectIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingProjectProjectIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ProjectBillingResponse**](ProjectBillingResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingProjectProjectIdPut

> ProjectBillingResponse InternalV1BillingProjectProjectIdPut(ctx, projectId).ProjectBillingUpdateRequest(projectBillingUpdateRequest).Execute()

Update project billing



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
	projectBillingUpdateRequest := *openapiclient.NewProjectBillingUpdateRequest() // ProjectBillingUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingProjectProjectIdPut(context.Background(), projectId).ProjectBillingUpdateRequest(projectBillingUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingProjectProjectIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingProjectProjectIdPut`: ProjectBillingResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingProjectProjectIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingProjectProjectIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **projectBillingUpdateRequest** | [**ProjectBillingUpdateRequest**](ProjectBillingUpdateRequest.md) |  | 

### Return type

[**ProjectBillingResponse**](ProjectBillingResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingUsageCompanyCompanyIdGet

> UsageRecordsResponse InternalV1BillingUsageCompanyCompanyIdGet(ctx, companyId).Execute()

Get company usage



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
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingUsageCompanyCompanyIdGet(context.Background(), companyId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingUsageCompanyCompanyIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingUsageCompanyCompanyIdGet`: UsageRecordsResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingUsageCompanyCompanyIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**companyId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingUsageCompanyCompanyIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**UsageRecordsResponse**](UsageRecordsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingUsageGet

> UsageRecordsResponse InternalV1BillingUsageGet(ctx).Limit(limit).Offset(offset).CompanyId(companyId).ProjectId(projectId).Execute()

Get usage records



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
	companyId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string |  (optional)
	projectId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingUsageGet(context.Background()).Limit(limit).Offset(offset).CompanyId(companyId).ProjectId(projectId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingUsageGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingUsageGet`: UsageRecordsResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingUsageGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingUsageGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **int32** | Number of items to return | [default to 10]
 **offset** | **int32** | Number of items to skip | [default to 0]
 **companyId** | **string** |  | 
 **projectId** | **string** |  | 

### Return type

[**UsageRecordsResponse**](UsageRecordsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingUsageProjectProjectIdGet

> UsageRecordsResponse InternalV1BillingUsageProjectProjectIdGet(ctx, projectId).Execute()

Get project usage



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
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingUsageProjectProjectIdGet(context.Background(), projectId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingUsageProjectProjectIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingUsageProjectProjectIdGet`: UsageRecordsResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingUsageProjectProjectIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingUsageProjectProjectIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**UsageRecordsResponse**](UsageRecordsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1BillingUsageSummaryGet

> UsageSummaryResponse InternalV1BillingUsageSummaryGet(ctx).CompanyId(companyId).ProjectId(projectId).StartDate(startDate).EndDate(endDate).Execute()

Get usage summary



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
	companyId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string |  (optional)
	projectId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string |  (optional)
	startDate := time.Now() // string |  (optional)
	endDate := time.Now() // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.InternalAPIBillingAPI.InternalV1BillingUsageSummaryGet(context.Background()).CompanyId(companyId).ProjectId(projectId).StartDate(startDate).EndDate(endDate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIBillingAPI.InternalV1BillingUsageSummaryGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1BillingUsageSummaryGet`: UsageSummaryResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIBillingAPI.InternalV1BillingUsageSummaryGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1BillingUsageSummaryGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **companyId** | **string** |  | 
 **projectId** | **string** |  | 
 **startDate** | **string** |  | 
 **endDate** | **string** |  | 

### Return type

[**UsageSummaryResponse**](UsageSummaryResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

