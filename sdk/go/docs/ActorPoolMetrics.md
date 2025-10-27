# ActorPoolMetrics

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ActiveActors** | Pointer to **int32** |  | [optional] 
**IdleActors** | Pointer to **int32** |  | [optional] 
**TotalActors** | Pointer to **int32** |  | [optional] 
**QueueSize** | Pointer to **int32** |  | [optional] 

## Methods

### NewActorPoolMetrics

`func NewActorPoolMetrics() *ActorPoolMetrics`

NewActorPoolMetrics instantiates a new ActorPoolMetrics object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewActorPoolMetricsWithDefaults

`func NewActorPoolMetricsWithDefaults() *ActorPoolMetrics`

NewActorPoolMetricsWithDefaults instantiates a new ActorPoolMetrics object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActiveActors

`func (o *ActorPoolMetrics) GetActiveActors() int32`

GetActiveActors returns the ActiveActors field if non-nil, zero value otherwise.

### GetActiveActorsOk

`func (o *ActorPoolMetrics) GetActiveActorsOk() (*int32, bool)`

GetActiveActorsOk returns a tuple with the ActiveActors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveActors

`func (o *ActorPoolMetrics) SetActiveActors(v int32)`

SetActiveActors sets ActiveActors field to given value.

### HasActiveActors

`func (o *ActorPoolMetrics) HasActiveActors() bool`

HasActiveActors returns a boolean if a field has been set.

### GetIdleActors

`func (o *ActorPoolMetrics) GetIdleActors() int32`

GetIdleActors returns the IdleActors field if non-nil, zero value otherwise.

### GetIdleActorsOk

`func (o *ActorPoolMetrics) GetIdleActorsOk() (*int32, bool)`

GetIdleActorsOk returns a tuple with the IdleActors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdleActors

`func (o *ActorPoolMetrics) SetIdleActors(v int32)`

SetIdleActors sets IdleActors field to given value.

### HasIdleActors

`func (o *ActorPoolMetrics) HasIdleActors() bool`

HasIdleActors returns a boolean if a field has been set.

### GetTotalActors

`func (o *ActorPoolMetrics) GetTotalActors() int32`

GetTotalActors returns the TotalActors field if non-nil, zero value otherwise.

### GetTotalActorsOk

`func (o *ActorPoolMetrics) GetTotalActorsOk() (*int32, bool)`

GetTotalActorsOk returns a tuple with the TotalActors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalActors

`func (o *ActorPoolMetrics) SetTotalActors(v int32)`

SetTotalActors sets TotalActors field to given value.

### HasTotalActors

`func (o *ActorPoolMetrics) HasTotalActors() bool`

HasTotalActors returns a boolean if a field has been set.

### GetQueueSize

`func (o *ActorPoolMetrics) GetQueueSize() int32`

GetQueueSize returns the QueueSize field if non-nil, zero value otherwise.

### GetQueueSizeOk

`func (o *ActorPoolMetrics) GetQueueSizeOk() (*int32, bool)`

GetQueueSizeOk returns a tuple with the QueueSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueueSize

`func (o *ActorPoolMetrics) SetQueueSize(v int32)`

SetQueueSize sets QueueSize field to given value.

### HasQueueSize

`func (o *ActorPoolMetrics) HasQueueSize() bool`

HasQueueSize returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


