# ProjectResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**CompanyId** | Pointer to **string** |  | [optional] 
**NotifyKey** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**DailyLimit** | Pointer to **int32** |  | [optional] 
**MonthlyLimit** | Pointer to **int32** |  | [optional] 
**Priority** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewProjectResponse

`func NewProjectResponse() *ProjectResponse`

NewProjectResponse instantiates a new ProjectResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectResponseWithDefaults

`func NewProjectResponseWithDefaults() *ProjectResponse`

NewProjectResponseWithDefaults instantiates a new ProjectResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ProjectResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ProjectResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ProjectResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ProjectResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetCompanyId

`func (o *ProjectResponse) GetCompanyId() string`

GetCompanyId returns the CompanyId field if non-nil, zero value otherwise.

### GetCompanyIdOk

`func (o *ProjectResponse) GetCompanyIdOk() (*string, bool)`

GetCompanyIdOk returns a tuple with the CompanyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyId

`func (o *ProjectResponse) SetCompanyId(v string)`

SetCompanyId sets CompanyId field to given value.

### HasCompanyId

`func (o *ProjectResponse) HasCompanyId() bool`

HasCompanyId returns a boolean if a field has been set.

### GetNotifyKey

`func (o *ProjectResponse) GetNotifyKey() string`

GetNotifyKey returns the NotifyKey field if non-nil, zero value otherwise.

### GetNotifyKeyOk

`func (o *ProjectResponse) GetNotifyKeyOk() (*string, bool)`

GetNotifyKeyOk returns a tuple with the NotifyKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifyKey

`func (o *ProjectResponse) SetNotifyKey(v string)`

SetNotifyKey sets NotifyKey field to given value.

### HasNotifyKey

`func (o *ProjectResponse) HasNotifyKey() bool`

HasNotifyKey returns a boolean if a field has been set.

### GetDescription

`func (o *ProjectResponse) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ProjectResponse) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ProjectResponse) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ProjectResponse) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetDailyLimit

`func (o *ProjectResponse) GetDailyLimit() int32`

GetDailyLimit returns the DailyLimit field if non-nil, zero value otherwise.

### GetDailyLimitOk

`func (o *ProjectResponse) GetDailyLimitOk() (*int32, bool)`

GetDailyLimitOk returns a tuple with the DailyLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDailyLimit

`func (o *ProjectResponse) SetDailyLimit(v int32)`

SetDailyLimit sets DailyLimit field to given value.

### HasDailyLimit

`func (o *ProjectResponse) HasDailyLimit() bool`

HasDailyLimit returns a boolean if a field has been set.

### GetMonthlyLimit

`func (o *ProjectResponse) GetMonthlyLimit() int32`

GetMonthlyLimit returns the MonthlyLimit field if non-nil, zero value otherwise.

### GetMonthlyLimitOk

`func (o *ProjectResponse) GetMonthlyLimitOk() (*int32, bool)`

GetMonthlyLimitOk returns a tuple with the MonthlyLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonthlyLimit

`func (o *ProjectResponse) SetMonthlyLimit(v int32)`

SetMonthlyLimit sets MonthlyLimit field to given value.

### HasMonthlyLimit

`func (o *ProjectResponse) HasMonthlyLimit() bool`

HasMonthlyLimit returns a boolean if a field has been set.

### GetPriority

`func (o *ProjectResponse) GetPriority() string`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *ProjectResponse) GetPriorityOk() (*string, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *ProjectResponse) SetPriority(v string)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *ProjectResponse) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetStatus

`func (o *ProjectResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ProjectResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ProjectResponse) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ProjectResponse) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetCreatedAt

`func (o *ProjectResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ProjectResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ProjectResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *ProjectResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *ProjectResponse) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *ProjectResponse) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *ProjectResponse) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *ProjectResponse) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


