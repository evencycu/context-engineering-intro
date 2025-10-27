# BillingAnalyticsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TotalRevenue** | Pointer to **float32** |  | [optional] 
**ActiveCompanies** | Pointer to **int32** |  | [optional] 
**AverageUsage** | Pointer to **float32** |  | [optional] 
**TopPlans** | Pointer to [**[]BillingAnalyticsResponseTopPlansInner**](BillingAnalyticsResponseTopPlansInner.md) |  | [optional] 

## Methods

### NewBillingAnalyticsResponse

`func NewBillingAnalyticsResponse() *BillingAnalyticsResponse`

NewBillingAnalyticsResponse instantiates a new BillingAnalyticsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBillingAnalyticsResponseWithDefaults

`func NewBillingAnalyticsResponseWithDefaults() *BillingAnalyticsResponse`

NewBillingAnalyticsResponseWithDefaults instantiates a new BillingAnalyticsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTotalRevenue

`func (o *BillingAnalyticsResponse) GetTotalRevenue() float32`

GetTotalRevenue returns the TotalRevenue field if non-nil, zero value otherwise.

### GetTotalRevenueOk

`func (o *BillingAnalyticsResponse) GetTotalRevenueOk() (*float32, bool)`

GetTotalRevenueOk returns a tuple with the TotalRevenue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalRevenue

`func (o *BillingAnalyticsResponse) SetTotalRevenue(v float32)`

SetTotalRevenue sets TotalRevenue field to given value.

### HasTotalRevenue

`func (o *BillingAnalyticsResponse) HasTotalRevenue() bool`

HasTotalRevenue returns a boolean if a field has been set.

### GetActiveCompanies

`func (o *BillingAnalyticsResponse) GetActiveCompanies() int32`

GetActiveCompanies returns the ActiveCompanies field if non-nil, zero value otherwise.

### GetActiveCompaniesOk

`func (o *BillingAnalyticsResponse) GetActiveCompaniesOk() (*int32, bool)`

GetActiveCompaniesOk returns a tuple with the ActiveCompanies field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveCompanies

`func (o *BillingAnalyticsResponse) SetActiveCompanies(v int32)`

SetActiveCompanies sets ActiveCompanies field to given value.

### HasActiveCompanies

`func (o *BillingAnalyticsResponse) HasActiveCompanies() bool`

HasActiveCompanies returns a boolean if a field has been set.

### GetAverageUsage

`func (o *BillingAnalyticsResponse) GetAverageUsage() float32`

GetAverageUsage returns the AverageUsage field if non-nil, zero value otherwise.

### GetAverageUsageOk

`func (o *BillingAnalyticsResponse) GetAverageUsageOk() (*float32, bool)`

GetAverageUsageOk returns a tuple with the AverageUsage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAverageUsage

`func (o *BillingAnalyticsResponse) SetAverageUsage(v float32)`

SetAverageUsage sets AverageUsage field to given value.

### HasAverageUsage

`func (o *BillingAnalyticsResponse) HasAverageUsage() bool`

HasAverageUsage returns a boolean if a field has been set.

### GetTopPlans

`func (o *BillingAnalyticsResponse) GetTopPlans() []BillingAnalyticsResponseTopPlansInner`

GetTopPlans returns the TopPlans field if non-nil, zero value otherwise.

### GetTopPlansOk

`func (o *BillingAnalyticsResponse) GetTopPlansOk() (*[]BillingAnalyticsResponseTopPlansInner, bool)`

GetTopPlansOk returns a tuple with the TopPlans field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTopPlans

`func (o *BillingAnalyticsResponse) SetTopPlans(v []BillingAnalyticsResponseTopPlansInner)`

SetTopPlans sets TopPlans field to given value.

### HasTopPlans

`func (o *BillingAnalyticsResponse) HasTopPlans() bool`

HasTopPlans returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


