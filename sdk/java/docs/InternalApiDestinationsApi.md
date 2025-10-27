# InternalApiDestinationsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1DestinationsBotBotIdGet**](InternalApiDestinationsApi.md#internalV1DestinationsBotBotIdGet) | **GET** /internal/v1/destinations/bot/{botId} | Get destinations by bot |
| [**internalV1DestinationsGet**](InternalApiDestinationsApi.md#internalV1DestinationsGet) | **GET** /internal/v1/destinations | List destinations |
| [**internalV1DestinationsIdDelete**](InternalApiDestinationsApi.md#internalV1DestinationsIdDelete) | **DELETE** /internal/v1/destinations/{id} | Delete destination |
| [**internalV1DestinationsIdGet**](InternalApiDestinationsApi.md#internalV1DestinationsIdGet) | **GET** /internal/v1/destinations/{id} | Get destination |
| [**internalV1DestinationsIdPut**](InternalApiDestinationsApi.md#internalV1DestinationsIdPut) | **PUT** /internal/v1/destinations/{id} | Update destination |
| [**internalV1DestinationsIdTargetsPatch**](InternalApiDestinationsApi.md#internalV1DestinationsIdTargetsPatch) | **PATCH** /internal/v1/destinations/{id}/targets | Update destination targets |
| [**internalV1DestinationsIdValidatePost**](InternalApiDestinationsApi.md#internalV1DestinationsIdValidatePost) | **POST** /internal/v1/destinations/{id}/validate | Validate destination targets |
| [**internalV1DestinationsPost**](InternalApiDestinationsApi.md#internalV1DestinationsPost) | **POST** /internal/v1/destinations | Create destination |
| [**internalV1DestinationsProjectProjectIdGet**](InternalApiDestinationsApi.md#internalV1DestinationsProjectProjectIdGet) | **GET** /internal/v1/destinations/project/{projectId} | Get destinations by project |
| [**internalV1DestinationsSearchGet**](InternalApiDestinationsApi.md#internalV1DestinationsSearchGet) | **GET** /internal/v1/destinations/search | Search destinations |


<a id="internalV1DestinationsBotBotIdGet"></a>
# **internalV1DestinationsBotBotIdGet**
> DestinationListResponse internalV1DestinationsBotBotIdGet(botId)

Get destinations by bot

Get destinations by bot ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    UUID botId = UUID.randomUUID(); // UUID | 
    try {
      DestinationListResponse result = apiInstance.internalV1DestinationsBotBotIdGet(botId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsBotBotIdGet");
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
| **botId** | **UUID**|  | |

### Return type

[**DestinationListResponse**](DestinationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of destinations for bot |  -  |
| **404** | Not found |  -  |

<a id="internalV1DestinationsGet"></a>
# **internalV1DestinationsGet**
> DestinationListResponse internalV1DestinationsGet(limit, offset, search)

List destinations

Get a list of destinations with pagination

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    Integer limit = 10; // Integer | Number of items to return
    Integer offset = 0; // Integer | Number of items to skip
    String search = "search_example"; // String | Search term
    try {
      DestinationListResponse result = apiInstance.internalV1DestinationsGet(limit, offset, search);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsGet");
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

[**DestinationListResponse**](DestinationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of destinations |  -  |

<a id="internalV1DestinationsIdDelete"></a>
# **internalV1DestinationsIdDelete**
> internalV1DestinationsIdDelete(id)

Delete destination

Delete destination by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      apiInstance.internalV1DestinationsIdDelete(id);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsIdDelete");
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
| **204** | Destination deleted successfully |  -  |
| **404** | Not found |  -  |

<a id="internalV1DestinationsIdGet"></a>
# **internalV1DestinationsIdGet**
> DestinationResponse internalV1DestinationsIdGet(id)

Get destination

Get destination by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      DestinationResponse result = apiInstance.internalV1DestinationsIdGet(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsIdGet");
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

[**DestinationResponse**](DestinationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Destination details |  -  |
| **404** | Not found |  -  |

<a id="internalV1DestinationsIdPut"></a>
# **internalV1DestinationsIdPut**
> DestinationResponse internalV1DestinationsIdPut(id, destinationUpdateRequest)

Update destination

Update destination by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    DestinationUpdateRequest destinationUpdateRequest = new DestinationUpdateRequest(); // DestinationUpdateRequest | 
    try {
      DestinationResponse result = apiInstance.internalV1DestinationsIdPut(id, destinationUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsIdPut");
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
| **destinationUpdateRequest** | [**DestinationUpdateRequest**](DestinationUpdateRequest.md)|  | |

### Return type

[**DestinationResponse**](DestinationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Destination updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1DestinationsIdTargetsPatch"></a>
# **internalV1DestinationsIdTargetsPatch**
> DestinationResponse internalV1DestinationsIdTargetsPatch(id, destinationTargetsUpdateRequest)

Update destination targets

Update destination targets by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    DestinationTargetsUpdateRequest destinationTargetsUpdateRequest = new DestinationTargetsUpdateRequest(); // DestinationTargetsUpdateRequest | 
    try {
      DestinationResponse result = apiInstance.internalV1DestinationsIdTargetsPatch(id, destinationTargetsUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsIdTargetsPatch");
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
| **destinationTargetsUpdateRequest** | [**DestinationTargetsUpdateRequest**](DestinationTargetsUpdateRequest.md)|  | |

### Return type

[**DestinationResponse**](DestinationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Destination targets updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1DestinationsIdValidatePost"></a>
# **internalV1DestinationsIdValidatePost**
> DestinationValidationResponse internalV1DestinationsIdValidatePost(id)

Validate destination targets

Validate destination targets by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      DestinationValidationResponse result = apiInstance.internalV1DestinationsIdValidatePost(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsIdValidatePost");
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

[**DestinationValidationResponse**](DestinationValidationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Destination targets validation result |  -  |
| **404** | Not found |  -  |

<a id="internalV1DestinationsPost"></a>
# **internalV1DestinationsPost**
> DestinationResponse internalV1DestinationsPost(destinationCreateRequest)

Create destination

Create a new destination

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    DestinationCreateRequest destinationCreateRequest = new DestinationCreateRequest(); // DestinationCreateRequest | 
    try {
      DestinationResponse result = apiInstance.internalV1DestinationsPost(destinationCreateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsPost");
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
| **destinationCreateRequest** | [**DestinationCreateRequest**](DestinationCreateRequest.md)|  | |

### Return type

[**DestinationResponse**](DestinationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Destination created successfully |  -  |
| **400** | Bad request |  -  |

<a id="internalV1DestinationsProjectProjectIdGet"></a>
# **internalV1DestinationsProjectProjectIdGet**
> DestinationListResponse internalV1DestinationsProjectProjectIdGet(projectId)

Get destinations by project

Get destinations by project ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    UUID projectId = UUID.randomUUID(); // UUID | 
    try {
      DestinationListResponse result = apiInstance.internalV1DestinationsProjectProjectIdGet(projectId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsProjectProjectIdGet");
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

[**DestinationListResponse**](DestinationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of destinations for project |  -  |
| **404** | Not found |  -  |

<a id="internalV1DestinationsSearchGet"></a>
# **internalV1DestinationsSearchGet**
> DestinationListResponse internalV1DestinationsSearchGet(search, type, status)

Search destinations

Search destinations with filters

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiDestinationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiDestinationsApi apiInstance = new InternalApiDestinationsApi(defaultClient);
    String search = "search_example"; // String | Search term
    String type = "personal"; // String | 
    String status = "active"; // String | 
    try {
      DestinationListResponse result = apiInstance.internalV1DestinationsSearchGet(search, type, status);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiDestinationsApi#internalV1DestinationsSearchGet");
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
| **search** | **String**| Search term | [optional] |
| **type** | **String**|  | [optional] [enum: personal, group, channel] |
| **status** | **String**|  | [optional] [enum: active, inactive] |

### Return type

[**DestinationListResponse**](DestinationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Search results |  -  |

