# QueueClearResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Message** | Pointer to **string** |  | [optional] 
**ClearedItems** | Pointer to **int32** |  | [optional] 
**Timestamp** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewQueueClearResponse

`func NewQueueClearResponse() *QueueClearResponse`

NewQueueClearResponse instantiates a new QueueClearResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQueueClearResponseWithDefaults

`func NewQueueClearResponseWithDefaults() *QueueClearResponse`

NewQueueClearResponseWithDefaults instantiates a new QueueClearResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMessage

`func (o *QueueClearResponse) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *QueueClearResponse) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *QueueClearResponse) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *QueueClearResponse) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetClearedItems

`func (o *QueueClearResponse) GetClearedItems() int32`

GetClearedItems returns the ClearedItems field if non-nil, zero value otherwise.

### GetClearedItemsOk

`func (o *QueueClearResponse) GetClearedItemsOk() (*int32, bool)`

GetClearedItemsOk returns a tuple with the ClearedItems field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClearedItems

`func (o *QueueClearResponse) SetClearedItems(v int32)`

SetClearedItems sets ClearedItems field to given value.

### HasClearedItems

`func (o *QueueClearResponse) HasClearedItems() bool`

HasClearedItems returns a boolean if a field has been set.

### GetTimestamp

`func (o *QueueClearResponse) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *QueueClearResponse) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *QueueClearResponse) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *QueueClearResponse) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


