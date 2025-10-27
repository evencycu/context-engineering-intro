# SystemHealthResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | Pointer to **string** |  | [optional] 
**Timestamp** | Pointer to **time.Time** |  | [optional] 
**Services** | Pointer to [**map[string]SystemHealthResponseServicesValue**](SystemHealthResponseServicesValue.md) |  | [optional] 
**Uptime** | Pointer to **string** |  | [optional] 
**Version** | Pointer to **string** |  | [optional] 

## Methods

### NewSystemHealthResponse

`func NewSystemHealthResponse() *SystemHealthResponse`

NewSystemHealthResponse instantiates a new SystemHealthResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSystemHealthResponseWithDefaults

`func NewSystemHealthResponseWithDefaults() *SystemHealthResponse`

NewSystemHealthResponseWithDefaults instantiates a new SystemHealthResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *SystemHealthResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *SystemHealthResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *SystemHealthResponse) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *SystemHealthResponse) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetTimestamp

`func (o *SystemHealthResponse) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *SystemHealthResponse) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *SystemHealthResponse) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *SystemHealthResponse) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetServices

`func (o *SystemHealthResponse) GetServices() map[string]SystemHealthResponseServicesValue`

GetServices returns the Services field if non-nil, zero value otherwise.

### GetServicesOk

`func (o *SystemHealthResponse) GetServicesOk() (*map[string]SystemHealthResponseServicesValue, bool)`

GetServicesOk returns a tuple with the Services field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServices

`func (o *SystemHealthResponse) SetServices(v map[string]SystemHealthResponseServicesValue)`

SetServices sets Services field to given value.

### HasServices

`func (o *SystemHealthResponse) HasServices() bool`

HasServices returns a boolean if a field has been set.

### GetUptime

`func (o *SystemHealthResponse) GetUptime() string`

GetUptime returns the Uptime field if non-nil, zero value otherwise.

### GetUptimeOk

`func (o *SystemHealthResponse) GetUptimeOk() (*string, bool)`

GetUptimeOk returns a tuple with the Uptime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUptime

`func (o *SystemHealthResponse) SetUptime(v string)`

SetUptime sets Uptime field to given value.

### HasUptime

`func (o *SystemHealthResponse) HasUptime() bool`

HasUptime returns a boolean if a field has been set.

### GetVersion

`func (o *SystemHealthResponse) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *SystemHealthResponse) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *SystemHealthResponse) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *SystemHealthResponse) HasVersion() bool`

HasVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


