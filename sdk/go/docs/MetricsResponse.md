# MetricsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Timestamp** | Pointer to **time.Time** |  | [optional] 
**System** | Pointer to [**SystemMetrics**](SystemMetrics.md) |  | [optional] 
**TokenCache** | Pointer to [**TokenCacheMetrics**](TokenCacheMetrics.md) |  | [optional] 
**ActorPool** | Pointer to [**ActorPoolMetrics**](ActorPoolMetrics.md) |  | [optional] 

## Methods

### NewMetricsResponse

`func NewMetricsResponse() *MetricsResponse`

NewMetricsResponse instantiates a new MetricsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMetricsResponseWithDefaults

`func NewMetricsResponseWithDefaults() *MetricsResponse`

NewMetricsResponseWithDefaults instantiates a new MetricsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTimestamp

`func (o *MetricsResponse) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *MetricsResponse) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *MetricsResponse) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *MetricsResponse) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetSystem

`func (o *MetricsResponse) GetSystem() SystemMetrics`

GetSystem returns the System field if non-nil, zero value otherwise.

### GetSystemOk

`func (o *MetricsResponse) GetSystemOk() (*SystemMetrics, bool)`

GetSystemOk returns a tuple with the System field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSystem

`func (o *MetricsResponse) SetSystem(v SystemMetrics)`

SetSystem sets System field to given value.

### HasSystem

`func (o *MetricsResponse) HasSystem() bool`

HasSystem returns a boolean if a field has been set.

### GetTokenCache

`func (o *MetricsResponse) GetTokenCache() TokenCacheMetrics`

GetTokenCache returns the TokenCache field if non-nil, zero value otherwise.

### GetTokenCacheOk

`func (o *MetricsResponse) GetTokenCacheOk() (*TokenCacheMetrics, bool)`

GetTokenCacheOk returns a tuple with the TokenCache field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenCache

`func (o *MetricsResponse) SetTokenCache(v TokenCacheMetrics)`

SetTokenCache sets TokenCache field to given value.

### HasTokenCache

`func (o *MetricsResponse) HasTokenCache() bool`

HasTokenCache returns a boolean if a field has been set.

### GetActorPool

`func (o *MetricsResponse) GetActorPool() ActorPoolMetrics`

GetActorPool returns the ActorPool field if non-nil, zero value otherwise.

### GetActorPoolOk

`func (o *MetricsResponse) GetActorPoolOk() (*ActorPoolMetrics, bool)`

GetActorPoolOk returns a tuple with the ActorPool field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActorPool

`func (o *MetricsResponse) SetActorPool(v ActorPoolMetrics)`

SetActorPool sets ActorPool field to given value.

### HasActorPool

`func (o *MetricsResponse) HasActorPool() bool`

HasActorPool returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


