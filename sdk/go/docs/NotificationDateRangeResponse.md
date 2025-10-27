# NotificationDateRangeResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]NotificationResponse**](NotificationResponse.md) |  | [optional] 
**Count** | Pointer to **int32** |  | [optional] 
**DateRange** | Pointer to [**NotificationDateRangeResponseDateRange**](NotificationDateRangeResponseDateRange.md) |  | [optional] 

## Methods

### NewNotificationDateRangeResponse

`func NewNotificationDateRangeResponse() *NotificationDateRangeResponse`

NewNotificationDateRangeResponse instantiates a new NotificationDateRangeResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNotificationDateRangeResponseWithDefaults

`func NewNotificationDateRangeResponseWithDefaults() *NotificationDateRangeResponse`

NewNotificationDateRangeResponseWithDefaults instantiates a new NotificationDateRangeResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *NotificationDateRangeResponse) GetData() []NotificationResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *NotificationDateRangeResponse) GetDataOk() (*[]NotificationResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *NotificationDateRangeResponse) SetData(v []NotificationResponse)`

SetData sets Data field to given value.

### HasData

`func (o *NotificationDateRangeResponse) HasData() bool`

HasData returns a boolean if a field has been set.

### GetCount

`func (o *NotificationDateRangeResponse) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *NotificationDateRangeResponse) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *NotificationDateRangeResponse) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *NotificationDateRangeResponse) HasCount() bool`

HasCount returns a boolean if a field has been set.

### GetDateRange

`func (o *NotificationDateRangeResponse) GetDateRange() NotificationDateRangeResponseDateRange`

GetDateRange returns the DateRange field if non-nil, zero value otherwise.

### GetDateRangeOk

`func (o *NotificationDateRangeResponse) GetDateRangeOk() (*NotificationDateRangeResponseDateRange, bool)`

GetDateRangeOk returns a tuple with the DateRange field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDateRange

`func (o *NotificationDateRangeResponse) SetDateRange(v NotificationDateRangeResponseDateRange)`

SetDateRange sets DateRange field to given value.

### HasDateRange

`func (o *NotificationDateRangeResponse) HasDateRange() bool`

HasDateRange returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


