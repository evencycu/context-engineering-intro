# InternalApiNotificationsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1NotificationsDateRangeGet**](InternalApiNotificationsApi.md#internalV1NotificationsDateRangeGet) | **GET** /internal/v1/notifications/date-range | Get notifications by date range |
| [**internalV1NotificationsGet**](InternalApiNotificationsApi.md#internalV1NotificationsGet) | **GET** /internal/v1/notifications | List notifications |
| [**internalV1NotificationsIdDelete**](InternalApiNotificationsApi.md#internalV1NotificationsIdDelete) | **DELETE** /internal/v1/notifications/{id} | Cancel notification |
| [**internalV1NotificationsIdGet**](InternalApiNotificationsApi.md#internalV1NotificationsIdGet) | **GET** /internal/v1/notifications/{id} | Get notification |
| [**internalV1NotificationsIdRetryPost**](InternalApiNotificationsApi.md#internalV1NotificationsIdRetryPost) | **POST** /internal/v1/notifications/{id}/retry | Retry notification |
| [**internalV1NotificationsPost**](InternalApiNotificationsApi.md#internalV1NotificationsPost) | **POST** /internal/v1/notifications | Send notification |
| [**internalV1NotificationsProjectProjectIdGet**](InternalApiNotificationsApi.md#internalV1NotificationsProjectProjectIdGet) | **GET** /internal/v1/notifications/project/{projectId} | Get notifications by project |
| [**internalV1NotificationsSenderSenderIdGet**](InternalApiNotificationsApi.md#internalV1NotificationsSenderSenderIdGet) | **GET** /internal/v1/notifications/sender/{senderId} | Get notifications by sender |
| [**internalV1NotificationsStatusStatusGet**](InternalApiNotificationsApi.md#internalV1NotificationsStatusStatusGet) | **GET** /internal/v1/notifications/status/{status} | Get notifications by status |


<a id="internalV1NotificationsDateRangeGet"></a>
# **internalV1NotificationsDateRangeGet**
> NotificationDateRangeResponse internalV1NotificationsDateRangeGet(startDate, endDate)

Get notifications by date range

Get notifications within date range

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiNotificationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiNotificationsApi apiInstance = new InternalApiNotificationsApi(defaultClient);
    OffsetDateTime startDate = OffsetDateTime.now(); // OffsetDateTime | 
    OffsetDateTime endDate = OffsetDateTime.now(); // OffsetDateTime | 
    try {
      NotificationDateRangeResponse result = apiInstance.internalV1NotificationsDateRangeGet(startDate, endDate);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiNotificationsApi#internalV1NotificationsDateRangeGet");
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
| **startDate** | **OffsetDateTime**|  | |
| **endDate** | **OffsetDateTime**|  | |

### Return type

[**NotificationDateRangeResponse**](NotificationDateRangeResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of notifications in date range |  -  |

<a id="internalV1NotificationsGet"></a>
# **internalV1NotificationsGet**
> NotificationListResponse internalV1NotificationsGet(limit, offset, search, status)

List notifications

Get a list of notifications with pagination

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiNotificationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiNotificationsApi apiInstance = new InternalApiNotificationsApi(defaultClient);
    Integer limit = 10; // Integer | Number of items to return
    Integer offset = 0; // Integer | Number of items to skip
    String search = "search_example"; // String | Search term
    String status = "pending"; // String | 
    try {
      NotificationListResponse result = apiInstance.internalV1NotificationsGet(limit, offset, search, status);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiNotificationsApi#internalV1NotificationsGet");
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
| **status** | **String**|  | [optional] [enum: pending, sent, failed, cancelled] |

### Return type

[**NotificationListResponse**](NotificationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of notifications |  -  |

<a id="internalV1NotificationsIdDelete"></a>
# **internalV1NotificationsIdDelete**
> NotificationCancelResponse internalV1NotificationsIdDelete(id)

Cancel notification

Cancel notification by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiNotificationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiNotificationsApi apiInstance = new InternalApiNotificationsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      NotificationCancelResponse result = apiInstance.internalV1NotificationsIdDelete(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiNotificationsApi#internalV1NotificationsIdDelete");
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

[**NotificationCancelResponse**](NotificationCancelResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Notification cancelled successfully |  -  |
| **404** | Not found |  -  |

<a id="internalV1NotificationsIdGet"></a>
# **internalV1NotificationsIdGet**
> NotificationResponse internalV1NotificationsIdGet(id)

Get notification

Get notification by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiNotificationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiNotificationsApi apiInstance = new InternalApiNotificationsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      NotificationResponse result = apiInstance.internalV1NotificationsIdGet(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiNotificationsApi#internalV1NotificationsIdGet");
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

[**NotificationResponse**](NotificationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Notification details |  -  |
| **404** | Not found |  -  |

<a id="internalV1NotificationsIdRetryPost"></a>
# **internalV1NotificationsIdRetryPost**
> NotificationRetryResponse internalV1NotificationsIdRetryPost(id)

Retry notification

Retry failed notification by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiNotificationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiNotificationsApi apiInstance = new InternalApiNotificationsApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      NotificationRetryResponse result = apiInstance.internalV1NotificationsIdRetryPost(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiNotificationsApi#internalV1NotificationsIdRetryPost");
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

[**NotificationRetryResponse**](NotificationRetryResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Notification retry initiated |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1NotificationsPost"></a>
# **internalV1NotificationsPost**
> NotificationSendResponse internalV1NotificationsPost(notificationSendRequest)

Send notification

Send a new notification

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiNotificationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiNotificationsApi apiInstance = new InternalApiNotificationsApi(defaultClient);
    NotificationSendRequest notificationSendRequest = new NotificationSendRequest(); // NotificationSendRequest | 
    try {
      NotificationSendResponse result = apiInstance.internalV1NotificationsPost(notificationSendRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiNotificationsApi#internalV1NotificationsPost");
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
| **notificationSendRequest** | [**NotificationSendRequest**](NotificationSendRequest.md)|  | |

### Return type

[**NotificationSendResponse**](NotificationSendResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **202** | Notification queued successfully |  -  |
| **400** | Bad request |  -  |

<a id="internalV1NotificationsProjectProjectIdGet"></a>
# **internalV1NotificationsProjectProjectIdGet**
> NotificationListResponse internalV1NotificationsProjectProjectIdGet(projectId)

Get notifications by project

Get notifications by project ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiNotificationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiNotificationsApi apiInstance = new InternalApiNotificationsApi(defaultClient);
    UUID projectId = UUID.randomUUID(); // UUID | 
    try {
      NotificationListResponse result = apiInstance.internalV1NotificationsProjectProjectIdGet(projectId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiNotificationsApi#internalV1NotificationsProjectProjectIdGet");
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

[**NotificationListResponse**](NotificationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of notifications for project |  -  |
| **404** | Not found |  -  |

<a id="internalV1NotificationsSenderSenderIdGet"></a>
# **internalV1NotificationsSenderSenderIdGet**
> NotificationListResponse internalV1NotificationsSenderSenderIdGet(senderId)

Get notifications by sender

Get notifications by sender ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiNotificationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiNotificationsApi apiInstance = new InternalApiNotificationsApi(defaultClient);
    UUID senderId = UUID.randomUUID(); // UUID | 
    try {
      NotificationListResponse result = apiInstance.internalV1NotificationsSenderSenderIdGet(senderId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiNotificationsApi#internalV1NotificationsSenderSenderIdGet");
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
| **senderId** | **UUID**|  | |

### Return type

[**NotificationListResponse**](NotificationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of notifications by sender |  -  |
| **404** | Not found |  -  |

<a id="internalV1NotificationsStatusStatusGet"></a>
# **internalV1NotificationsStatusStatusGet**
> NotificationListResponse internalV1NotificationsStatusStatusGet(status)

Get notifications by status

Get notifications by status

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiNotificationsApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiNotificationsApi apiInstance = new InternalApiNotificationsApi(defaultClient);
    String status = "pending"; // String | 
    try {
      NotificationListResponse result = apiInstance.internalV1NotificationsStatusStatusGet(status);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiNotificationsApi#internalV1NotificationsStatusStatusGet");
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
| **status** | **String**|  | [enum: pending, sent, failed, cancelled] |

### Return type

[**NotificationListResponse**](NotificationListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of notifications with status |  -  |

