# DestinationUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to **string** |  | [optional] 
**Targets** | Pointer to [**[]DestinationTarget**](DestinationTarget.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Methods

### NewDestinationUpdateRequest

`func NewDestinationUpdateRequest() *DestinationUpdateRequest`

NewDestinationUpdateRequest instantiates a new DestinationUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDestinationUpdateRequestWithDefaults

`func NewDestinationUpdateRequestWithDefaults() *DestinationUpdateRequest`

NewDestinationUpdateRequestWithDefaults instantiates a new DestinationUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *DestinationUpdateRequest) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DestinationUpdateRequest) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DestinationUpdateRequest) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *DestinationUpdateRequest) HasType() bool`

HasType returns a boolean if a field has been set.

### GetTargets

`func (o *DestinationUpdateRequest) GetTargets() []DestinationTarget`

GetTargets returns the Targets field if non-nil, zero value otherwise.

### GetTargetsOk

`func (o *DestinationUpdateRequest) GetTargetsOk() (*[]DestinationTarget, bool)`

GetTargetsOk returns a tuple with the Targets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargets

`func (o *DestinationUpdateRequest) SetTargets(v []DestinationTarget)`

SetTargets sets Targets field to given value.

### HasTargets

`func (o *DestinationUpdateRequest) HasTargets() bool`

HasTargets returns a boolean if a field has been set.

### GetStatus

`func (o *DestinationUpdateRequest) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DestinationUpdateRequest) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DestinationUpdateRequest) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *DestinationUpdateRequest) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


