# NotificationListResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]NotificationResponse**](NotificationResponse.md) |  | [optional] 
**Pagination** | Pointer to [**PaginationResponse**](PaginationResponse.md) |  | [optional] 

## Methods

### NewNotificationListResponse

`func NewNotificationListResponse() *NotificationListResponse`

NewNotificationListResponse instantiates a new NotificationListResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNotificationListResponseWithDefaults

`func NewNotificationListResponseWithDefaults() *NotificationListResponse`

NewNotificationListResponseWithDefaults instantiates a new NotificationListResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *NotificationListResponse) GetData() []NotificationResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *NotificationListResponse) GetDataOk() (*[]NotificationResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *NotificationListResponse) SetData(v []NotificationResponse)`

SetData sets Data field to given value.

### HasData

`func (o *NotificationListResponse) HasData() bool`

HasData returns a boolean if a field has been set.

### GetPagination

`func (o *NotificationListResponse) GetPagination() PaginationResponse`

GetPagination returns the Pagination field if non-nil, zero value otherwise.

### GetPaginationOk

`func (o *NotificationListResponse) GetPaginationOk() (*PaginationResponse, bool)`

GetPaginationOk returns a tuple with the Pagination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPagination

`func (o *NotificationListResponse) SetPagination(v PaginationResponse)`

SetPagination sets Pagination field to given value.

### HasPagination

`func (o *NotificationListResponse) HasPagination() bool`

HasPagination returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


