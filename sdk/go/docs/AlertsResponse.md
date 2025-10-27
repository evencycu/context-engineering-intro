# AlertsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Alerts** | Pointer to [**[]AlertsResponseAlertsInner**](AlertsResponseAlertsInner.md) |  | [optional] 
**TotalAlerts** | Pointer to **int32** |  | [optional] 
**CriticalAlerts** | Pointer to **int32** |  | [optional] 

## Methods

### NewAlertsResponse

`func NewAlertsResponse() *AlertsResponse`

NewAlertsResponse instantiates a new AlertsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAlertsResponseWithDefaults

`func NewAlertsResponseWithDefaults() *AlertsResponse`

NewAlertsResponseWithDefaults instantiates a new AlertsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAlerts

`func (o *AlertsResponse) GetAlerts() []AlertsResponseAlertsInner`

GetAlerts returns the Alerts field if non-nil, zero value otherwise.

### GetAlertsOk

`func (o *AlertsResponse) GetAlertsOk() (*[]AlertsResponseAlertsInner, bool)`

GetAlertsOk returns a tuple with the Alerts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlerts

`func (o *AlertsResponse) SetAlerts(v []AlertsResponseAlertsInner)`

SetAlerts sets Alerts field to given value.

### HasAlerts

`func (o *AlertsResponse) HasAlerts() bool`

HasAlerts returns a boolean if a field has been set.

### GetTotalAlerts

`func (o *AlertsResponse) GetTotalAlerts() int32`

GetTotalAlerts returns the TotalAlerts field if non-nil, zero value otherwise.

### GetTotalAlertsOk

`func (o *AlertsResponse) GetTotalAlertsOk() (*int32, bool)`

GetTotalAlertsOk returns a tuple with the TotalAlerts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalAlerts

`func (o *AlertsResponse) SetTotalAlerts(v int32)`

SetTotalAlerts sets TotalAlerts field to given value.

### HasTotalAlerts

`func (o *AlertsResponse) HasTotalAlerts() bool`

HasTotalAlerts returns a boolean if a field has been set.

### GetCriticalAlerts

`func (o *AlertsResponse) GetCriticalAlerts() int32`

GetCriticalAlerts returns the CriticalAlerts field if non-nil, zero value otherwise.

### GetCriticalAlertsOk

`func (o *AlertsResponse) GetCriticalAlertsOk() (*int32, bool)`

GetCriticalAlertsOk returns a tuple with the CriticalAlerts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCriticalAlerts

`func (o *AlertsResponse) SetCriticalAlerts(v int32)`

SetCriticalAlerts sets CriticalAlerts field to given value.

### HasCriticalAlerts

`func (o *AlertsResponse) HasCriticalAlerts() bool`

HasCriticalAlerts returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


