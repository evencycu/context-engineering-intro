# InternalApiProjectsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1ProjectsCompanyCompanyIdGet**](InternalApiProjectsApi.md#internalV1ProjectsCompanyCompanyIdGet) | **GET** /internal/v1/projects/company/{companyId} | Get projects by company |
| [**internalV1ProjectsGet**](InternalApiProjectsApi.md#internalV1ProjectsGet) | **GET** /internal/v1/projects | List projects |
| [**internalV1ProjectsIdDelete**](InternalApiProjectsApi.md#internalV1ProjectsIdDelete) | **DELETE** /internal/v1/projects/{id} | Delete project |
| [**internalV1ProjectsIdGet**](InternalApiProjectsApi.md#internalV1ProjectsIdGet) | **GET** /internal/v1/projects/{id} | Get project |
| [**internalV1ProjectsIdLimitsPatch**](InternalApiProjectsApi.md#internalV1ProjectsIdLimitsPatch) | **PATCH** /internal/v1/projects/{id}/limits | Update project limits |
| [**internalV1ProjectsIdPut**](InternalApiProjectsApi.md#internalV1ProjectsIdPut) | **PUT** /internal/v1/projects/{id} | Update project |
| [**internalV1ProjectsKeyKeyNameGet**](InternalApiProjectsApi.md#internalV1ProjectsKeyKeyNameGet) | **GET** /internal/v1/projects/key/{keyName} | Get project by key name |
| [**internalV1ProjectsPost**](InternalApiProjectsApi.md#internalV1ProjectsPost) | **POST** /internal/v1/projects | Create project |


<a id="internalV1ProjectsCompanyCompanyIdGet"></a>
# **internalV1ProjectsCompanyCompanyIdGet**
> ProjectListResponse internalV1ProjectsCompanyCompanyIdGet(companyId)

Get projects by company

Get projects by company ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiProjectsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiProjectsApi apiInstance = new InternalApiProjectsApi(defaultClient);
    UUID companyId = UUID.randomUUID(); // UUID | 
    try {
      ProjectListResponse result = apiInstance.internalV1ProjectsCompanyCompanyIdGet(companyId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiProjectsApi#internalV1ProjectsCompanyCompanyIdGet");
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

[**ProjectListResponse**](ProjectListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of projects in company |  -  |
| **404** | Not found |  -  |

<a id="internalV1ProjectsGet"></a>
# **internalV1ProjectsGet**
> ProjectListResponse internalV1ProjectsGet(limit, offset, search)

List projects

Get a list of projects with pagination

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiProjectsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiProjectsApi apiInstance = new InternalApiProjectsApi(defaultClient);
    Integer limit = 10; // Integer | Number of items to return
    Integer offset = 0; // Integer | Number of items to skip
    String search = "search_example"; // String | Search term
    try {
      ProjectListResponse result = apiInstance.internalV1ProjectsGet(limit, offset, search);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiProjectsApi#internalV1ProjectsGet");
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

[**ProjectListResponse**](ProjectListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of projects |  -  |

<a id="internalV1ProjectsIdDelete"></a>
# **internalV1ProjectsIdDelete**
> internalV1ProjectsIdDelete(id)

Delete project

Delete project by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiProjectsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiProjectsApi apiInstance = new InternalApiProjectsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      apiInstance.internalV1ProjectsIdDelete(id);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiProjectsApi#internalV1ProjectsIdDelete");
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

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Project deleted successfully |  -  |
| **404** | Not found |  -  |

<a id="internalV1ProjectsIdGet"></a>
# **internalV1ProjectsIdGet**
> ProjectResponse internalV1ProjectsIdGet(id)

Get project

Get project by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiProjectsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiProjectsApi apiInstance = new InternalApiProjectsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      ProjectResponse result = apiInstance.internalV1ProjectsIdGet(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiProjectsApi#internalV1ProjectsIdGet");
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

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project details |  -  |
| **404** | Not found |  -  |

<a id="internalV1ProjectsIdLimitsPatch"></a>
# **internalV1ProjectsIdLimitsPatch**
> ProjectResponse internalV1ProjectsIdLimitsPatch(id, projectLimitsUpdateRequest)

Update project limits

Update project limits by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiProjectsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiProjectsApi apiInstance = new InternalApiProjectsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    ProjectLimitsUpdateRequest projectLimitsUpdateRequest = new ProjectLimitsUpdateRequest(); // ProjectLimitsUpdateRequest | 
    try {
      ProjectResponse result = apiInstance.internalV1ProjectsIdLimitsPatch(id, projectLimitsUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiProjectsApi#internalV1ProjectsIdLimitsPatch");
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
| **projectLimitsUpdateRequest** | [**ProjectLimitsUpdateRequest**](ProjectLimitsUpdateRequest.md)|  | |

### Return type

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project limits updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1ProjectsIdPut"></a>
# **internalV1ProjectsIdPut**
> ProjectResponse internalV1ProjectsIdPut(id, projectUpdateRequest)

Update project

Update project by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiProjectsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiProjectsApi apiInstance = new InternalApiProjectsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    ProjectUpdateRequest projectUpdateRequest = new ProjectUpdateRequest(); // ProjectUpdateRequest | 
    try {
      ProjectResponse result = apiInstance.internalV1ProjectsIdPut(id, projectUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiProjectsApi#internalV1ProjectsIdPut");
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
| **projectUpdateRequest** | [**ProjectUpdateRequest**](ProjectUpdateRequest.md)|  | |

### Return type

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1ProjectsKeyKeyNameGet"></a>
# **internalV1ProjectsKeyKeyNameGet**
> ProjectResponse internalV1ProjectsKeyKeyNameGet(keyName)

Get project by key name

Get project by notify key name

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiProjectsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiProjectsApi apiInstance = new InternalApiProjectsApi(defaultClient);
    String keyName = "keyName_example"; // String | 
    try {
      ProjectResponse result = apiInstance.internalV1ProjectsKeyKeyNameGet(keyName);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiProjectsApi#internalV1ProjectsKeyKeyNameGet");
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
| **keyName** | **String**|  | |

### Return type

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project details |  -  |
| **404** | Not found |  -  |

<a id="internalV1ProjectsPost"></a>
# **internalV1ProjectsPost**
> ProjectResponse internalV1ProjectsPost(projectCreateRequest)

Create project

Create a new project

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiProjectsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiProjectsApi apiInstance = new InternalApiProjectsApi(defaultClient);
    ProjectCreateRequest projectCreateRequest = new ProjectCreateRequest(); // ProjectCreateRequest | 
    try {
      ProjectResponse result = apiInstance.internalV1ProjectsPost(projectCreateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiProjectsApi#internalV1ProjectsPost");
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
| **projectCreateRequest** | [**ProjectCreateRequest**](ProjectCreateRequest.md)|  | |

### Return type

[**ProjectResponse**](ProjectResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Project created successfully |  -  |
| **400** | Bad request |  -  |

