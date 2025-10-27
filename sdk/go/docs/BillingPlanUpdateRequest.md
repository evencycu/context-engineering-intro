# BillingPlanUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**NotificationLimit** | Pointer to **int32** |  | [optional] 
**Features** | Pointer to **[]string** |  | [optional] 

## Methods

### NewBillingPlanUpdateRequest

`func NewBillingPlanUpdateRequest() *BillingPlanUpdateRequest`

NewBillingPlanUpdateRequest instantiates a new BillingPlanUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBillingPlanUpdateRequestWithDefaults

`func NewBillingPlanUpdateRequestWithDefaults() *BillingPlanUpdateRequest`

NewBillingPlanUpdateRequestWithDefaults instantiates a new BillingPlanUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BillingPlanUpdateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BillingPlanUpdateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BillingPlanUpdateRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *BillingPlanUpdateRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *BillingPlanUpdateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *BillingPlanUpdateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *BillingPlanUpdateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *BillingPlanUpdateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetPrice

`func (o *BillingPlanUpdateRequest) GetPrice() float32`

GetPrice returns the Price field if non-nil, zero value otherwise.

### GetPriceOk

`func (o *BillingPlanUpdateRequest) GetPriceOk() (*float32, bool)`

GetPriceOk returns a tuple with the Price field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrice

`func (o *BillingPlanUpdateRequest) SetPrice(v float32)`

SetPrice sets Price field to given value.

### HasPrice

`func (o *BillingPlanUpdateRequest) HasPrice() bool`

HasPrice returns a boolean if a field has been set.

### GetNotificationLimit

`func (o *BillingPlanUpdateRequest) GetNotificationLimit() int32`

GetNotificationLimit returns the NotificationLimit field if non-nil, zero value otherwise.

### GetNotificationLimitOk

`func (o *BillingPlanUpdateRequest) GetNotificationLimitOk() (*int32, bool)`

GetNotificationLimitOk returns a tuple with the NotificationLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotificationLimit

`func (o *BillingPlanUpdateRequest) SetNotificationLimit(v int32)`

SetNotificationLimit sets NotificationLimit field to given value.

### HasNotificationLimit

`func (o *BillingPlanUpdateRequest) HasNotificationLimit() bool`

HasNotificationLimit returns a boolean if a field has been set.

### GetFeatures

`func (o *BillingPlanUpdateRequest) GetFeatures() []string`

GetFeatures returns the Features field if non-nil, zero value otherwise.

### GetFeaturesOk

`func (o *BillingPlanUpdateRequest) GetFeaturesOk() (*[]string, bool)`

GetFeaturesOk returns a tuple with the Features field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeatures

`func (o *BillingPlanUpdateRequest) SetFeatures(v []string)`

SetFeatures sets Features field to given value.

### HasFeatures

`func (o *BillingPlanUpdateRequest) HasFeatures() bool`

HasFeatures returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


