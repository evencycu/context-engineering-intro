# QueueStatusResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | Pointer to **string** |  | [optional] 
**Workers** | Pointer to **int32** |  | [optional] 
**ActiveWorkers** | Pointer to **int32** |  | [optional] 
**QueueSize** | Pointer to **int32** |  | [optional] 
**LastProcessed** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewQueueStatusResponse

`func NewQueueStatusResponse() *QueueStatusResponse`

NewQueueStatusResponse instantiates a new QueueStatusResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQueueStatusResponseWithDefaults

`func NewQueueStatusResponseWithDefaults() *QueueStatusResponse`

NewQueueStatusResponseWithDefaults instantiates a new QueueStatusResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *QueueStatusResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *QueueStatusResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *QueueStatusResponse) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *QueueStatusResponse) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetWorkers

`func (o *QueueStatusResponse) GetWorkers() int32`

GetWorkers returns the Workers field if non-nil, zero value otherwise.

### GetWorkersOk

`func (o *QueueStatusResponse) GetWorkersOk() (*int32, bool)`

GetWorkersOk returns a tuple with the Workers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkers

`func (o *QueueStatusResponse) SetWorkers(v int32)`

SetWorkers sets Workers field to given value.

### HasWorkers

`func (o *QueueStatusResponse) HasWorkers() bool`

HasWorkers returns a boolean if a field has been set.

### GetActiveWorkers

`func (o *QueueStatusResponse) GetActiveWorkers() int32`

GetActiveWorkers returns the ActiveWorkers field if non-nil, zero value otherwise.

### GetActiveWorkersOk

`func (o *QueueStatusResponse) GetActiveWorkersOk() (*int32, bool)`

GetActiveWorkersOk returns a tuple with the ActiveWorkers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActiveWorkers

`func (o *QueueStatusResponse) SetActiveWorkers(v int32)`

SetActiveWorkers sets ActiveWorkers field to given value.

### HasActiveWorkers

`func (o *QueueStatusResponse) HasActiveWorkers() bool`

HasActiveWorkers returns a boolean if a field has been set.

### GetQueueSize

`func (o *QueueStatusResponse) GetQueueSize() int32`

GetQueueSize returns the QueueSize field if non-nil, zero value otherwise.

### GetQueueSizeOk

`func (o *QueueStatusResponse) GetQueueSizeOk() (*int32, bool)`

GetQueueSizeOk returns a tuple with the QueueSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueueSize

`func (o *QueueStatusResponse) SetQueueSize(v int32)`

SetQueueSize sets QueueSize field to given value.

### HasQueueSize

`func (o *QueueStatusResponse) HasQueueSize() bool`

HasQueueSize returns a boolean if a field has been set.

### GetLastProcessed

`func (o *QueueStatusResponse) GetLastProcessed() time.Time`

GetLastProcessed returns the LastProcessed field if non-nil, zero value otherwise.

### GetLastProcessedOk

`func (o *QueueStatusResponse) GetLastProcessedOk() (*time.Time, bool)`

GetLastProcessedOk returns a tuple with the LastProcessed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastProcessed

`func (o *QueueStatusResponse) SetLastProcessed(v time.Time)`

SetLastProcessed sets LastProcessed field to given value.

### HasLastProcessed

`func (o *QueueStatusResponse) HasLastProcessed() bool`

HasLastProcessed returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


