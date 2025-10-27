# \InternalAPIMonitoringAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**InternalV1MonitoringAlertsGet**](InternalAPIMonitoringAPI.md#InternalV1MonitoringAlertsGet) | **Get** /internal/v1/monitoring/alerts | Get alert status
[**InternalV1MonitoringBusinessGet**](InternalAPIMonitoringAPI.md#InternalV1MonitoringBusinessGet) | **Get** /internal/v1/monitoring/business | Get business health
[**InternalV1MonitoringDashboardGet**](InternalAPIMonitoringAPI.md#InternalV1MonitoringDashboardGet) | **Get** /internal/v1/monitoring/dashboard | Get monitoring dashboard
[**InternalV1MonitoringHealthGet**](InternalAPIMonitoringAPI.md#InternalV1MonitoringHealthGet) | **Get** /internal/v1/monitoring/health | Get system health
[**InternalV1MonitoringPerformanceGet**](InternalAPIMonitoringAPI.md#InternalV1MonitoringPerformanceGet) | **Get** /internal/v1/monitoring/performance | Get performance metrics



## InternalV1MonitoringAlertsGet

> AlertStatusResponse InternalV1MonitoringAlertsGet(ctx).Execute()

Get alert status



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
	resp, r, err := apiClient.InternalAPIMonitoringAPI.InternalV1MonitoringAlertsGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIMonitoringAPI.InternalV1MonitoringAlertsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1MonitoringAlertsGet`: AlertStatusResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIMonitoringAPI.InternalV1MonitoringAlertsGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1MonitoringAlertsGetRequest struct via the builder pattern


### Return type

[**AlertStatusResponse**](AlertStatusResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1MonitoringBusinessGet

> BusinessHealthResponse InternalV1MonitoringBusinessGet(ctx).Execute()

Get business health



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
	resp, r, err := apiClient.InternalAPIMonitoringAPI.InternalV1MonitoringBusinessGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIMonitoringAPI.InternalV1MonitoringBusinessGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1MonitoringBusinessGet`: BusinessHealthResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIMonitoringAPI.InternalV1MonitoringBusinessGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1MonitoringBusinessGetRequest struct via the builder pattern


### Return type

[**BusinessHealthResponse**](BusinessHealthResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1MonitoringDashboardGet

> DashboardResponse InternalV1MonitoringDashboardGet(ctx).Execute()

Get monitoring dashboard



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
	resp, r, err := apiClient.InternalAPIMonitoringAPI.InternalV1MonitoringDashboardGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIMonitoringAPI.InternalV1MonitoringDashboardGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1MonitoringDashboardGet`: DashboardResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIMonitoringAPI.InternalV1MonitoringDashboardGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1MonitoringDashboardGetRequest struct via the builder pattern


### Return type

[**DashboardResponse**](DashboardResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1MonitoringHealthGet

> SystemHealthResponse InternalV1MonitoringHealthGet(ctx).Execute()

Get system health



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
	resp, r, err := apiClient.InternalAPIMonitoringAPI.InternalV1MonitoringHealthGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIMonitoringAPI.InternalV1MonitoringHealthGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1MonitoringHealthGet`: SystemHealthResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIMonitoringAPI.InternalV1MonitoringHealthGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1MonitoringHealthGetRequest struct via the builder pattern


### Return type

[**SystemHealthResponse**](SystemHealthResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InternalV1MonitoringPerformanceGet

> PerformanceMetricsResponse InternalV1MonitoringPerformanceGet(ctx).Execute()

Get performance metrics



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
	resp, r, err := apiClient.InternalAPIMonitoringAPI.InternalV1MonitoringPerformanceGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InternalAPIMonitoringAPI.InternalV1MonitoringPerformanceGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InternalV1MonitoringPerformanceGet`: PerformanceMetricsResponse
	fmt.Fprintf(os.Stdout, "Response from `InternalAPIMonitoringAPI.InternalV1MonitoringPerformanceGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiInternalV1MonitoringPerformanceGetRequest struct via the builder pattern


### Return type

[**PerformanceMetricsResponse**](PerformanceMetricsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

