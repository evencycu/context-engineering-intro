# FileListResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]FileResponse**](FileResponse.md) |  | [optional] 
**Pagination** | Pointer to [**PaginationResponse**](PaginationResponse.md) |  | [optional] 

## Methods

### NewFileListResponse

`func NewFileListResponse() *FileListResponse`

NewFileListResponse instantiates a new FileListResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFileListResponseWithDefaults

`func NewFileListResponseWithDefaults() *FileListResponse`

NewFileListResponseWithDefaults instantiates a new FileListResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *FileListResponse) GetData() []FileResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *FileListResponse) GetDataOk() (*[]FileResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *FileListResponse) SetData(v []FileResponse)`

SetData sets Data field to given value.

### HasData

`func (o *FileListResponse) HasData() bool`

HasData returns a boolean if a field has been set.

### GetPagination

`func (o *FileListResponse) GetPagination() PaginationResponse`

GetPagination returns the Pagination field if non-nil, zero value otherwise.

### GetPaginationOk

`func (o *FileListResponse) GetPaginationOk() (*PaginationResponse, bool)`

GetPaginationOk returns a tuple with the Pagination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPagination

`func (o *FileListResponse) SetPagination(v PaginationResponse)`

SetPagination sets Pagination field to given value.

### HasPagination

`func (o *FileListResponse) HasPagination() bool`

HasPagination returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


