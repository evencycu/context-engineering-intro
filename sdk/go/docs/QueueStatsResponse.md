# QueueStatsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Timestamp** | Pointer to **time.Time** |  | [optional] 
**TotalItems** | Pointer to **int32** |  | [optional] 
**PendingItems** | Pointer to **int32** |  | [optional] 
**ProcessingItems** | Pointer to **int32** |  | [optional] 
**CompletedItems** | Pointer to **int32** |  | [optional] 
**FailedItems** | Pointer to **int32** |  | [optional] 
**AverageProcessingTime** | Pointer to **float32** |  | [optional] 
**Throughput** | Pointer to **float32** |  | [optional] 

## Methods

### NewQueueStatsResponse

`func NewQueueStatsResponse() *QueueStatsResponse`

NewQueueStatsResponse instantiates a new QueueStatsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQueueStatsResponseWithDefaults

`func NewQueueStatsResponseWithDefaults() *QueueStatsResponse`

NewQueueStatsResponseWithDefaults instantiates a new QueueStatsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTimestamp

`func (o *QueueStatsResponse) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *QueueStatsResponse) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *QueueStatsResponse) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *QueueStatsResponse) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetTotalItems

`func (o *QueueStatsResponse) GetTotalItems() int32`

GetTotalItems returns the TotalItems field if non-nil, zero value otherwise.

### GetTotalItemsOk

`func (o *QueueStatsResponse) GetTotalItemsOk() (*int32, bool)`

GetTotalItemsOk returns a tuple with the TotalItems field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalItems

`func (o *QueueStatsResponse) SetTotalItems(v int32)`

SetTotalItems sets TotalItems field to given value.

### HasTotalItems

`func (o *QueueStatsResponse) HasTotalItems() bool`

HasTotalItems returns a boolean if a field has been set.

### GetPendingItems

`func (o *QueueStatsResponse) GetPendingItems() int32`

GetPendingItems returns the PendingItems field if non-nil, zero value otherwise.

### GetPendingItemsOk

`func (o *QueueStatsResponse) GetPendingItemsOk() (*int32, bool)`

GetPendingItemsOk returns a tuple with the PendingItems field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPendingItems

`func (o *QueueStatsResponse) SetPendingItems(v int32)`

SetPendingItems sets PendingItems field to given value.

### HasPendingItems

`func (o *QueueStatsResponse) HasPendingItems() bool`

HasPendingItems returns a boolean if a field has been set.

### GetProcessingItems

`func (o *QueueStatsResponse) GetProcessingItems() int32`

GetProcessingItems returns the ProcessingItems field if non-nil, zero value otherwise.

### GetProcessingItemsOk

`func (o *QueueStatsResponse) GetProcessingItemsOk() (*int32, bool)`

GetProcessingItemsOk returns a tuple with the ProcessingItems field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProcessingItems

`func (o *QueueStatsResponse) SetProcessingItems(v int32)`

SetProcessingItems sets ProcessingItems field to given value.

### HasProcessingItems

`func (o *QueueStatsResponse) HasProcessingItems() bool`

HasProcessingItems returns a boolean if a field has been set.

### GetCompletedItems

`func (o *QueueStatsResponse) GetCompletedItems() int32`

GetCompletedItems returns the CompletedItems field if non-nil, zero value otherwise.

### GetCompletedItemsOk

`func (o *QueueStatsResponse) GetCompletedItemsOk() (*int32, bool)`

GetCompletedItemsOk returns a tuple with the CompletedItems field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedItems

`func (o *QueueStatsResponse) SetCompletedItems(v int32)`

SetCompletedItems sets CompletedItems field to given value.

### HasCompletedItems

`func (o *QueueStatsResponse) HasCompletedItems() bool`

HasCompletedItems returns a boolean if a field has been set.

### GetFailedItems

`func (o *QueueStatsResponse) GetFailedItems() int32`

GetFailedItems returns the FailedItems field if non-nil, zero value otherwise.

### GetFailedItemsOk

`func (o *QueueStatsResponse) GetFailedItemsOk() (*int32, bool)`

GetFailedItemsOk returns a tuple with the FailedItems field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailedItems

`func (o *QueueStatsResponse) SetFailedItems(v int32)`

SetFailedItems sets FailedItems field to given value.

### HasFailedItems

`func (o *QueueStatsResponse) HasFailedItems() bool`

HasFailedItems returns a boolean if a field has been set.

### GetAverageProcessingTime

`func (o *QueueStatsResponse) GetAverageProcessingTime() float32`

GetAverageProcessingTime returns the AverageProcessingTime field if non-nil, zero value otherwise.

### GetAverageProcessingTimeOk

`func (o *QueueStatsResponse) GetAverageProcessingTimeOk() (*float32, bool)`

GetAverageProcessingTimeOk returns a tuple with the AverageProcessingTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAverageProcessingTime

`func (o *QueueStatsResponse) SetAverageProcessingTime(v float32)`

SetAverageProcessingTime sets AverageProcessingTime field to given value.

### HasAverageProcessingTime

`func (o *QueueStatsResponse) HasAverageProcessingTime() bool`

HasAverageProcessingTime returns a boolean if a field has been set.

### GetThroughput

`func (o *QueueStatsResponse) GetThroughput() float32`

GetThroughput returns the Throughput field if non-nil, zero value otherwise.

### GetThroughputOk

`func (o *QueueStatsResponse) GetThroughputOk() (*float32, bool)`

GetThroughputOk returns a tuple with the Throughput field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThroughput

`func (o *QueueStatsResponse) SetThroughput(v float32)`

SetThroughput sets Throughput field to given value.

### HasThroughput

`func (o *QueueStatsResponse) HasThroughput() bool`

HasThroughput returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


