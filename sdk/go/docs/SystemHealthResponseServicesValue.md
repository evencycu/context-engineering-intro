# SystemHealthResponseServicesValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | Pointer to **string** |  | [optional] 
**ResponseTime** | Pointer to **float32** |  | [optional] 
**LastCheck** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewSystemHealthResponseServicesValue

`func NewSystemHealthResponseServicesValue() *SystemHealthResponseServicesValue`

NewSystemHealthResponseServicesValue instantiates a new SystemHealthResponseServicesValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSystemHealthResponseServicesValueWithDefaults

`func NewSystemHealthResponseServicesValueWithDefaults() *SystemHealthResponseServicesValue`

NewSystemHealthResponseServicesValueWithDefaults instantiates a new SystemHealthResponseServicesValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *SystemHealthResponseServicesValue) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *SystemHealthResponseServicesValue) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *SystemHealthResponseServicesValue) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *SystemHealthResponseServicesValue) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetResponseTime

`func (o *SystemHealthResponseServicesValue) GetResponseTime() float32`

GetResponseTime returns the ResponseTime field if non-nil, zero value otherwise.

### GetResponseTimeOk

`func (o *SystemHealthResponseServicesValue) GetResponseTimeOk() (*float32, bool)`

GetResponseTimeOk returns a tuple with the ResponseTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseTime

`func (o *SystemHealthResponseServicesValue) SetResponseTime(v float32)`

SetResponseTime sets ResponseTime field to given value.

### HasResponseTime

`func (o *SystemHealthResponseServicesValue) HasResponseTime() bool`

HasResponseTime returns a boolean if a field has been set.

### GetLastCheck

`func (o *SystemHealthResponseServicesValue) GetLastCheck() time.Time`

GetLastCheck returns the LastCheck field if non-nil, zero value otherwise.

### GetLastCheckOk

`func (o *SystemHealthResponseServicesValue) GetLastCheckOk() (*time.Time, bool)`

GetLastCheckOk returns a tuple with the LastCheck field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastCheck

`func (o *SystemHealthResponseServicesValue) SetLastCheck(v time.Time)`

SetLastCheck sets LastCheck field to given value.

### HasLastCheck

`func (o *SystemHealthResponseServicesValue) HasLastCheck() bool`

HasLastCheck returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


