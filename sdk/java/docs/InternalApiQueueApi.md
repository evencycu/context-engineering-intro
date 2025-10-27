# InternalApiQueueApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1QueueClearPost**](InternalApiQueueApi.md#internalV1QueueClearPost) | **POST** /internal/v1/queue/clear | Clear queue |
| [**internalV1QueueStatsGet**](InternalApiQueueApi.md#internalV1QueueStatsGet) | **GET** /internal/v1/queue/stats | Get queue statistics |
| [**internalV1QueueStatusGet**](InternalApiQueueApi.md#internalV1QueueStatusGet) | **GET** /internal/v1/queue/status | Get queue status |


<a id="internalV1QueueClearPost"></a>
# **internalV1QueueClearPost**
> QueueClearResponse internalV1QueueClearPost()

Clear queue

Clear all items from the queue

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiQueueApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiQueueApi apiInstance = new InternalApiQueueApi(defaultClient);
    try {
      QueueClearResponse result = apiInstance.internalV1QueueClearPost();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiQueueApi#internalV1QueueClearPost");
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

[**QueueClearResponse**](QueueClearResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Queue cleared successfully |  -  |

<a id="internalV1QueueStatsGet"></a>
# **internalV1QueueStatsGet**
> QueueStatsResponse internalV1QueueStatsGet()

Get queue statistics

Get queue statistics and status

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiQueueApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiQueueApi apiInstance = new InternalApiQueueApi(defaultClient);
    try {
      QueueStatsResponse result = apiInstance.internalV1QueueStatsGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiQueueApi#internalV1QueueStatsGet");
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

[**QueueStatsResponse**](QueueStatsResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Queue statistics |  -  |

<a id="internalV1QueueStatusGet"></a>
# **internalV1QueueStatusGet**
> QueueStatusResponse internalV1QueueStatusGet()

Get queue status

Get current queue status

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiQueueApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiQueueApi apiInstance = new InternalApiQueueApi(defaultClient);
    try {
      QueueStatusResponse result = apiInstance.internalV1QueueStatusGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiQueueApi#internalV1QueueStatusGet");
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

[**QueueStatusResponse**](QueueStatusResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Queue status |  -  |

