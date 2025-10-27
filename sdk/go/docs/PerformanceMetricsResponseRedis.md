# PerformanceMetricsResponseRedis

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Memory** | Pointer to **int32** |  | [optional] 
**Keys** | Pointer to **int32** |  | [optional] 
**HitRate** | Pointer to **float32** |  | [optional] 

## Methods

### NewPerformanceMetricsResponseRedis

`func NewPerformanceMetricsResponseRedis() *PerformanceMetricsResponseRedis`

NewPerformanceMetricsResponseRedis instantiates a new PerformanceMetricsResponseRedis object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPerformanceMetricsResponseRedisWithDefaults

`func NewPerformanceMetricsResponseRedisWithDefaults() *PerformanceMetricsResponseRedis`

NewPerformanceMetricsResponseRedisWithDefaults instantiates a new PerformanceMetricsResponseRedis object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMemory

`func (o *PerformanceMetricsResponseRedis) GetMemory() int32`

GetMemory returns the Memory field if non-nil, zero value otherwise.

### GetMemoryOk

`func (o *PerformanceMetricsResponseRedis) GetMemoryOk() (*int32, bool)`

GetMemoryOk returns a tuple with the Memory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemory

`func (o *PerformanceMetricsResponseRedis) SetMemory(v int32)`

SetMemory sets Memory field to given value.

### HasMemory

`func (o *PerformanceMetricsResponseRedis) HasMemory() bool`

HasMemory returns a boolean if a field has been set.

### GetKeys

`func (o *PerformanceMetricsResponseRedis) GetKeys() int32`

GetKeys returns the Keys field if non-nil, zero value otherwise.

### GetKeysOk

`func (o *PerformanceMetricsResponseRedis) GetKeysOk() (*int32, bool)`

GetKeysOk returns a tuple with the Keys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeys

`func (o *PerformanceMetricsResponseRedis) SetKeys(v int32)`

SetKeys sets Keys field to given value.

### HasKeys

`func (o *PerformanceMetricsResponseRedis) HasKeys() bool`

HasKeys returns a boolean if a field has been set.

### GetHitRate

`func (o *PerformanceMetricsResponseRedis) GetHitRate() float32`

GetHitRate returns the HitRate field if non-nil, zero value otherwise.

### GetHitRateOk

`func (o *PerformanceMetricsResponseRedis) GetHitRateOk() (*float32, bool)`

GetHitRateOk returns a tuple with the HitRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHitRate

`func (o *PerformanceMetricsResponseRedis) SetHitRate(v float32)`

SetHitRate sets HitRate field to given value.

### HasHitRate

`func (o *PerformanceMetricsResponseRedis) HasHitRate() bool`

HasHitRate returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


