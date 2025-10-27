# ProjectListResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]ProjectResponse**](ProjectResponse.md) |  | [optional] 
**Pagination** | Pointer to [**PaginationResponse**](PaginationResponse.md) |  | [optional] 

## Methods

### NewProjectListResponse

`func NewProjectListResponse() *ProjectListResponse`

NewProjectListResponse instantiates a new ProjectListResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectListResponseWithDefaults

`func NewProjectListResponseWithDefaults() *ProjectListResponse`

NewProjectListResponseWithDefaults instantiates a new ProjectListResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *ProjectListResponse) GetData() []ProjectResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ProjectListResponse) GetDataOk() (*[]ProjectResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ProjectListResponse) SetData(v []ProjectResponse)`

SetData sets Data field to given value.

### HasData

`func (o *ProjectListResponse) HasData() bool`

HasData returns a boolean if a field has been set.

### GetPagination

`func (o *ProjectListResponse) GetPagination() PaginationResponse`

GetPagination returns the Pagination field if non-nil, zero value otherwise.

### GetPaginationOk

`func (o *ProjectListResponse) GetPaginationOk() (*PaginationResponse, bool)`

GetPaginationOk returns a tuple with the Pagination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPagination

`func (o *ProjectListResponse) SetPagination(v PaginationResponse)`

SetPagination sets Pagination field to given value.

### HasPagination

`func (o *ProjectListResponse) HasPagination() bool`

HasPagination returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


