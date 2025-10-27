# DestinationListResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]DestinationResponse**](DestinationResponse.md) |  | [optional] 
**Pagination** | Pointer to [**PaginationResponse**](PaginationResponse.md) |  | [optional] 

## Methods

### NewDestinationListResponse

`func NewDestinationListResponse() *DestinationListResponse`

NewDestinationListResponse instantiates a new DestinationListResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDestinationListResponseWithDefaults

`func NewDestinationListResponseWithDefaults() *DestinationListResponse`

NewDestinationListResponseWithDefaults instantiates a new DestinationListResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *DestinationListResponse) GetData() []DestinationResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DestinationListResponse) GetDataOk() (*[]DestinationResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DestinationListResponse) SetData(v []DestinationResponse)`

SetData sets Data field to given value.

### HasData

`func (o *DestinationListResponse) HasData() bool`

HasData returns a boolean if a field has been set.

### GetPagination

`func (o *DestinationListResponse) GetPagination() PaginationResponse`

GetPagination returns the Pagination field if non-nil, zero value otherwise.

### GetPaginationOk

`func (o *DestinationListResponse) GetPaginationOk() (*PaginationResponse, bool)`

GetPaginationOk returns a tuple with the Pagination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPagination

`func (o *DestinationListResponse) SetPagination(v PaginationResponse)`

SetPagination sets Pagination field to given value.

### HasPagination

`func (o *DestinationListResponse) HasPagination() bool`

HasPagination returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


