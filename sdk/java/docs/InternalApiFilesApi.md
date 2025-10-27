# InternalApiFilesApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1FilesGet**](InternalApiFilesApi.md#internalV1FilesGet) | **GET** /internal/v1/files | List files |
| [**internalV1FilesIdDelete**](InternalApiFilesApi.md#internalV1FilesIdDelete) | **DELETE** /internal/v1/files/{id} | Delete file |
| [**internalV1FilesIdDownloadGet**](InternalApiFilesApi.md#internalV1FilesIdDownloadGet) | **GET** /internal/v1/files/{id}/download | Download file |
| [**internalV1FilesIdGet**](InternalApiFilesApi.md#internalV1FilesIdGet) | **GET** /internal/v1/files/{id} | Get file |
| [**internalV1FilesUploadMultiplePost**](InternalApiFilesApi.md#internalV1FilesUploadMultiplePost) | **POST** /internal/v1/files/upload/multiple | Upload multiple files |
| [**internalV1FilesUploadPost**](InternalApiFilesApi.md#internalV1FilesUploadPost) | **POST** /internal/v1/files/upload | Upload file |
| [**internalV1FilesValidatePost**](InternalApiFilesApi.md#internalV1FilesValidatePost) | **POST** /internal/v1/files/validate | Validate file |


<a id="internalV1FilesGet"></a>
# **internalV1FilesGet**
> FileListResponse internalV1FilesGet(limit, offset, search)

List files

Get a list of files with pagination

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiFilesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiFilesApi apiInstance = new InternalApiFilesApi(defaultClient);
    Integer limit = 10; // Integer | Number of items to return
    Integer offset = 0; // Integer | Number of items to skip
    String search = "search_example"; // String | Search term
    try {
      FileListResponse result = apiInstance.internalV1FilesGet(limit, offset, search);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiFilesApi#internalV1FilesGet");
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

[**FileListResponse**](FileListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of files |  -  |

<a id="internalV1FilesIdDelete"></a>
# **internalV1FilesIdDelete**
> internalV1FilesIdDelete(id)

Delete file

Delete file by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiFilesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiFilesApi apiInstance = new InternalApiFilesApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      apiInstance.internalV1FilesIdDelete(id);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiFilesApi#internalV1FilesIdDelete");
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
| **204** | File deleted successfully |  -  |
| **404** | Not found |  -  |

<a id="internalV1FilesIdDownloadGet"></a>
# **internalV1FilesIdDownloadGet**
> File internalV1FilesIdDownloadGet(id)

Download file

Download file by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiFilesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiFilesApi apiInstance = new InternalApiFilesApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      File result = apiInstance.internalV1FilesIdDownloadGet(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiFilesApi#internalV1FilesIdDownloadGet");
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

[**File**](File.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream, application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | File content |  -  |
| **404** | Not found |  -  |

<a id="internalV1FilesIdGet"></a>
# **internalV1FilesIdGet**
> FileResponse internalV1FilesIdGet(id)

Get file

Get file by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiFilesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiFilesApi apiInstance = new InternalApiFilesApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      FileResponse result = apiInstance.internalV1FilesIdGet(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiFilesApi#internalV1FilesIdGet");
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

[**FileResponse**](FileResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | File details |  -  |
| **404** | Not found |  -  |

<a id="internalV1FilesUploadMultiplePost"></a>
# **internalV1FilesUploadMultiplePost**
> MultipleFileUploadResponse internalV1FilesUploadMultiplePost(files, description)

Upload multiple files

Upload multiple files

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiFilesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiFilesApi apiInstance = new InternalApiFilesApi(defaultClient);
    List<File> files = Arrays.asList(); // List<File> | Files to upload
    String description = "description_example"; // String | Upload description
    try {
      MultipleFileUploadResponse result = apiInstance.internalV1FilesUploadMultiplePost(files, description);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiFilesApi#internalV1FilesUploadMultiplePost");
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
| **files** | **List&lt;File&gt;**| Files to upload | |
| **description** | **String**| Upload description | [optional] |

### Return type

[**MultipleFileUploadResponse**](MultipleFileUploadResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Files uploaded successfully |  -  |
| **400** | Bad request |  -  |

<a id="internalV1FilesUploadPost"></a>
# **internalV1FilesUploadPost**
> FileResponse internalV1FilesUploadPost(_file, description, tags)

Upload file

Upload a single file

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiFilesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiFilesApi apiInstance = new InternalApiFilesApi(defaultClient);
    File _file = new File("/path/to/file"); // File | File to upload
    String description = "description_example"; // String | File description
    List<String> tags = Arrays.asList(); // List<String> | File tags
    try {
      FileResponse result = apiInstance.internalV1FilesUploadPost(_file, description, tags);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiFilesApi#internalV1FilesUploadPost");
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
| **_file** | **File**| File to upload | |
| **description** | **String**| File description | [optional] |
| **tags** | [**List&lt;String&gt;**](String.md)| File tags | [optional] |

### Return type

[**FileResponse**](FileResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | File uploaded successfully |  -  |
| **400** | Bad request |  -  |

<a id="internalV1FilesValidatePost"></a>
# **internalV1FilesValidatePost**
> FileValidationResponse internalV1FilesValidatePost(_file)

Validate file

Validate file before upload

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiFilesApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiFilesApi apiInstance = new InternalApiFilesApi(defaultClient);
    File _file = new File("/path/to/file"); // File | File to validate
    try {
      FileValidationResponse result = apiInstance.internalV1FilesValidatePost(_file);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiFilesApi#internalV1FilesValidatePost");
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
| **_file** | **File**| File to validate | |

### Return type

[**FileValidationResponse**](FileValidationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | File validation result |  -  |

