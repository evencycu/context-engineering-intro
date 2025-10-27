# CompanyListResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]CompanyResponse**](CompanyResponse.md) |  | [optional] 
**Pagination** | Pointer to [**PaginationResponse**](PaginationResponse.md) |  | [optional] 

## Methods

### NewCompanyListResponse

`func NewCompanyListResponse() *CompanyListResponse`

NewCompanyListResponse instantiates a new CompanyListResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCompanyListResponseWithDefaults

`func NewCompanyListResponseWithDefaults() *CompanyListResponse`

NewCompanyListResponseWithDefaults instantiates a new CompanyListResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *CompanyListResponse) GetData() []CompanyResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *CompanyListResponse) GetDataOk() (*[]CompanyResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *CompanyListResponse) SetData(v []CompanyResponse)`

SetData sets Data field to given value.

### HasData

`func (o *CompanyListResponse) HasData() bool`

HasData returns a boolean if a field has been set.

### GetPagination

`func (o *CompanyListResponse) GetPagination() PaginationResponse`

GetPagination returns the Pagination field if non-nil, zero value otherwise.

### GetPaginationOk

`func (o *CompanyListResponse) GetPaginationOk() (*PaginationResponse, bool)`

GetPaginationOk returns a tuple with the Pagination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPagination

`func (o *CompanyListResponse) SetPagination(v PaginationResponse)`

SetPagination sets Pagination field to given value.

### HasPagination

`func (o *CompanyListResponse) HasPagination() bool`

HasPagination returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


