# ProjectCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CompanyId** | **string** | Company ID | 
**NotifyKey** | **string** | Project notify key | 
**Description** | **string** | Project description | 
**DailyLimit** | Pointer to **int32** | Daily notification limit | [optional] 
**MonthlyLimit** | Pointer to **int32** | Monthly notification limit | [optional] 
**Priority** | Pointer to **string** | Project priority | [optional] [default to "normal"]
**CreatedBy** | **string** | Creator user ID | 

## Methods

### NewProjectCreateRequest

`func NewProjectCreateRequest(companyId string, notifyKey string, description string, createdBy string, ) *ProjectCreateRequest`

NewProjectCreateRequest instantiates a new ProjectCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectCreateRequestWithDefaults

`func NewProjectCreateRequestWithDefaults() *ProjectCreateRequest`

NewProjectCreateRequestWithDefaults instantiates a new ProjectCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCompanyId

`func (o *ProjectCreateRequest) GetCompanyId() string`

GetCompanyId returns the CompanyId field if non-nil, zero value otherwise.

### GetCompanyIdOk

`func (o *ProjectCreateRequest) GetCompanyIdOk() (*string, bool)`

GetCompanyIdOk returns a tuple with the CompanyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyId

`func (o *ProjectCreateRequest) SetCompanyId(v string)`

SetCompanyId sets CompanyId field to given value.


### GetNotifyKey

`func (o *ProjectCreateRequest) GetNotifyKey() string`

GetNotifyKey returns the NotifyKey field if non-nil, zero value otherwise.

### GetNotifyKeyOk

`func (o *ProjectCreateRequest) GetNotifyKeyOk() (*string, bool)`

GetNotifyKeyOk returns a tuple with the NotifyKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifyKey

`func (o *ProjectCreateRequest) SetNotifyKey(v string)`

SetNotifyKey sets NotifyKey field to given value.


### GetDescription

`func (o *ProjectCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ProjectCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ProjectCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetDailyLimit

`func (o *ProjectCreateRequest) GetDailyLimit() int32`

GetDailyLimit returns the DailyLimit field if non-nil, zero value otherwise.

### GetDailyLimitOk

`func (o *ProjectCreateRequest) GetDailyLimitOk() (*int32, bool)`

GetDailyLimitOk returns a tuple with the DailyLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDailyLimit

`func (o *ProjectCreateRequest) SetDailyLimit(v int32)`

SetDailyLimit sets DailyLimit field to given value.

### HasDailyLimit

`func (o *ProjectCreateRequest) HasDailyLimit() bool`

HasDailyLimit returns a boolean if a field has been set.

### GetMonthlyLimit

`func (o *ProjectCreateRequest) GetMonthlyLimit() int32`

GetMonthlyLimit returns the MonthlyLimit field if non-nil, zero value otherwise.

### GetMonthlyLimitOk

`func (o *ProjectCreateRequest) GetMonthlyLimitOk() (*int32, bool)`

GetMonthlyLimitOk returns a tuple with the MonthlyLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonthlyLimit

`func (o *ProjectCreateRequest) SetMonthlyLimit(v int32)`

SetMonthlyLimit sets MonthlyLimit field to given value.

### HasMonthlyLimit

`func (o *ProjectCreateRequest) HasMonthlyLimit() bool`

HasMonthlyLimit returns a boolean if a field has been set.

### GetPriority

`func (o *ProjectCreateRequest) GetPriority() string`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *ProjectCreateRequest) GetPriorityOk() (*string, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *ProjectCreateRequest) SetPriority(v string)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *ProjectCreateRequest) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetCreatedBy

`func (o *ProjectCreateRequest) GetCreatedBy() string`

GetCreatedBy returns the CreatedBy field if non-nil, zero value otherwise.

### GetCreatedByOk

`func (o *ProjectCreateRequest) GetCreatedByOk() (*string, bool)`

GetCreatedByOk returns a tuple with the CreatedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedBy

`func (o *ProjectCreateRequest) SetCreatedBy(v string)`

SetCreatedBy sets CreatedBy field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


