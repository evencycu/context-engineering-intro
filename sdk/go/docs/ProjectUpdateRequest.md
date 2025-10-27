# ProjectUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NotifyKey** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**DailyLimit** | Pointer to **int32** |  | [optional] 
**MonthlyLimit** | Pointer to **int32** |  | [optional] 
**Priority** | Pointer to **string** |  | [optional] 

## Methods

### NewProjectUpdateRequest

`func NewProjectUpdateRequest() *ProjectUpdateRequest`

NewProjectUpdateRequest instantiates a new ProjectUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectUpdateRequestWithDefaults

`func NewProjectUpdateRequestWithDefaults() *ProjectUpdateRequest`

NewProjectUpdateRequestWithDefaults instantiates a new ProjectUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNotifyKey

`func (o *ProjectUpdateRequest) GetNotifyKey() string`

GetNotifyKey returns the NotifyKey field if non-nil, zero value otherwise.

### GetNotifyKeyOk

`func (o *ProjectUpdateRequest) GetNotifyKeyOk() (*string, bool)`

GetNotifyKeyOk returns a tuple with the NotifyKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifyKey

`func (o *ProjectUpdateRequest) SetNotifyKey(v string)`

SetNotifyKey sets NotifyKey field to given value.

### HasNotifyKey

`func (o *ProjectUpdateRequest) HasNotifyKey() bool`

HasNotifyKey returns a boolean if a field has been set.

### GetDescription

`func (o *ProjectUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ProjectUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ProjectUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ProjectUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetDailyLimit

`func (o *ProjectUpdateRequest) GetDailyLimit() int32`

GetDailyLimit returns the DailyLimit field if non-nil, zero value otherwise.

### GetDailyLimitOk

`func (o *ProjectUpdateRequest) GetDailyLimitOk() (*int32, bool)`

GetDailyLimitOk returns a tuple with the DailyLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDailyLimit

`func (o *ProjectUpdateRequest) SetDailyLimit(v int32)`

SetDailyLimit sets DailyLimit field to given value.

### HasDailyLimit

`func (o *ProjectUpdateRequest) HasDailyLimit() bool`

HasDailyLimit returns a boolean if a field has been set.

### GetMonthlyLimit

`func (o *ProjectUpdateRequest) GetMonthlyLimit() int32`

GetMonthlyLimit returns the MonthlyLimit field if non-nil, zero value otherwise.

### GetMonthlyLimitOk

`func (o *ProjectUpdateRequest) GetMonthlyLimitOk() (*int32, bool)`

GetMonthlyLimitOk returns a tuple with the MonthlyLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonthlyLimit

`func (o *ProjectUpdateRequest) SetMonthlyLimit(v int32)`

SetMonthlyLimit sets MonthlyLimit field to given value.

### HasMonthlyLimit

`func (o *ProjectUpdateRequest) HasMonthlyLimit() bool`

HasMonthlyLimit returns a boolean if a field has been set.

### GetPriority

`func (o *ProjectUpdateRequest) GetPriority() string`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *ProjectUpdateRequest) GetPriorityOk() (*string, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *ProjectUpdateRequest) SetPriority(v string)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *ProjectUpdateRequest) HasPriority() bool`

HasPriority returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


