# BillingPlanResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**NotificationLimit** | Pointer to **int32** |  | [optional] 
**Features** | Pointer to **[]string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewBillingPlanResponse

`func NewBillingPlanResponse() *BillingPlanResponse`

NewBillingPlanResponse instantiates a new BillingPlanResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBillingPlanResponseWithDefaults

`func NewBillingPlanResponseWithDefaults() *BillingPlanResponse`

NewBillingPlanResponseWithDefaults instantiates a new BillingPlanResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *BillingPlanResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *BillingPlanResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *BillingPlanResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *BillingPlanResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *BillingPlanResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BillingPlanResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BillingPlanResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *BillingPlanResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *BillingPlanResponse) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *BillingPlanResponse) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *BillingPlanResponse) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *BillingPlanResponse) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetPrice

`func (o *BillingPlanResponse) GetPrice() float32`

GetPrice returns the Price field if non-nil, zero value otherwise.

### GetPriceOk

`func (o *BillingPlanResponse) GetPriceOk() (*float32, bool)`

GetPriceOk returns a tuple with the Price field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrice

`func (o *BillingPlanResponse) SetPrice(v float32)`

SetPrice sets Price field to given value.

### HasPrice

`func (o *BillingPlanResponse) HasPrice() bool`

HasPrice returns a boolean if a field has been set.

### GetNotificationLimit

`func (o *BillingPlanResponse) GetNotificationLimit() int32`

GetNotificationLimit returns the NotificationLimit field if non-nil, zero value otherwise.

### GetNotificationLimitOk

`func (o *BillingPlanResponse) GetNotificationLimitOk() (*int32, bool)`

GetNotificationLimitOk returns a tuple with the NotificationLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotificationLimit

`func (o *BillingPlanResponse) SetNotificationLimit(v int32)`

SetNotificationLimit sets NotificationLimit field to given value.

### HasNotificationLimit

`func (o *BillingPlanResponse) HasNotificationLimit() bool`

HasNotificationLimit returns a boolean if a field has been set.

### GetFeatures

`func (o *BillingPlanResponse) GetFeatures() []string`

GetFeatures returns the Features field if non-nil, zero value otherwise.

### GetFeaturesOk

`func (o *BillingPlanResponse) GetFeaturesOk() (*[]string, bool)`

GetFeaturesOk returns a tuple with the Features field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeatures

`func (o *BillingPlanResponse) SetFeatures(v []string)`

SetFeatures sets Features field to given value.

### HasFeatures

`func (o *BillingPlanResponse) HasFeatures() bool`

HasFeatures returns a boolean if a field has been set.

### GetCreatedAt

`func (o *BillingPlanResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *BillingPlanResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *BillingPlanResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *BillingPlanResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *BillingPlanResponse) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *BillingPlanResponse) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *BillingPlanResponse) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *BillingPlanResponse) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


