# DashboardResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Timestamp** | Pointer to **time.Time** |  | [optional] 
**Overview** | Pointer to [**BusinessHealthResponse**](BusinessHealthResponse.md) |  | [optional] 
**Performance** | Pointer to [**PerformanceMetricsResponse**](PerformanceMetricsResponse.md) |  | [optional] 
**Alerts** | Pointer to [**AlertsResponse**](AlertsResponse.md) |  | [optional] 
**Charts** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewDashboardResponse

`func NewDashboardResponse() *DashboardResponse`

NewDashboardResponse instantiates a new DashboardResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDashboardResponseWithDefaults

`func NewDashboardResponseWithDefaults() *DashboardResponse`

NewDashboardResponseWithDefaults instantiates a new DashboardResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTimestamp

`func (o *DashboardResponse) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *DashboardResponse) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *DashboardResponse) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *DashboardResponse) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetOverview

`func (o *DashboardResponse) GetOverview() BusinessHealthResponse`

GetOverview returns the Overview field if non-nil, zero value otherwise.

### GetOverviewOk

`func (o *DashboardResponse) GetOverviewOk() (*BusinessHealthResponse, bool)`

GetOverviewOk returns a tuple with the Overview field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOverview

`func (o *DashboardResponse) SetOverview(v BusinessHealthResponse)`

SetOverview sets Overview field to given value.

### HasOverview

`func (o *DashboardResponse) HasOverview() bool`

HasOverview returns a boolean if a field has been set.

### GetPerformance

`func (o *DashboardResponse) GetPerformance() PerformanceMetricsResponse`

GetPerformance returns the Performance field if non-nil, zero value otherwise.

### GetPerformanceOk

`func (o *DashboardResponse) GetPerformanceOk() (*PerformanceMetricsResponse, bool)`

GetPerformanceOk returns a tuple with the Performance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPerformance

`func (o *DashboardResponse) SetPerformance(v PerformanceMetricsResponse)`

SetPerformance sets Performance field to given value.

### HasPerformance

`func (o *DashboardResponse) HasPerformance() bool`

HasPerformance returns a boolean if a field has been set.

### GetAlerts

`func (o *DashboardResponse) GetAlerts() AlertsResponse`

GetAlerts returns the Alerts field if non-nil, zero value otherwise.

### GetAlertsOk

`func (o *DashboardResponse) GetAlertsOk() (*AlertsResponse, bool)`

GetAlertsOk returns a tuple with the Alerts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlerts

`func (o *DashboardResponse) SetAlerts(v AlertsResponse)`

SetAlerts sets Alerts field to given value.

### HasAlerts

`func (o *DashboardResponse) HasAlerts() bool`

HasAlerts returns a boolean if a field has been set.

### GetCharts

`func (o *DashboardResponse) GetCharts() map[string]interface{}`

GetCharts returns the Charts field if non-nil, zero value otherwise.

### GetChartsOk

`func (o *DashboardResponse) GetChartsOk() (*map[string]interface{}, bool)`

GetChartsOk returns a tuple with the Charts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCharts

`func (o *DashboardResponse) SetCharts(v map[string]interface{})`

SetCharts sets Charts field to given value.

### HasCharts

`func (o *DashboardResponse) HasCharts() bool`

HasCharts returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


