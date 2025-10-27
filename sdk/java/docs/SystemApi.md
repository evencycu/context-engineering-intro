# SystemApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**configGet**](SystemApi.md#configGet) | **GET** /config | Get system configuration |
| [**configValidatePost**](SystemApi.md#configValidatePost) | **POST** /config/validate | Validate system configuration |
| [**healthGet**](SystemApi.md#healthGet) | **GET** /health | Service health check |
| [**metricsGet**](SystemApi.md#metricsGet) | **GET** /metrics | System metrics |


<a id="configGet"></a>
# **configGet**
> ConfigResponse configGet()

Get system configuration

Returns current system configuration (non-sensitive values only)

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.auth.*;
import org.openapitools.client.models.*;
import org.openapitools.client.api.SystemApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");
    
    // Configure HTTP bearer authorization: bearerAuth
    HttpBearerAuth bearerAuth = (HttpBearerAuth) defaultClient.getAuthentication("bearerAuth");
    bearerAuth.setBearerToken("BEARER TOKEN");

    SystemApi apiInstance = new SystemApi(defaultClient);
    try {
      ConfigResponse result = apiInstance.configGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling SystemApi#configGet");
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

[**ConfigResponse**](ConfigResponse.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | System configuration |  -  |
| **401** | Unauthorized |  -  |

<a id="configValidatePost"></a>
# **configValidatePost**
> ConfigValidationResponse configValidatePost()

Validate system configuration

Validates the current system configuration

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.auth.*;
import org.openapitools.client.models.*;
import org.openapitools.client.api.SystemApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");
    
    // Configure HTTP bearer authorization: bearerAuth
    HttpBearerAuth bearerAuth = (HttpBearerAuth) defaultClient.getAuthentication("bearerAuth");
    bearerAuth.setBearerToken("BEARER TOKEN");

    SystemApi apiInstance = new SystemApi(defaultClient);
    try {
      ConfigValidationResponse result = apiInstance.configValidatePost();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling SystemApi#configValidatePost");
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

[**ConfigValidationResponse**](ConfigValidationResponse.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Configuration validation result |  -  |
| **401** | Unauthorized |  -  |

<a id="healthGet"></a>
# **healthGet**
> HealthResponse healthGet()

Service health check

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.SystemApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    SystemApi apiInstance = new SystemApi(defaultClient);
    try {
      HealthResponse result = apiInstance.healthGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling SystemApi#healthGet");
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

[**HealthResponse**](HealthResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Health status |  -  |

<a id="metricsGet"></a>
# **metricsGet**
> MetricsResponse metricsGet()

System metrics

Returns system metrics including token cache statistics, actor pool status, and performance indicators

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.SystemApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    SystemApi apiInstance = new SystemApi(defaultClient);
    try {
      MetricsResponse result = apiInstance.metricsGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling SystemApi#metricsGet");
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

[**MetricsResponse**](MetricsResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | System metrics |  -  |

