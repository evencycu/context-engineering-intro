# ExternalNotifyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**Data** | Pointer to [**NotifyResponse**](NotifyResponse.md) |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 

## Methods

### NewExternalNotifyResponse

`func NewExternalNotifyResponse() *ExternalNotifyResponse`

NewExternalNotifyResponse instantiates a new ExternalNotifyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExternalNotifyResponseWithDefaults

`func NewExternalNotifyResponseWithDefaults() *ExternalNotifyResponse`

NewExternalNotifyResponseWithDefaults instantiates a new ExternalNotifyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSuccess

`func (o *ExternalNotifyResponse) GetSuccess() bool`

GetSuccess returns the Success field if non-nil, zero value otherwise.

### GetSuccessOk

`func (o *ExternalNotifyResponse) GetSuccessOk() (*bool, bool)`

GetSuccessOk returns a tuple with the Success field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccess

`func (o *ExternalNotifyResponse) SetSuccess(v bool)`

SetSuccess sets Success field to given value.

### HasSuccess

`func (o *ExternalNotifyResponse) HasSuccess() bool`

HasSuccess returns a boolean if a field has been set.

### GetData

`func (o *ExternalNotifyResponse) GetData() NotifyResponse`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ExternalNotifyResponse) GetDataOk() (*NotifyResponse, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ExternalNotifyResponse) SetData(v NotifyResponse)`

SetData sets Data field to given value.

### HasData

`func (o *ExternalNotifyResponse) HasData() bool`

HasData returns a boolean if a field has been set.

### GetError

`func (o *ExternalNotifyResponse) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *ExternalNotifyResponse) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *ExternalNotifyResponse) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *ExternalNotifyResponse) HasError() bool`

HasError returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


