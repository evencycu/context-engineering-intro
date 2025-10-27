# BillingPlansResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]BillingPlanResponse**](BillingPlanResponse.md) |  | [optional] 
**Pagination** | Pointer to [**PaginationResponse**](PaginationResponse.md) |  | [optional] 

## Methods

### NewBillingPlansResponse

`func NewBillingPlansResponse() *BillingPlansResponse`

NewBillingPlansResponse instantiates a new BillingPlansResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBillingPlansResponseWithDefaults

`func NewBillingPlansResponseWithDefaults() *BillingPlansResponse`

NewBillingPlansResponseWithDefaults instantiates a new BillingPlansResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *BillingPlansResponse) GetData() []BillingPlanResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *BillingPlansResponse) GetDataOk() (*[]BillingPlanResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *BillingPlansResponse) SetData(v []BillingPlanResponse)`

SetData sets Data field to given value.

### HasData

`func (o *BillingPlansResponse) HasData() bool`

HasData returns a boolean if a field has been set.

### GetPagination

`func (o *BillingPlansResponse) GetPagination() PaginationResponse`

GetPagination returns the Pagination field if non-nil, zero value otherwise.

### GetPaginationOk

`func (o *BillingPlansResponse) GetPaginationOk() (*PaginationResponse, bool)`

GetPaginationOk returns a tuple with the Pagination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPagination

`func (o *BillingPlansResponse) SetPagination(v PaginationResponse)`

SetPagination sets Pagination field to given value.

### HasPagination

`func (o *BillingPlansResponse) HasPagination() bool`

HasPagination returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


