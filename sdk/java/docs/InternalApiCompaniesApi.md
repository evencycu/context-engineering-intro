# InternalApiCompaniesApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1CompaniesGet**](InternalApiCompaniesApi.md#internalV1CompaniesGet) | **GET** /internal/v1/companies | List companies |
| [**internalV1CompaniesIdBillingPatch**](InternalApiCompaniesApi.md#internalV1CompaniesIdBillingPatch) | **PATCH** /internal/v1/companies/{id}/billing | Update company billing status |
| [**internalV1CompaniesIdDelete**](InternalApiCompaniesApi.md#internalV1CompaniesIdDelete) | **DELETE** /internal/v1/companies/{id} | Delete company |
| [**internalV1CompaniesIdGet**](InternalApiCompaniesApi.md#internalV1CompaniesIdGet) | **GET** /internal/v1/companies/{id} | Get company |
| [**internalV1CompaniesIdPut**](InternalApiCompaniesApi.md#internalV1CompaniesIdPut) | **PUT** /internal/v1/companies/{id} | Update company |
| [**internalV1CompaniesIdStatusPatch**](InternalApiCompaniesApi.md#internalV1CompaniesIdStatusPatch) | **PATCH** /internal/v1/companies/{id}/status | Update company status |
| [**internalV1CompaniesPost**](InternalApiCompaniesApi.md#internalV1CompaniesPost) | **POST** /internal/v1/companies | Create company |


<a id="internalV1CompaniesGet"></a>
# **internalV1CompaniesGet**
> CompanyListResponse internalV1CompaniesGet(limit, offset, search)

List companies

Get a list of companies with pagination

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiCompaniesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiCompaniesApi apiInstance = new InternalApiCompaniesApi(defaultClient);
    Integer limit = 10; // Integer | Number of items to return
    Integer offset = 0; // Integer | Number of items to skip
    String search = "search_example"; // String | Search term
    try {
      CompanyListResponse result = apiInstance.internalV1CompaniesGet(limit, offset, search);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiCompaniesApi#internalV1CompaniesGet");
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
| **search** | **String**| Search term | [optional] |

### Return type

[**CompanyListResponse**](CompanyListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of companies |  -  |

<a id="internalV1CompaniesIdBillingPatch"></a>
# **internalV1CompaniesIdBillingPatch**
> CompanyResponse internalV1CompaniesIdBillingPatch(id, projectBillingUpdateRequest)

Update company billing status

Update company billing status by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiCompaniesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiCompaniesApi apiInstance = new InternalApiCompaniesApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | Company ID
    ProjectBillingUpdateRequest projectBillingUpdateRequest = new ProjectBillingUpdateRequest(); // ProjectBillingUpdateRequest | 
    try {
      CompanyResponse result = apiInstance.internalV1CompaniesIdBillingPatch(id, projectBillingUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiCompaniesApi#internalV1CompaniesIdBillingPatch");
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
| **id** | **UUID**| Company ID | |
| **projectBillingUpdateRequest** | [**ProjectBillingUpdateRequest**](ProjectBillingUpdateRequest.md)|  | |

### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Company billing status updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1CompaniesIdDelete"></a>
# **internalV1CompaniesIdDelete**
> internalV1CompaniesIdDelete(id)

Delete company

Delete company by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiCompaniesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiCompaniesApi apiInstance = new InternalApiCompaniesApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | Company ID
    try {
      apiInstance.internalV1CompaniesIdDelete(id);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiCompaniesApi#internalV1CompaniesIdDelete");
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
| **id** | **UUID**| Company ID | |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Company deleted successfully |  -  |
| **404** | Not found |  -  |

<a id="internalV1CompaniesIdGet"></a>
# **internalV1CompaniesIdGet**
> CompanyResponse internalV1CompaniesIdGet(id)

Get company

Get company by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiCompaniesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiCompaniesApi apiInstance = new InternalApiCompaniesApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | Company ID
    try {
      CompanyResponse result = apiInstance.internalV1CompaniesIdGet(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiCompaniesApi#internalV1CompaniesIdGet");
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
| **id** | **UUID**| Company ID | |

### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Company details |  -  |
| **404** | Not found |  -  |

<a id="internalV1CompaniesIdPut"></a>
# **internalV1CompaniesIdPut**
> CompanyResponse internalV1CompaniesIdPut(id, companyUpdateRequest)

Update company

Update company by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiCompaniesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiCompaniesApi apiInstance = new InternalApiCompaniesApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | Company ID
    CompanyUpdateRequest companyUpdateRequest = new CompanyUpdateRequest(); // CompanyUpdateRequest | 
    try {
      CompanyResponse result = apiInstance.internalV1CompaniesIdPut(id, companyUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiCompaniesApi#internalV1CompaniesIdPut");
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
| **id** | **UUID**| Company ID | |
| **companyUpdateRequest** | [**CompanyUpdateRequest**](CompanyUpdateRequest.md)|  | |

### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Company updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1CompaniesIdStatusPatch"></a>
# **internalV1CompaniesIdStatusPatch**
> CompanyResponse internalV1CompaniesIdStatusPatch(id, companyStatusUpdateRequest)

Update company status

Update company status by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiCompaniesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiCompaniesApi apiInstance = new InternalApiCompaniesApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | Company ID
    CompanyStatusUpdateRequest companyStatusUpdateRequest = new CompanyStatusUpdateRequest(); // CompanyStatusUpdateRequest | 
    try {
      CompanyResponse result = apiInstance.internalV1CompaniesIdStatusPatch(id, companyStatusUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiCompaniesApi#internalV1CompaniesIdStatusPatch");
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
| **id** | **UUID**| Company ID | |
| **companyStatusUpdateRequest** | [**CompanyStatusUpdateRequest**](CompanyStatusUpdateRequest.md)|  | |

### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Company status updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1CompaniesPost"></a>
# **internalV1CompaniesPost**
> CompanyResponse internalV1CompaniesPost(companyCreateRequest)

Create company

Create a new company

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiCompaniesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiCompaniesApi apiInstance = new InternalApiCompaniesApi(defaultClient);
    CompanyCreateRequest companyCreateRequest = new CompanyCreateRequest(); // CompanyCreateRequest | 
    try {
      CompanyResponse result = apiInstance.internalV1CompaniesPost(companyCreateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiCompaniesApi#internalV1CompaniesPost");
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
| **companyCreateRequest** | [**CompanyCreateRequest**](CompanyCreateRequest.md)|  | |

### Return type

[**CompanyResponse**](CompanyResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Company created successfully |  -  |
| **400** | Bad request |  -  |

