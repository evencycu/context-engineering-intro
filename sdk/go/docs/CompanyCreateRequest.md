# CompanyCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Company name | 
**ContactEmail** | **string** | Contact email | 
**ContactPhone** | Pointer to **string** | Contact phone number | [optional] 
**Address** | Pointer to **string** | Company address | [optional] 
**Status** | Pointer to **string** |  | [optional] [default to "active"]
**BillingEnabled** | Pointer to **bool** |  | [optional] [default to false]

## Methods

### NewCompanyCreateRequest

`func NewCompanyCreateRequest(name string, contactEmail string, ) *CompanyCreateRequest`

NewCompanyCreateRequest instantiates a new CompanyCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCompanyCreateRequestWithDefaults

`func NewCompanyCreateRequestWithDefaults() *CompanyCreateRequest`

NewCompanyCreateRequestWithDefaults instantiates a new CompanyCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CompanyCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CompanyCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CompanyCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetContactEmail

`func (o *CompanyCreateRequest) GetContactEmail() string`

GetContactEmail returns the ContactEmail field if non-nil, zero value otherwise.

### GetContactEmailOk

`func (o *CompanyCreateRequest) GetContactEmailOk() (*string, bool)`

GetContactEmailOk returns a tuple with the ContactEmail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContactEmail

`func (o *CompanyCreateRequest) SetContactEmail(v string)`

SetContactEmail sets ContactEmail field to given value.


### GetContactPhone

`func (o *CompanyCreateRequest) GetContactPhone() string`

GetContactPhone returns the ContactPhone field if non-nil, zero value otherwise.

### GetContactPhoneOk

`func (o *CompanyCreateRequest) GetContactPhoneOk() (*string, bool)`

GetContactPhoneOk returns a tuple with the ContactPhone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContactPhone

`func (o *CompanyCreateRequest) SetContactPhone(v string)`

SetContactPhone sets ContactPhone field to given value.

### HasContactPhone

`func (o *CompanyCreateRequest) HasContactPhone() bool`

HasContactPhone returns a boolean if a field has been set.

### GetAddress

`func (o *CompanyCreateRequest) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *CompanyCreateRequest) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *CompanyCreateRequest) SetAddress(v string)`

SetAddress sets Address field to given value.

### HasAddress

`func (o *CompanyCreateRequest) HasAddress() bool`

HasAddress returns a boolean if a field has been set.

### GetStatus

`func (o *CompanyCreateRequest) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CompanyCreateRequest) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CompanyCreateRequest) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *CompanyCreateRequest) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetBillingEnabled

`func (o *CompanyCreateRequest) GetBillingEnabled() bool`

GetBillingEnabled returns the BillingEnabled field if non-nil, zero value otherwise.

### GetBillingEnabledOk

`func (o *CompanyCreateRequest) GetBillingEnabledOk() (*bool, bool)`

GetBillingEnabledOk returns a tuple with the BillingEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBillingEnabled

`func (o *CompanyCreateRequest) SetBillingEnabled(v bool)`

SetBillingEnabled sets BillingEnabled field to given value.

### HasBillingEnabled

`func (o *CompanyCreateRequest) HasBillingEnabled() bool`

HasBillingEnabled returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


