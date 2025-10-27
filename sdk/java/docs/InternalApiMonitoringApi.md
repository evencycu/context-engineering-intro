# InternalApiMonitoringApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1MonitoringAlertsGet**](InternalApiMonitoringApi.md#internalV1MonitoringAlertsGet) | **GET** /internal/v1/monitoring/alerts | Get alert status |
| [**internalV1MonitoringBusinessGet**](InternalApiMonitoringApi.md#internalV1MonitoringBusinessGet) | **GET** /internal/v1/monitoring/business | Get business health |
| [**internalV1MonitoringDashboardGet**](InternalApiMonitoringApi.md#internalV1MonitoringDashboardGet) | **GET** /internal/v1/monitoring/dashboard | Get monitoring dashboard |
| [**internalV1MonitoringHealthGet**](InternalApiMonitoringApi.md#internalV1MonitoringHealthGet) | **GET** /internal/v1/monitoring/health | Get system health |
| [**internalV1MonitoringPerformanceGet**](InternalApiMonitoringApi.md#internalV1MonitoringPerformanceGet) | **GET** /internal/v1/monitoring/performance | Get performance metrics |


<a id="internalV1MonitoringAlertsGet"></a>
# **internalV1MonitoringAlertsGet**
> AlertStatusResponse internalV1MonitoringAlertsGet()

Get alert status

Get current alert status

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiMonitoringApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiMonitoringApi apiInstance = new InternalApiMonitoringApi(defaultClient);
    try {
      AlertStatusResponse result = apiInstance.internalV1MonitoringAlertsGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiMonitoringApi#internalV1MonitoringAlertsGet");
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

[**AlertStatusResponse**](AlertStatusResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Alert status |  -  |

<a id="internalV1MonitoringBusinessGet"></a>
# **internalV1MonitoringBusinessGet**
> BusinessHealthResponse internalV1MonitoringBusinessGet()

Get business health

Get business metrics and health indicators

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiMonitoringApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiMonitoringApi apiInstance = new InternalApiMonitoringApi(defaultClient);
    try {
      BusinessHealthResponse result = apiInstance.internalV1MonitoringBusinessGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiMonitoringApi#internalV1MonitoringBusinessGet");
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

[**BusinessHealthResponse**](BusinessHealthResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Business health metrics |  -  |

<a id="internalV1MonitoringDashboardGet"></a>
# **internalV1MonitoringDashboardGet**
> DashboardResponse internalV1MonitoringDashboardGet()

Get monitoring dashboard

Get comprehensive monitoring dashboard data

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiMonitoringApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiMonitoringApi apiInstance = new InternalApiMonitoringApi(defaultClient);
    try {
      DashboardResponse result = apiInstance.internalV1MonitoringDashboardGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiMonitoringApi#internalV1MonitoringDashboardGet");
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

[**DashboardResponse**](DashboardResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Dashboard data |  -  |

<a id="internalV1MonitoringHealthGet"></a>
# **internalV1MonitoringHealthGet**
> SystemHealthResponse internalV1MonitoringHealthGet()

Get system health

Get comprehensive system health status

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiMonitoringApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiMonitoringApi apiInstance = new InternalApiMonitoringApi(defaultClient);
    try {
      SystemHealthResponse result = apiInstance.internalV1MonitoringHealthGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiMonitoringApi#internalV1MonitoringHealthGet");
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

[**SystemHealthResponse**](SystemHealthResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | System health status |  -  |

<a id="internalV1MonitoringPerformanceGet"></a>
# **internalV1MonitoringPerformanceGet**
> PerformanceMetricsResponse internalV1MonitoringPerformanceGet()

Get performance metrics

Get system performance metrics

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiMonitoringApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiMonitoringApi apiInstance = new InternalApiMonitoringApi(defaultClient);
    try {
      PerformanceMetricsResponse result = apiInstance.internalV1MonitoringPerformanceGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiMonitoringApi#internalV1MonitoringPerformanceGet");
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

[**PerformanceMetricsResponse**](PerformanceMetricsResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Performance metrics |  -  |

