# InternalApiBillingApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1BillingAnalyticsOverviewGet**](InternalApiBillingApi.md#internalV1BillingAnalyticsOverviewGet) | **GET** /internal/v1/billing/analytics/overview | Get billing analytics overview |
| [**internalV1BillingPlansGet**](InternalApiBillingApi.md#internalV1BillingPlansGet) | **GET** /internal/v1/billing/plans | List billing plans |
| [**internalV1BillingPlansIdGet**](InternalApiBillingApi.md#internalV1BillingPlansIdGet) | **GET** /internal/v1/billing/plans/{id} | Get billing plan |
| [**internalV1BillingPlansIdPut**](InternalApiBillingApi.md#internalV1BillingPlansIdPut) | **PUT** /internal/v1/billing/plans/{id} | Update billing plan |
| [**internalV1BillingPlansPost**](InternalApiBillingApi.md#internalV1BillingPlansPost) | **POST** /internal/v1/billing/plans | Create billing plan |
| [**internalV1BillingProjectProjectIdGet**](InternalApiBillingApi.md#internalV1BillingProjectProjectIdGet) | **GET** /internal/v1/billing/project/{projectId} | Get project billing |
| [**internalV1BillingProjectProjectIdPut**](InternalApiBillingApi.md#internalV1BillingProjectProjectIdPut) | **PUT** /internal/v1/billing/project/{projectId} | Update project billing |
| [**internalV1BillingUsageCompanyCompanyIdGet**](InternalApiBillingApi.md#internalV1BillingUsageCompanyCompanyIdGet) | **GET** /internal/v1/billing/usage/company/{companyId} | Get company usage |
| [**internalV1BillingUsageGet**](InternalApiBillingApi.md#internalV1BillingUsageGet) | **GET** /internal/v1/billing/usage | Get usage records |
| [**internalV1BillingUsageProjectProjectIdGet**](InternalApiBillingApi.md#internalV1BillingUsageProjectProjectIdGet) | **GET** /internal/v1/billing/usage/project/{projectId} | Get project usage |
| [**internalV1BillingUsageSummaryGet**](InternalApiBillingApi.md#internalV1BillingUsageSummaryGet) | **GET** /internal/v1/billing/usage/summary | Get usage summary |


<a id="internalV1BillingAnalyticsOverviewGet"></a>
# **internalV1BillingAnalyticsOverviewGet**
> BillingAnalyticsResponse internalV1BillingAnalyticsOverviewGet(startDate, endDate)

Get billing analytics overview

Get billing analytics overview

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    LocalDate startDate = LocalDate.now(); // LocalDate | 
    LocalDate endDate = LocalDate.now(); // LocalDate | 
    try {
      BillingAnalyticsResponse result = apiInstance.internalV1BillingAnalyticsOverviewGet(startDate, endDate);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingAnalyticsOverviewGet");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **startDate** | **LocalDate**|  | [optional] |
| **endDate** | **LocalDate**|  | [optional] |

### Return type

[**BillingAnalyticsResponse**](BillingAnalyticsResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Billing analytics overview |  -  |

<a id="internalV1BillingPlansGet"></a>
# **internalV1BillingPlansGet**
> BillingPlansResponse internalV1BillingPlansGet()

List billing plans

Get a list of billing plans

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    try {
      BillingPlansResponse result = apiInstance.internalV1BillingPlansGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingPlansGet");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**BillingPlansResponse**](BillingPlansResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of billing plans |  -  |

<a id="internalV1BillingPlansIdGet"></a>
# **internalV1BillingPlansIdGet**
> BillingPlanResponse internalV1BillingPlansIdGet(id)

Get billing plan

Get billing plan by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      BillingPlanResponse result = apiInstance.internalV1BillingPlansIdGet(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingPlansIdGet");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **UUID**|  | |

### Return type

[**BillingPlanResponse**](BillingPlanResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Billing plan details |  -  |
| **404** | Not found |  -  |

<a id="internalV1BillingPlansIdPut"></a>
# **internalV1BillingPlansIdPut**
> BillingPlanResponse internalV1BillingPlansIdPut(id, billingPlanUpdateRequest)

Update billing plan

Update billing plan by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    BillingPlanUpdateRequest billingPlanUpdateRequest = new BillingPlanUpdateRequest(); // BillingPlanUpdateRequest | 
    try {
      BillingPlanResponse result = apiInstance.internalV1BillingPlansIdPut(id, billingPlanUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingPlansIdPut");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **UUID**|  | |
| **billingPlanUpdateRequest** | [**BillingPlanUpdateRequest**](BillingPlanUpdateRequest.md)|  | |

### Return type

[**BillingPlanResponse**](BillingPlanResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Billing plan updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1BillingPlansPost"></a>
# **internalV1BillingPlansPost**
> BillingPlanResponse internalV1BillingPlansPost(billingPlanCreateRequest)

Create billing plan

Create a new billing plan

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    BillingPlanCreateRequest billingPlanCreateRequest = new BillingPlanCreateRequest(); // BillingPlanCreateRequest | 
    try {
      BillingPlanResponse result = apiInstance.internalV1BillingPlansPost(billingPlanCreateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingPlansPost");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **billingPlanCreateRequest** | [**BillingPlanCreateRequest**](BillingPlanCreateRequest.md)|  | |

### Return type

[**BillingPlanResponse**](BillingPlanResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Billing plan created successfully |  -  |
| **400** | Bad request |  -  |

<a id="internalV1BillingProjectProjectIdGet"></a>
# **internalV1BillingProjectProjectIdGet**
> ProjectBillingResponse internalV1BillingProjectProjectIdGet(projectId)

Get project billing

Get project billing information

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    UUID projectId = UUID.randomUUID(); // UUID | 
    try {
      ProjectBillingResponse result = apiInstance.internalV1BillingProjectProjectIdGet(projectId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingProjectProjectIdGet");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **projectId** | **UUID**|  | |

### Return type

[**ProjectBillingResponse**](ProjectBillingResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project billing information |  -  |
| **404** | Not found |  -  |

<a id="internalV1BillingProjectProjectIdPut"></a>
# **internalV1BillingProjectProjectIdPut**
> ProjectBillingResponse internalV1BillingProjectProjectIdPut(projectId, projectBillingUpdateRequest)

Update project billing

Update project billing information

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    UUID projectId = UUID.randomUUID(); // UUID | 
    ProjectBillingUpdateRequest projectBillingUpdateRequest = new ProjectBillingUpdateRequest(); // ProjectBillingUpdateRequest | 
    try {
      ProjectBillingResponse result = apiInstance.internalV1BillingProjectProjectIdPut(projectId, projectBillingUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingProjectProjectIdPut");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **projectId** | **UUID**|  | |
| **projectBillingUpdateRequest** | [**ProjectBillingUpdateRequest**](ProjectBillingUpdateRequest.md)|  | |

### Return type

[**ProjectBillingResponse**](ProjectBillingResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project billing updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1BillingUsageCompanyCompanyIdGet"></a>
# **internalV1BillingUsageCompanyCompanyIdGet**
> UsageRecordsResponse internalV1BillingUsageCompanyCompanyIdGet(companyId)

Get company usage

Get usage records for a company

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    UUID companyId = UUID.randomUUID(); // UUID | 
    try {
      UsageRecordsResponse result = apiInstance.internalV1BillingUsageCompanyCompanyIdGet(companyId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingUsageCompanyCompanyIdGet");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **companyId** | **UUID**|  | |

### Return type

[**UsageRecordsResponse**](UsageRecordsResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Company usage records |  -  |
| **404** | Not found |  -  |

<a id="internalV1BillingUsageGet"></a>
# **internalV1BillingUsageGet**
> UsageRecordsResponse internalV1BillingUsageGet(limit, offset, companyId, projectId)

Get usage records

Get usage records with pagination

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    Integer limit = 10; // Integer | Number of items to return
    Integer offset = 0; // Integer | Number of items to skip
    UUID companyId = UUID.randomUUID(); // UUID | 
    UUID projectId = UUID.randomUUID(); // UUID | 
    try {
      UsageRecordsResponse result = apiInstance.internalV1BillingUsageGet(limit, offset, companyId, projectId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingUsageGet");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **limit** | **Integer**| Number of items to return | [optional] [default to 10] |
| **offset** | **Integer**| Number of items to skip | [optional] [default to 0] |
| **companyId** | **UUID**|  | [optional] |
| **projectId** | **UUID**|  | [optional] |

### Return type

[**UsageRecordsResponse**](UsageRecordsResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of usage records |  -  |

<a id="internalV1BillingUsageProjectProjectIdGet"></a>
# **internalV1BillingUsageProjectProjectIdGet**
> UsageRecordsResponse internalV1BillingUsageProjectProjectIdGet(projectId)

Get project usage

Get usage records for a project

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    UUID projectId = UUID.randomUUID(); // UUID | 
    try {
      UsageRecordsResponse result = apiInstance.internalV1BillingUsageProjectProjectIdGet(projectId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingUsageProjectProjectIdGet");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **projectId** | **UUID**|  | |

### Return type

[**UsageRecordsResponse**](UsageRecordsResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project usage records |  -  |
| **404** | Not found |  -  |

<a id="internalV1BillingUsageSummaryGet"></a>
# **internalV1BillingUsageSummaryGet**
> UsageSummaryResponse internalV1BillingUsageSummaryGet(companyId, projectId, startDate, endDate)

Get usage summary

Get usage summary statistics

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBillingApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBillingApi apiInstance = new InternalApiBillingApi(defaultClient);
    UUID companyId = UUID.randomUUID(); // UUID | 
    UUID projectId = UUID.randomUUID(); // UUID | 
    LocalDate startDate = LocalDate.now(); // LocalDate | 
    LocalDate endDate = LocalDate.now(); // LocalDate | 
    try {
      UsageSummaryResponse result = apiInstance.internalV1BillingUsageSummaryGet(companyId, projectId, startDate, endDate);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBillingApi#internalV1BillingUsageSummaryGet");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters

| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **companyId** | **UUID**|  | [optional] |
| **projectId** | **UUID**|  | [optional] |
| **startDate** | **LocalDate**|  | [optional] |
| **endDate** | **LocalDate**|  | [optional] |

### Return type

[**UsageSummaryResponse**](UsageSummaryResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Usage summary |  -  |

