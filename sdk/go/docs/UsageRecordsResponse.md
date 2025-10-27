# UsageRecordsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]UsageRecordResponse**](UsageRecordResponse.md) |  | [optional] 
**Pagination** | Pointer to [**PaginationResponse**](PaginationResponse.md) |  | [optional] 

## Methods

### NewUsageRecordsResponse

`func NewUsageRecordsResponse() *UsageRecordsResponse`

NewUsageRecordsResponse instantiates a new UsageRecordsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUsageRecordsResponseWithDefaults

`func NewUsageRecordsResponseWithDefaults() *UsageRecordsResponse`

NewUsageRecordsResponseWithDefaults instantiates a new UsageRecordsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *UsageRecordsResponse) GetData() []UsageRecordResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *UsageRecordsResponse) GetDataOk() (*[]UsageRecordResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *UsageRecordsResponse) SetData(v []UsageRecordResponse)`

SetData sets Data field to given value.

### HasData

`func (o *UsageRecordsResponse) HasData() bool`

HasData returns a boolean if a field has been set.

### GetPagination

`func (o *UsageRecordsResponse) GetPagination() PaginationResponse`

GetPagination returns the Pagination field if non-nil, zero value otherwise.

### GetPaginationOk

`func (o *UsageRecordsResponse) GetPaginationOk() (*PaginationResponse, bool)`

GetPaginationOk returns a tuple with the Pagination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPagination

`func (o *UsageRecordsResponse) SetPagination(v PaginationResponse)`

SetPagination sets Pagination field to given value.

### HasPagination

`func (o *UsageRecordsResponse) HasPagination() bool`

HasPagination returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


