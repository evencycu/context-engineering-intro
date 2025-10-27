# SystemMetrics

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MemoryUsage** | Pointer to **float32** |  | [optional] 
**CpuUsage** | Pointer to **float32** |  | [optional] 
**Goroutines** | Pointer to **int32** |  | [optional] 
**Uptime** | Pointer to **string** |  | [optional] 

## Methods

### NewSystemMetrics

`func NewSystemMetrics() *SystemMetrics`

NewSystemMetrics instantiates a new SystemMetrics object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSystemMetricsWithDefaults

`func NewSystemMetricsWithDefaults() *SystemMetrics`

NewSystemMetricsWithDefaults instantiates a new SystemMetrics object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMemoryUsage

`func (o *SystemMetrics) GetMemoryUsage() float32`

GetMemoryUsage returns the MemoryUsage field if non-nil, zero value otherwise.

### GetMemoryUsageOk

`func (o *SystemMetrics) GetMemoryUsageOk() (*float32, bool)`

GetMemoryUsageOk returns a tuple with the MemoryUsage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemoryUsage

`func (o *SystemMetrics) SetMemoryUsage(v float32)`

SetMemoryUsage sets MemoryUsage field to given value.

### HasMemoryUsage

`func (o *SystemMetrics) HasMemoryUsage() bool`

HasMemoryUsage returns a boolean if a field has been set.

### GetCpuUsage

`func (o *SystemMetrics) GetCpuUsage() float32`

GetCpuUsage returns the CpuUsage field if non-nil, zero value otherwise.

### GetCpuUsageOk

`func (o *SystemMetrics) GetCpuUsageOk() (*float32, bool)`

GetCpuUsageOk returns a tuple with the CpuUsage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpuUsage

`func (o *SystemMetrics) SetCpuUsage(v float32)`

SetCpuUsage sets CpuUsage field to given value.

### HasCpuUsage

`func (o *SystemMetrics) HasCpuUsage() bool`

HasCpuUsage returns a boolean if a field has been set.

### GetGoroutines

`func (o *SystemMetrics) GetGoroutines() int32`

GetGoroutines returns the Goroutines field if non-nil, zero value otherwise.

### GetGoroutinesOk

`func (o *SystemMetrics) GetGoroutinesOk() (*int32, bool)`

GetGoroutinesOk returns a tuple with the Goroutines field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGoroutines

`func (o *SystemMetrics) SetGoroutines(v int32)`

SetGoroutines sets Goroutines field to given value.

### HasGoroutines

`func (o *SystemMetrics) HasGoroutines() bool`

HasGoroutines returns a boolean if a field has been set.

### GetUptime

`func (o *SystemMetrics) GetUptime() string`

GetUptime returns the Uptime field if non-nil, zero value otherwise.

### GetUptimeOk

`func (o *SystemMetrics) GetUptimeOk() (*string, bool)`

GetUptimeOk returns a tuple with the Uptime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUptime

`func (o *SystemMetrics) SetUptime(v string)`

SetUptime sets Uptime field to given value.

### HasUptime

`func (o *SystemMetrics) HasUptime() bool`

HasUptime returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


