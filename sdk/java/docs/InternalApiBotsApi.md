# InternalApiBotsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1BotsPlatformGet**](InternalApiBotsApi.md#internalV1BotsPlatformGet) | **GET** /internal/v1/bots/platform | List Teams bots |
| [**internalV1BotsPlatformIdCapabilitiesPatch**](InternalApiBotsApi.md#internalV1BotsPlatformIdCapabilitiesPatch) | **PATCH** /internal/v1/bots/platform/{id}/capabilities | Update bot capabilities |
| [**internalV1BotsPlatformIdDelete**](InternalApiBotsApi.md#internalV1BotsPlatformIdDelete) | **DELETE** /internal/v1/bots/platform/{id} | Delete Teams bot |
| [**internalV1BotsPlatformIdGet**](InternalApiBotsApi.md#internalV1BotsPlatformIdGet) | **GET** /internal/v1/bots/platform/{id} | Get Teams bot |
| [**internalV1BotsPlatformIdPut**](InternalApiBotsApi.md#internalV1BotsPlatformIdPut) | **PUT** /internal/v1/bots/platform/{id} | Update Teams bot |
| [**internalV1BotsPlatformIdStatusPatch**](InternalApiBotsApi.md#internalV1BotsPlatformIdStatusPatch) | **PATCH** /internal/v1/bots/platform/{id}/status | Update bot status |
| [**internalV1BotsPlatformIdTestPost**](InternalApiBotsApi.md#internalV1BotsPlatformIdTestPost) | **POST** /internal/v1/bots/platform/{id}/test | Test bot connection |
| [**internalV1BotsPlatformPost**](InternalApiBotsApi.md#internalV1BotsPlatformPost) | **POST** /internal/v1/bots/platform | Create Teams bot |
| [**internalV1BotsStatusStatusGet**](InternalApiBotsApi.md#internalV1BotsStatusStatusGet) | **GET** /internal/v1/bots/status/{status} | Get bots by status |


<a id="internalV1BotsPlatformGet"></a>
# **internalV1BotsPlatformGet**
> BotListResponse internalV1BotsPlatformGet(limit, offset, search)

List Teams bots

Get a list of Teams bot services

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBotsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBotsApi apiInstance = new InternalApiBotsApi(defaultClient);
    Integer limit = 10; // Integer | Number of items to return
    Integer offset = 0; // Integer | Number of items to skip
    String search = "search_example"; // String | Search term
    try {
      BotListResponse result = apiInstance.internalV1BotsPlatformGet(limit, offset, search);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBotsApi#internalV1BotsPlatformGet");
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

[**BotListResponse**](BotListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of Teams bots |  -  |

<a id="internalV1BotsPlatformIdCapabilitiesPatch"></a>
# **internalV1BotsPlatformIdCapabilitiesPatch**
> BotResponse internalV1BotsPlatformIdCapabilitiesPatch(id, botCapabilitiesUpdateRequest)

Update bot capabilities

Update bot capabilities by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBotsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBotsApi apiInstance = new InternalApiBotsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    BotCapabilitiesUpdateRequest botCapabilitiesUpdateRequest = new BotCapabilitiesUpdateRequest(); // BotCapabilitiesUpdateRequest | 
    try {
      BotResponse result = apiInstance.internalV1BotsPlatformIdCapabilitiesPatch(id, botCapabilitiesUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBotsApi#internalV1BotsPlatformIdCapabilitiesPatch");
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
| **botCapabilitiesUpdateRequest** | [**BotCapabilitiesUpdateRequest**](BotCapabilitiesUpdateRequest.md)|  | |

### Return type

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Bot capabilities updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1BotsPlatformIdDelete"></a>
# **internalV1BotsPlatformIdDelete**
> internalV1BotsPlatformIdDelete(id)

Delete Teams bot

Delete Teams bot by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBotsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBotsApi apiInstance = new InternalApiBotsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      apiInstance.internalV1BotsPlatformIdDelete(id);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBotsApi#internalV1BotsPlatformIdDelete");
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
| **204** | Bot deleted successfully |  -  |
| **404** | Not found |  -  |

<a id="internalV1BotsPlatformIdGet"></a>
# **internalV1BotsPlatformIdGet**
> BotResponse internalV1BotsPlatformIdGet(id)

Get Teams bot

Get Teams bot by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBotsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBotsApi apiInstance = new InternalApiBotsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      BotResponse result = apiInstance.internalV1BotsPlatformIdGet(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBotsApi#internalV1BotsPlatformIdGet");
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

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Bot details |  -  |
| **404** | Not found |  -  |

<a id="internalV1BotsPlatformIdPut"></a>
# **internalV1BotsPlatformIdPut**
> BotResponse internalV1BotsPlatformIdPut(id, botUpdateRequest)

Update Teams bot

Update Teams bot by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBotsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBotsApi apiInstance = new InternalApiBotsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    BotUpdateRequest botUpdateRequest = new BotUpdateRequest(); // BotUpdateRequest | 
    try {
      BotResponse result = apiInstance.internalV1BotsPlatformIdPut(id, botUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBotsApi#internalV1BotsPlatformIdPut");
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
| **botUpdateRequest** | [**BotUpdateRequest**](BotUpdateRequest.md)|  | |

### Return type

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Bot updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1BotsPlatformIdStatusPatch"></a>
# **internalV1BotsPlatformIdStatusPatch**
> BotResponse internalV1BotsPlatformIdStatusPatch(id, botStatusUpdateRequest)

Update bot status

Update bot status by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBotsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBotsApi apiInstance = new InternalApiBotsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    BotStatusUpdateRequest botStatusUpdateRequest = new BotStatusUpdateRequest(); // BotStatusUpdateRequest | 
    try {
      BotResponse result = apiInstance.internalV1BotsPlatformIdStatusPatch(id, botStatusUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBotsApi#internalV1BotsPlatformIdStatusPatch");
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
| **botStatusUpdateRequest** | [**BotStatusUpdateRequest**](BotStatusUpdateRequest.md)|  | |

### Return type

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Bot status updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1BotsPlatformIdTestPost"></a>
# **internalV1BotsPlatformIdTestPost**
> BotTestResponse internalV1BotsPlatformIdTestPost(id)

Test bot connection

Test bot connection by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBotsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBotsApi apiInstance = new InternalApiBotsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      BotTestResponse result = apiInstance.internalV1BotsPlatformIdTestPost(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBotsApi#internalV1BotsPlatformIdTestPost");
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

[**BotTestResponse**](BotTestResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Bot connection test result |  -  |
| **404** | Not found |  -  |

<a id="internalV1BotsPlatformPost"></a>
# **internalV1BotsPlatformPost**
> BotResponse internalV1BotsPlatformPost(botCreateRequest)

Create Teams bot

Create a new Teams bot service

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBotsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBotsApi apiInstance = new InternalApiBotsApi(defaultClient);
    BotCreateRequest botCreateRequest = new BotCreateRequest(); // BotCreateRequest | 
    try {
      BotResponse result = apiInstance.internalV1BotsPlatformPost(botCreateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBotsApi#internalV1BotsPlatformPost");
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
| **botCreateRequest** | [**BotCreateRequest**](BotCreateRequest.md)|  | |

### Return type

[**BotResponse**](BotResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Bot created successfully |  -  |
| **400** | Bad request |  -  |

<a id="internalV1BotsStatusStatusGet"></a>
# **internalV1BotsStatusStatusGet**
> BotListResponse internalV1BotsStatusStatusGet(status)

Get bots by status

Get bots by status

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiBotsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiBotsApi apiInstance = new InternalApiBotsApi(defaultClient);
    String status = "active"; // String | 
    try {
      BotListResponse result = apiInstance.internalV1BotsStatusStatusGet(status);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiBotsApi#internalV1BotsStatusStatusGet");
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
| **status** | **String**|  | [enum: active, inactive, error] |

### Return type

[**BotListResponse**](BotListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of bots with status |  -  |

