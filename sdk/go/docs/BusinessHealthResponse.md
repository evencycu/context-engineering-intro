# BusinessHealthResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Timestamp** | Pointer to **time.Time** |  | [optional] 
**Notifications** | Pointer to [**BusinessHealthResponseNotifications**](BusinessHealthResponseNotifications.md) |  | [optional] 
**Users** | Pointer to [**BusinessHealthResponseUsers**](BusinessHealthResponseUsers.md) |  | [optional] 
**Companies** | Pointer to [**BusinessHealthResponseUsers**](BusinessHealthResponseUsers.md) |  | [optional] 

## Methods

### NewBusinessHealthResponse

`func NewBusinessHealthResponse() *BusinessHealthResponse`

NewBusinessHealthResponse instantiates a new BusinessHealthResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBusinessHealthResponseWithDefaults

`func NewBusinessHealthResponseWithDefaults() *BusinessHealthResponse`

NewBusinessHealthResponseWithDefaults instantiates a new BusinessHealthResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTimestamp

`func (o *BusinessHealthResponse) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *BusinessHealthResponse) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *BusinessHealthResponse) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *BusinessHealthResponse) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetNotifications

`func (o *BusinessHealthResponse) GetNotifications() BusinessHealthResponseNotifications`

GetNotifications returns the Notifications field if non-nil, zero value otherwise.

### GetNotificationsOk

`func (o *BusinessHealthResponse) GetNotificationsOk() (*BusinessHealthResponseNotifications, bool)`

GetNotificationsOk returns a tuple with the Notifications field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifications

`func (o *BusinessHealthResponse) SetNotifications(v BusinessHealthResponseNotifications)`

SetNotifications sets Notifications field to given value.

### HasNotifications

`func (o *BusinessHealthResponse) HasNotifications() bool`

HasNotifications returns a boolean if a field has been set.

### GetUsers

`func (o *BusinessHealthResponse) GetUsers() BusinessHealthResponseUsers`

GetUsers returns the Users field if non-nil, zero value otherwise.

### GetUsersOk

`func (o *BusinessHealthResponse) GetUsersOk() (*BusinessHealthResponseUsers, bool)`

GetUsersOk returns a tuple with the Users field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsers

`func (o *BusinessHealthResponse) SetUsers(v BusinessHealthResponseUsers)`

SetUsers sets Users field to given value.

### HasUsers

`func (o *BusinessHealthResponse) HasUsers() bool`

HasUsers returns a boolean if a field has been set.

### GetCompanies

`func (o *BusinessHealthResponse) GetCompanies() BusinessHealthResponseUsers`

GetCompanies returns the Companies field if non-nil, zero value otherwise.

### GetCompaniesOk

`func (o *BusinessHealthResponse) GetCompaniesOk() (*BusinessHealthResponseUsers, bool)`

GetCompaniesOk returns a tuple with the Companies field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanies

`func (o *BusinessHealthResponse) SetCompanies(v BusinessHealthResponseUsers)`

SetCompanies sets Companies field to given value.

### HasCompanies

`func (o *BusinessHealthResponse) HasCompanies() bool`

HasCompanies returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


