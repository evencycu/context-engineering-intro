# BillingPlanCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Plan name | 
**Description** | Pointer to **string** | Plan description | [optional] 
**Price** | **float32** | Monthly price | 
**NotificationLimit** | **int32** | Monthly notification limit | 
**Features** | Pointer to **[]string** | Plan features | [optional] 

## Methods

### NewBillingPlanCreateRequest

`func NewBillingPlanCreateRequest(name string, price float32, notificationLimit int32, ) *BillingPlanCreateRequest`

NewBillingPlanCreateRequest instantiates a new BillingPlanCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBillingPlanCreateRequestWithDefaults

`func NewBillingPlanCreateRequestWithDefaults() *BillingPlanCreateRequest`

NewBillingPlanCreateRequestWithDefaults instantiates a new BillingPlanCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BillingPlanCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BillingPlanCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BillingPlanCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *BillingPlanCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *BillingPlanCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *BillingPlanCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *BillingPlanCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetPrice

`func (o *BillingPlanCreateRequest) GetPrice() float32`

GetPrice returns the Price field if non-nil, zero value otherwise.

### GetPriceOk

`func (o *BillingPlanCreateRequest) GetPriceOk() (*float32, bool)`

GetPriceOk returns a tuple with the Price field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrice

`func (o *BillingPlanCreateRequest) SetPrice(v float32)`

SetPrice sets Price field to given value.


### GetNotificationLimit

`func (o *BillingPlanCreateRequest) GetNotificationLimit() int32`

GetNotificationLimit returns the NotificationLimit field if non-nil, zero value otherwise.

### GetNotificationLimitOk

`func (o *BillingPlanCreateRequest) GetNotificationLimitOk() (*int32, bool)`

GetNotificationLimitOk returns a tuple with the NotificationLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotificationLimit

`func (o *BillingPlanCreateRequest) SetNotificationLimit(v int32)`

SetNotificationLimit sets NotificationLimit field to given value.


### GetFeatures

`func (o *BillingPlanCreateRequest) GetFeatures() []string`

GetFeatures returns the Features field if non-nil, zero value otherwise.

### GetFeaturesOk

`func (o *BillingPlanCreateRequest) GetFeaturesOk() (*[]string, bool)`

GetFeaturesOk returns a tuple with the Features field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeatures

`func (o *BillingPlanCreateRequest) SetFeatures(v []string)`

SetFeatures sets Features field to given value.

### HasFeatures

`func (o *BillingPlanCreateRequest) HasFeatures() bool`

HasFeatures returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


