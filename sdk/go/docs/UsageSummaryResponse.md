# UsageSummaryResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TotalNotifications** | Pointer to **int32** |  | [optional] 
**DailyAverage** | Pointer to **float32** |  | [optional] 
**MonthlyTotal** | Pointer to **int32** |  | [optional] 
**Period** | Pointer to [**UsageSummaryResponsePeriod**](UsageSummaryResponsePeriod.md) |  | [optional] 

## Methods

### NewUsageSummaryResponse

`func NewUsageSummaryResponse() *UsageSummaryResponse`

NewUsageSummaryResponse instantiates a new UsageSummaryResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUsageSummaryResponseWithDefaults

`func NewUsageSummaryResponseWithDefaults() *UsageSummaryResponse`

NewUsageSummaryResponseWithDefaults instantiates a new UsageSummaryResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTotalNotifications

`func (o *UsageSummaryResponse) GetTotalNotifications() int32`

GetTotalNotifications returns the TotalNotifications field if non-nil, zero value otherwise.

### GetTotalNotificationsOk

`func (o *UsageSummaryResponse) GetTotalNotificationsOk() (*int32, bool)`

GetTotalNotificationsOk returns a tuple with the TotalNotifications field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalNotifications

`func (o *UsageSummaryResponse) SetTotalNotifications(v int32)`

SetTotalNotifications sets TotalNotifications field to given value.

### HasTotalNotifications

`func (o *UsageSummaryResponse) HasTotalNotifications() bool`

HasTotalNotifications returns a boolean if a field has been set.

### GetDailyAverage

`func (o *UsageSummaryResponse) GetDailyAverage() float32`

GetDailyAverage returns the DailyAverage field if non-nil, zero value otherwise.

### GetDailyAverageOk

`func (o *UsageSummaryResponse) GetDailyAverageOk() (*float32, bool)`

GetDailyAverageOk returns a tuple with the DailyAverage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDailyAverage

`func (o *UsageSummaryResponse) SetDailyAverage(v float32)`

SetDailyAverage sets DailyAverage field to given value.

### HasDailyAverage

`func (o *UsageSummaryResponse) HasDailyAverage() bool`

HasDailyAverage returns a boolean if a field has been set.

### GetMonthlyTotal

`func (o *UsageSummaryResponse) GetMonthlyTotal() int32`

GetMonthlyTotal returns the MonthlyTotal field if non-nil, zero value otherwise.

### GetMonthlyTotalOk

`func (o *UsageSummaryResponse) GetMonthlyTotalOk() (*int32, bool)`

GetMonthlyTotalOk returns a tuple with the MonthlyTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonthlyTotal

`func (o *UsageSummaryResponse) SetMonthlyTotal(v int32)`

SetMonthlyTotal sets MonthlyTotal field to given value.

### HasMonthlyTotal

`func (o *UsageSummaryResponse) HasMonthlyTotal() bool`

HasMonthlyTotal returns a boolean if a field has been set.

### GetPeriod

`func (o *UsageSummaryResponse) GetPeriod() UsageSummaryResponsePeriod`

GetPeriod returns the Period field if non-nil, zero value otherwise.

### GetPeriodOk

`func (o *UsageSummaryResponse) GetPeriodOk() (*UsageSummaryResponsePeriod, bool)`

GetPeriodOk returns a tuple with the Period field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPeriod

`func (o *UsageSummaryResponse) SetPeriod(v UsageSummaryResponsePeriod)`

SetPeriod sets Period field to given value.

### HasPeriod

`func (o *UsageSummaryResponse) HasPeriod() bool`

HasPeriod returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


