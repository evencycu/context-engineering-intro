# ProvisionCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NotifyKey** | **string** | Project notify key | 
**CompanyId** | **string** | Company ID | 
**ProjectName** | **string** | Project name | 
**Description** | Pointer to **string** | Project description | [optional] 

## Methods

### NewProvisionCreateRequest

`func NewProvisionCreateRequest(notifyKey string, companyId string, projectName string, ) *ProvisionCreateRequest`

NewProvisionCreateRequest instantiates a new ProvisionCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProvisionCreateRequestWithDefaults

`func NewProvisionCreateRequestWithDefaults() *ProvisionCreateRequest`

NewProvisionCreateRequestWithDefaults instantiates a new ProvisionCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNotifyKey

`func (o *ProvisionCreateRequest) GetNotifyKey() string`

GetNotifyKey returns the NotifyKey field if non-nil, zero value otherwise.

### GetNotifyKeyOk

`func (o *ProvisionCreateRequest) GetNotifyKeyOk() (*string, bool)`

GetNotifyKeyOk returns a tuple with the NotifyKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifyKey

`func (o *ProvisionCreateRequest) SetNotifyKey(v string)`

SetNotifyKey sets NotifyKey field to given value.


### GetCompanyId

`func (o *ProvisionCreateRequest) GetCompanyId() string`

GetCompanyId returns the CompanyId field if non-nil, zero value otherwise.

### GetCompanyIdOk

`func (o *ProvisionCreateRequest) GetCompanyIdOk() (*string, bool)`

GetCompanyIdOk returns a tuple with the CompanyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyId

`func (o *ProvisionCreateRequest) SetCompanyId(v string)`

SetCompanyId sets CompanyId field to given value.


### GetProjectName

`func (o *ProvisionCreateRequest) GetProjectName() string`

GetProjectName returns the ProjectName field if non-nil, zero value otherwise.

### GetProjectNameOk

`func (o *ProvisionCreateRequest) GetProjectNameOk() (*string, bool)`

GetProjectNameOk returns a tuple with the ProjectName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectName

`func (o *ProvisionCreateRequest) SetProjectName(v string)`

SetProjectName sets ProjectName field to given value.


### GetDescription

`func (o *ProvisionCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ProvisionCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ProvisionCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ProvisionCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


