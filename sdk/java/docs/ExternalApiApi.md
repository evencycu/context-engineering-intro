# ExternalApiApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**apiV1DestinationsNotifyKeyGet**](ExternalApiApi.md#apiV1DestinationsNotifyKeyGet) | **GET** /api/v1/destinations/{notifyKey} | Get project destinations (External API) |
| [**apiV1MessagesPost**](ExternalApiApi.md#apiV1MessagesPost) | **POST** /api/v1/messages | Handle Bot Framework messages |
| [**apiV1MessagesTestPost**](ExternalApiApi.md#apiV1MessagesTestPost) | **POST** /api/v1/messages/test | Test proactive message |
| [**apiV1NotifyPost**](ExternalApiApi.md#apiV1NotifyPost) | **POST** /api/v1/notify | Send notification (External API) |
| [**apiV1ProvisionNotifyKeyDisablePost**](ExternalApiApi.md#apiV1ProvisionNotifyKeyDisablePost) | **POST** /api/v1/provision/{notifyKey}/disable | Disable project provision |
| [**apiV1ProvisionNotifyKeyEnablePost**](ExternalApiApi.md#apiV1ProvisionNotifyKeyEnablePost) | **POST** /api/v1/provision/{notifyKey}/enable | Enable project provision |
| [**apiV1ProvisionNotifyKeyGet**](ExternalApiApi.md#apiV1ProvisionNotifyKeyGet) | **GET** /api/v1/provision/{notifyKey} | Get project provision |
| [**apiV1ProvisionNotifyKeyPut**](ExternalApiApi.md#apiV1ProvisionNotifyKeyPut) | **PUT** /api/v1/provision/{notifyKey} | Update project provision |
| [**apiV1ProvisionPost**](ExternalApiApi.md#apiV1ProvisionPost) | **POST** /api/v1/provision | Create project provision |


<a id="apiV1DestinationsNotifyKeyGet"></a>
# **apiV1DestinationsNotifyKeyGet**
> ProjectDestinationsResponse apiV1DestinationsNotifyKeyGet(notifyKey)

Get project destinations (External API)

Get destinations for a project using notify key

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.ExternalApiApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    ExternalApiApi apiInstance = new ExternalApiApi(defaultClient);
    String notifyKey = "notifyKey_example"; // String | Project notify key
    try {
      ProjectDestinationsResponse result = apiInstance.apiV1DestinationsNotifyKeyGet(notifyKey);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling ExternalApiApi#apiV1DestinationsNotifyKeyGet");
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
| **notifyKey** | **String**| Project notify key | |

### Return type

[**ProjectDestinationsResponse**](ProjectDestinationsResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project destinations |  -  |
| **404** | Not found |  -  |

<a id="apiV1MessagesPost"></a>
# **apiV1MessagesPost**
> MessageResponse apiV1MessagesPost(botFrameworkActivity)

Handle Bot Framework messages

Handle incoming Bot Framework Activity messages

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.ExternalApiApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    ExternalApiApi apiInstance = new ExternalApiApi(defaultClient);
    BotFrameworkActivity botFrameworkActivity = new BotFrameworkActivity(); // BotFrameworkActivity | 
    try {
      MessageResponse result = apiInstance.apiV1MessagesPost(botFrameworkActivity);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling ExternalApiApi#apiV1MessagesPost");
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
| **botFrameworkActivity** | [**BotFrameworkActivity**](BotFrameworkActivity.md)|  | |

### Return type

[**MessageResponse**](MessageResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Message processed successfully |  -  |

<a id="apiV1MessagesTestPost"></a>
# **apiV1MessagesTestPost**
> MessageResponse apiV1MessagesTestPost(proactiveTestRequest)

Test proactive message

Send a test proactive message

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.ExternalApiApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    ExternalApiApi apiInstance = new ExternalApiApi(defaultClient);
    ProactiveTestRequest proactiveTestRequest = new ProactiveTestRequest(); // ProactiveTestRequest | 
    try {
      MessageResponse result = apiInstance.apiV1MessagesTestPost(proactiveTestRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling ExternalApiApi#apiV1MessagesTestPost");
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
| **proactiveTestRequest** | [**ProactiveTestRequest**](ProactiveTestRequest.md)|  | |

### Return type

[**MessageResponse**](MessageResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Test message sent successfully |  -  |

<a id="apiV1NotifyPost"></a>
# **apiV1NotifyPost**
> ExternalNotifyResponse apiV1NotifyPost(externalNotifyRequest)

Send notification (External API)

Send a notification to Teams using project notify key

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.ExternalApiApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    ExternalApiApi apiInstance = new ExternalApiApi(defaultClient);
    ExternalNotifyRequest externalNotifyRequest = new ExternalNotifyRequest(); // ExternalNotifyRequest | 
    try {
      ExternalNotifyResponse result = apiInstance.apiV1NotifyPost(externalNotifyRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling ExternalApiApi#apiV1NotifyPost");
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
| **externalNotifyRequest** | [**ExternalNotifyRequest**](ExternalNotifyRequest.md)|  | |

### Return type

[**ExternalNotifyResponse**](ExternalNotifyResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Notification sent successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="apiV1ProvisionNotifyKeyDisablePost"></a>
# **apiV1ProvisionNotifyKeyDisablePost**
> ProvisionResponse apiV1ProvisionNotifyKeyDisablePost(notifyKey)

Disable project provision

Disable project provision by notify key

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.ExternalApiApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    ExternalApiApi apiInstance = new ExternalApiApi(defaultClient);
    String notifyKey = "notifyKey_example"; // String | Project notify key
    try {
      ProvisionResponse result = apiInstance.apiV1ProvisionNotifyKeyDisablePost(notifyKey);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling ExternalApiApi#apiV1ProvisionNotifyKeyDisablePost");
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
| **notifyKey** | **String**| Project notify key | |

### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project provision disabled successfully |  -  |
| **404** | Not found |  -  |

<a id="apiV1ProvisionNotifyKeyEnablePost"></a>
# **apiV1ProvisionNotifyKeyEnablePost**
> ProvisionResponse apiV1ProvisionNotifyKeyEnablePost(notifyKey)

Enable project provision

Enable project provision by notify key

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.ExternalApiApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    ExternalApiApi apiInstance = new ExternalApiApi(defaultClient);
    String notifyKey = "notifyKey_example"; // String | Project notify key
    try {
      ProvisionResponse result = apiInstance.apiV1ProvisionNotifyKeyEnablePost(notifyKey);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling ExternalApiApi#apiV1ProvisionNotifyKeyEnablePost");
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
| **notifyKey** | **String**| Project notify key | |

### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project provision enabled successfully |  -  |
| **404** | Not found |  -  |

<a id="apiV1ProvisionNotifyKeyGet"></a>
# **apiV1ProvisionNotifyKeyGet**
> ProvisionResponse apiV1ProvisionNotifyKeyGet(notifyKey)

Get project provision

Get project provision by notify key

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.ExternalApiApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    ExternalApiApi apiInstance = new ExternalApiApi(defaultClient);
    String notifyKey = "notifyKey_example"; // String | Project notify key
    try {
      ProvisionResponse result = apiInstance.apiV1ProvisionNotifyKeyGet(notifyKey);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling ExternalApiApi#apiV1ProvisionNotifyKeyGet");
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
| **notifyKey** | **String**| Project notify key | |

### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project provision details |  -  |
| **404** | Not found |  -  |

<a id="apiV1ProvisionNotifyKeyPut"></a>
# **apiV1ProvisionNotifyKeyPut**
> ProvisionResponse apiV1ProvisionNotifyKeyPut(notifyKey, provisionUpdateRequest)

Update project provision

Update project provision by notify key

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.ExternalApiApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    ExternalApiApi apiInstance = new ExternalApiApi(defaultClient);
    String notifyKey = "notifyKey_example"; // String | Project notify key
    ProvisionUpdateRequest provisionUpdateRequest = new ProvisionUpdateRequest(); // ProvisionUpdateRequest | 
    try {
      ProvisionResponse result = apiInstance.apiV1ProvisionNotifyKeyPut(notifyKey, provisionUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling ExternalApiApi#apiV1ProvisionNotifyKeyPut");
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
| **notifyKey** | **String**| Project notify key | |
| **provisionUpdateRequest** | [**ProvisionUpdateRequest**](ProvisionUpdateRequest.md)|  | |

### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Project provision updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="apiV1ProvisionPost"></a>
# **apiV1ProvisionPost**
> ProvisionResponse apiV1ProvisionPost(provisionCreateRequest)

Create project provision

Create a new project provision

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.ExternalApiApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    ExternalApiApi apiInstance = new ExternalApiApi(defaultClient);
    ProvisionCreateRequest provisionCreateRequest = new ProvisionCreateRequest(); // ProvisionCreateRequest | 
    try {
      ProvisionResponse result = apiInstance.apiV1ProvisionPost(provisionCreateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling ExternalApiApi#apiV1ProvisionPost");
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
| **provisionCreateRequest** | [**ProvisionCreateRequest**](ProvisionCreateRequest.md)|  | |

### Return type

[**ProvisionResponse**](ProvisionResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Project provision created successfully |  -  |
| **400** | Bad request |  -  |

