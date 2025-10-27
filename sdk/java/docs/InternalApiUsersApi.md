# InternalApiUsersApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**internalV1UsersCompanyCompanyIdGet**](InternalApiUsersApi.md#internalV1UsersCompanyCompanyIdGet) | **GET** /internal/v1/users/company/{companyId} | Get users by company |
| [**internalV1UsersGet**](InternalApiUsersApi.md#internalV1UsersGet) | **GET** /internal/v1/users | List users |
| [**internalV1UsersIdDelete**](InternalApiUsersApi.md#internalV1UsersIdDelete) | **DELETE** /internal/v1/users/{id} | Delete user |
| [**internalV1UsersIdGet**](InternalApiUsersApi.md#internalV1UsersIdGet) | **GET** /internal/v1/users/{id} | Get user |
| [**internalV1UsersIdPasswordPatch**](InternalApiUsersApi.md#internalV1UsersIdPasswordPatch) | **PATCH** /internal/v1/users/{id}/password | Change user password |
| [**internalV1UsersIdPut**](InternalApiUsersApi.md#internalV1UsersIdPut) | **PUT** /internal/v1/users/{id} | Update user |
| [**internalV1UsersPost**](InternalApiUsersApi.md#internalV1UsersPost) | **POST** /internal/v1/users | Create user |
| [**internalV1UsersRoleRoleGet**](InternalApiUsersApi.md#internalV1UsersRoleRoleGet) | **GET** /internal/v1/users/role/{role} | Get users by role |


<a id="internalV1UsersCompanyCompanyIdGet"></a>
# **internalV1UsersCompanyCompanyIdGet**
> UserListResponse internalV1UsersCompanyCompanyIdGet(companyId)

Get users by company

Get users by company ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiUsersApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiUsersApi apiInstance = new InternalApiUsersApi(defaultClient);
    UUID companyId = UUID.randomUUID(); // UUID | 
    try {
      UserListResponse result = apiInstance.internalV1UsersCompanyCompanyIdGet(companyId);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiUsersApi#internalV1UsersCompanyCompanyIdGet");
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

[**UserListResponse**](UserListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of users in company |  -  |
| **404** | Not found |  -  |

<a id="internalV1UsersGet"></a>
# **internalV1UsersGet**
> UserListResponse internalV1UsersGet(limit, offset, search)

List users

Get a list of users with pagination

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiUsersApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiUsersApi apiInstance = new InternalApiUsersApi(defaultClient);
    Integer limit = 10; // Integer | Number of items to return
    Integer offset = 0; // Integer | Number of items to skip
    String search = "search_example"; // String | Search term
    try {
      UserListResponse result = apiInstance.internalV1UsersGet(limit, offset, search);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiUsersApi#internalV1UsersGet");
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

[**UserListResponse**](UserListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of users |  -  |

<a id="internalV1UsersIdDelete"></a>
# **internalV1UsersIdDelete**
> internalV1UsersIdDelete(id)

Delete user

Delete user by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiUsersApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiUsersApi apiInstance = new InternalApiUsersApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      apiInstance.internalV1UsersIdDelete(id);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiUsersApi#internalV1UsersIdDelete");
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
| **204** | User deleted successfully |  -  |
| **404** | Not found |  -  |

<a id="internalV1UsersIdGet"></a>
# **internalV1UsersIdGet**
> UserResponse internalV1UsersIdGet(id)

Get user

Get user by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiUsersApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiUsersApi apiInstance = new InternalApiUsersApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    try {
      UserResponse result = apiInstance.internalV1UsersIdGet(id);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiUsersApi#internalV1UsersIdGet");
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

[**UserResponse**](UserResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | User details |  -  |
| **404** | Not found |  -  |

<a id="internalV1UsersIdPasswordPatch"></a>
# **internalV1UsersIdPasswordPatch**
> internalV1UsersIdPasswordPatch(id, changePasswordRequest)

Change user password

Change user password by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiUsersApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiUsersApi apiInstance = new InternalApiUsersApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    ChangePasswordRequest changePasswordRequest = new ChangePasswordRequest(); // ChangePasswordRequest | 
    try {
      apiInstance.internalV1UsersIdPasswordPatch(id, changePasswordRequest);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiUsersApi#internalV1UsersIdPasswordPatch");
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
| **changePasswordRequest** | [**ChangePasswordRequest**](ChangePasswordRequest.md)|  | |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Password changed successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1UsersIdPut"></a>
# **internalV1UsersIdPut**
> UserResponse internalV1UsersIdPut(id, userUpdateRequest)

Update user

Update user by ID

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiUsersApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiUsersApi apiInstance = new InternalApiUsersApi(defaultClient);
    UUID id = UUID.randomUUID(); // UUID | 
    UserUpdateRequest userUpdateRequest = new UserUpdateRequest(); // UserUpdateRequest | 
    try {
      UserResponse result = apiInstance.internalV1UsersIdPut(id, userUpdateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiUsersApi#internalV1UsersIdPut");
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
| **userUpdateRequest** | [**UserUpdateRequest**](UserUpdateRequest.md)|  | |

### Return type

[**UserResponse**](UserResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | User updated successfully |  -  |
| **400** | Bad request |  -  |
| **404** | Not found |  -  |

<a id="internalV1UsersPost"></a>
# **internalV1UsersPost**
> UserResponse internalV1UsersPost(userCreateRequest)

Create user

Create a new user

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiUsersApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiUsersApi apiInstance = new InternalApiUsersApi(defaultClient);
    UserCreateRequest userCreateRequest = new UserCreateRequest(); // UserCreateRequest | 
    try {
      UserResponse result = apiInstance.internalV1UsersPost(userCreateRequest);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiUsersApi#internalV1UsersPost");
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
| **userCreateRequest** | [**UserCreateRequest**](UserCreateRequest.md)|  | |

### Return type

[**UserResponse**](UserResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | User created successfully |  -  |
| **400** | Bad request |  -  |

<a id="internalV1UsersRoleRoleGet"></a>
# **internalV1UsersRoleRoleGet**
> UserListResponse internalV1UsersRoleRoleGet(role)

Get users by role

Get users by role

### Example
```java
// Import classes:
import org.openapitools.client.ApiClient;
import org.openapitools.client.ApiException;
import org.openapitools.client.Configuration;
import org.openapitools.client.models.*;
import org.openapitools.client.api.InternalApiUsersApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://localhost:8080");

    InternalApiUsersApi apiInstance = new InternalApiUsersApi(defaultClient);
    String role = "admin"; // String | 
    try {
      UserListResponse result = apiInstance.internalV1UsersRoleRoleGet(role);
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling InternalApiUsersApi#internalV1UsersRoleRoleGet");
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
| **role** | **String**|  | [enum: admin, manager, user] |

### Return type

[**UserListResponse**](UserListResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | List of users with role |  -  |

