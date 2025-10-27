# BotCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Bot name | 
**AppId** | **string** | Microsoft App ID | 
**AppPassword** | **string** | Microsoft App Password | 
**Capabilities** | **[]string** | Bot capabilities | 
**Status** | Pointer to **string** |  | [optional] [default to "active"]
**Description** | Pointer to **string** | Bot description | [optional] 

## Methods

### NewBotCreateRequest

`func NewBotCreateRequest(name string, appId string, appPassword string, capabilities []string, ) *BotCreateRequest`

NewBotCreateRequest instantiates a new BotCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBotCreateRequestWithDefaults

`func NewBotCreateRequestWithDefaults() *BotCreateRequest`

NewBotCreateRequestWithDefaults instantiates a new BotCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BotCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BotCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BotCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetAppId

`func (o *BotCreateRequest) GetAppId() string`

GetAppId returns the AppId field if non-nil, zero value otherwise.

### GetAppIdOk

`func (o *BotCreateRequest) GetAppIdOk() (*string, bool)`

GetAppIdOk returns a tuple with the AppId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAppId

`func (o *BotCreateRequest) SetAppId(v string)`

SetAppId sets AppId field to given value.


### GetAppPassword

`func (o *BotCreateRequest) GetAppPassword() string`

GetAppPassword returns the AppPassword field if non-nil, zero value otherwise.

### GetAppPasswordOk

`func (o *BotCreateRequest) GetAppPasswordOk() (*string, bool)`

GetAppPasswordOk returns a tuple with the AppPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAppPassword

`func (o *BotCreateRequest) SetAppPassword(v string)`

SetAppPassword sets AppPassword field to given value.


### GetCapabilities

`func (o *BotCreateRequest) GetCapabilities() []string`

GetCapabilities returns the Capabilities field if non-nil, zero value otherwise.

### GetCapabilitiesOk

`func (o *BotCreateRequest) GetCapabilitiesOk() (*[]string, bool)`

GetCapabilitiesOk returns a tuple with the Capabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapabilities

`func (o *BotCreateRequest) SetCapabilities(v []string)`

SetCapabilities sets Capabilities field to given value.


### GetStatus

`func (o *BotCreateRequest) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *BotCreateRequest) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *BotCreateRequest) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *BotCreateRequest) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetDescription

`func (o *BotCreateRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *BotCreateRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *BotCreateRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *BotCreateRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


