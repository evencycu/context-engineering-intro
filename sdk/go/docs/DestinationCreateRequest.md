# DestinationCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProjectId** | **string** | Project ID | 
**BotId** | **string** | Bot ID | 
**Type** | **string** | Destination type | 
**Targets** | [**[]DestinationTarget**](DestinationTarget.md) | List of targets | 
**Status** | Pointer to **string** | Destination status | [optional] [default to "active"]

## Methods

### NewDestinationCreateRequest

`func NewDestinationCreateRequest(projectId string, botId string, type_ string, targets []DestinationTarget, ) *DestinationCreateRequest`

NewDestinationCreateRequest instantiates a new DestinationCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDestinationCreateRequestWithDefaults

`func NewDestinationCreateRequestWithDefaults() *DestinationCreateRequest`

NewDestinationCreateRequestWithDefaults instantiates a new DestinationCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProjectId

`func (o *DestinationCreateRequest) GetProjectId() string`

GetProjectId returns the ProjectId field if non-nil, zero value otherwise.

### GetProjectIdOk

`func (o *DestinationCreateRequest) GetProjectIdOk() (*string, bool)`

GetProjectIdOk returns a tuple with the ProjectId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectId

`func (o *DestinationCreateRequest) SetProjectId(v string)`

SetProjectId sets ProjectId field to given value.


### GetBotId

`func (o *DestinationCreateRequest) GetBotId() string`

GetBotId returns the BotId field if non-nil, zero value otherwise.

### GetBotIdOk

`func (o *DestinationCreateRequest) GetBotIdOk() (*string, bool)`

GetBotIdOk returns a tuple with the BotId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBotId

`func (o *DestinationCreateRequest) SetBotId(v string)`

SetBotId sets BotId field to given value.


### GetType

`func (o *DestinationCreateRequest) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DestinationCreateRequest) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DestinationCreateRequest) SetType(v string)`

SetType sets Type field to given value.


### GetTargets

`func (o *DestinationCreateRequest) GetTargets() []DestinationTarget`

GetTargets returns the Targets field if non-nil, zero value otherwise.

### GetTargetsOk

`func (o *DestinationCreateRequest) GetTargetsOk() (*[]DestinationTarget, bool)`

GetTargetsOk returns a tuple with the Targets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargets

`func (o *DestinationCreateRequest) SetTargets(v []DestinationTarget)`

SetTargets sets Targets field to given value.


### GetStatus

`func (o *DestinationCreateRequest) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DestinationCreateRequest) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DestinationCreateRequest) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *DestinationCreateRequest) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


