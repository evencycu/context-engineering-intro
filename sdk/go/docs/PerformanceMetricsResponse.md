# PerformanceMetricsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Timestamp** | Pointer to **time.Time** |  | [optional] 
**Cpu** | Pointer to [**PerformanceMetricsResponseCpu**](PerformanceMetricsResponseCpu.md) |  | [optional] 
**Memory** | Pointer to [**PerformanceMetricsResponseMemory**](PerformanceMetricsResponseMemory.md) |  | [optional] 
**Database** | Pointer to [**PerformanceMetricsResponseDatabase**](PerformanceMetricsResponseDatabase.md) |  | [optional] 
**Redis** | Pointer to [**PerformanceMetricsResponseRedis**](PerformanceMetricsResponseRedis.md) |  | [optional] 

## Methods

### NewPerformanceMetricsResponse

`func NewPerformanceMetricsResponse() *PerformanceMetricsResponse`

NewPerformanceMetricsResponse instantiates a new PerformanceMetricsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPerformanceMetricsResponseWithDefaults

`func NewPerformanceMetricsResponseWithDefaults() *PerformanceMetricsResponse`

NewPerformanceMetricsResponseWithDefaults instantiates a new PerformanceMetricsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTimestamp

`func (o *PerformanceMetricsResponse) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *PerformanceMetricsResponse) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *PerformanceMetricsResponse) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *PerformanceMetricsResponse) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetCpu

`func (o *PerformanceMetricsResponse) GetCpu() PerformanceMetricsResponseCpu`

GetCpu returns the Cpu field if non-nil, zero value otherwise.

### GetCpuOk

`func (o *PerformanceMetricsResponse) GetCpuOk() (*PerformanceMetricsResponseCpu, bool)`

GetCpuOk returns a tuple with the Cpu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpu

`func (o *PerformanceMetricsResponse) SetCpu(v PerformanceMetricsResponseCpu)`

SetCpu sets Cpu field to given value.

### HasCpu

`func (o *PerformanceMetricsResponse) HasCpu() bool`

HasCpu returns a boolean if a field has been set.

### GetMemory

`func (o *PerformanceMetricsResponse) GetMemory() PerformanceMetricsResponseMemory`

GetMemory returns the Memory field if non-nil, zero value otherwise.

### GetMemoryOk

`func (o *PerformanceMetricsResponse) GetMemoryOk() (*PerformanceMetricsResponseMemory, bool)`

GetMemoryOk returns a tuple with the Memory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemory

`func (o *PerformanceMetricsResponse) SetMemory(v PerformanceMetricsResponseMemory)`

SetMemory sets Memory field to given value.

### HasMemory

`func (o *PerformanceMetricsResponse) HasMemory() bool`

HasMemory returns a boolean if a field has been set.

### GetDatabase

`func (o *PerformanceMetricsResponse) GetDatabase() PerformanceMetricsResponseDatabase`

GetDatabase returns the Database field if non-nil, zero value otherwise.

### GetDatabaseOk

`func (o *PerformanceMetricsResponse) GetDatabaseOk() (*PerformanceMetricsResponseDatabase, bool)`

GetDatabaseOk returns a tuple with the Database field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDatabase

`func (o *PerformanceMetricsResponse) SetDatabase(v PerformanceMetricsResponseDatabase)`

SetDatabase sets Database field to given value.

### HasDatabase

`func (o *PerformanceMetricsResponse) HasDatabase() bool`

HasDatabase returns a boolean if a field has been set.

### GetRedis

`func (o *PerformanceMetricsResponse) GetRedis() PerformanceMetricsResponseRedis`

GetRedis returns the Redis field if non-nil, zero value otherwise.

### GetRedisOk

`func (o *PerformanceMetricsResponse) GetRedisOk() (*PerformanceMetricsResponseRedis, bool)`

GetRedisOk returns a tuple with the Redis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedis

`func (o *PerformanceMetricsResponse) SetRedis(v PerformanceMetricsResponseRedis)`

SetRedis sets Redis field to given value.

### HasRedis

`func (o *PerformanceMetricsResponse) HasRedis() bool`

HasRedis returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


