# ExternalNotifyRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NotifyKey** | **string** | Project notify key | 
**Message** | **string** | Notification message | 
**MessageType** | Pointer to **string** | Type of message | [optional] [default to "text"]
**Priority** | Pointer to **string** | Message priority | [optional] [default to "normal"]
**Targets** | Pointer to **[]string** | List of target emails, conversation IDs, or &#39;all&#39; | [optional] 
**Mentions** | Pointer to **[]string** | List of user mentions | [optional] 
**Metadata** | Pointer to **map[string]interface{}** | Additional metadata | [optional] 

## Methods

### NewExternalNotifyRequest

`func NewExternalNotifyRequest(notifyKey string, message string, ) *ExternalNotifyRequest`

NewExternalNotifyRequest instantiates a new ExternalNotifyRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExternalNotifyRequestWithDefaults

`func NewExternalNotifyRequestWithDefaults() *ExternalNotifyRequest`

NewExternalNotifyRequestWithDefaults instantiates a new ExternalNotifyRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNotifyKey

`func (o *ExternalNotifyRequest) GetNotifyKey() string`

GetNotifyKey returns the NotifyKey field if non-nil, zero value otherwise.

### GetNotifyKeyOk

`func (o *ExternalNotifyRequest) GetNotifyKeyOk() (*string, bool)`

GetNotifyKeyOk returns a tuple with the NotifyKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifyKey

`func (o *ExternalNotifyRequest) SetNotifyKey(v string)`

SetNotifyKey sets NotifyKey field to given value.


### GetMessage

`func (o *ExternalNotifyRequest) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ExternalNotifyRequest) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ExternalNotifyRequest) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetMessageType

`func (o *ExternalNotifyRequest) GetMessageType() string`

GetMessageType returns the MessageType field if non-nil, zero value otherwise.

### GetMessageTypeOk

`func (o *ExternalNotifyRequest) GetMessageTypeOk() (*string, bool)`

GetMessageTypeOk returns a tuple with the MessageType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessageType

`func (o *ExternalNotifyRequest) SetMessageType(v string)`

SetMessageType sets MessageType field to given value.

### HasMessageType

`func (o *ExternalNotifyRequest) HasMessageType() bool`

HasMessageType returns a boolean if a field has been set.

### GetPriority

`func (o *ExternalNotifyRequest) GetPriority() string`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *ExternalNotifyRequest) GetPriorityOk() (*string, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *ExternalNotifyRequest) SetPriority(v string)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *ExternalNotifyRequest) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetTargets

`func (o *ExternalNotifyRequest) GetTargets() []string`

GetTargets returns the Targets field if non-nil, zero value otherwise.

### GetTargetsOk

`func (o *ExternalNotifyRequest) GetTargetsOk() (*[]string, bool)`

GetTargetsOk returns a tuple with the Targets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargets

`func (o *ExternalNotifyRequest) SetTargets(v []string)`

SetTargets sets Targets field to given value.

### HasTargets

`func (o *ExternalNotifyRequest) HasTargets() bool`

HasTargets returns a boolean if a field has been set.

### GetMentions

`func (o *ExternalNotifyRequest) GetMentions() []string`

GetMentions returns the Mentions field if non-nil, zero value otherwise.

### GetMentionsOk

`func (o *ExternalNotifyRequest) GetMentionsOk() (*[]string, bool)`

GetMentionsOk returns a tuple with the Mentions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMentions

`func (o *ExternalNotifyRequest) SetMentions(v []string)`

SetMentions sets Mentions field to given value.

### HasMentions

`func (o *ExternalNotifyRequest) HasMentions() bool`

HasMentions returns a boolean if a field has been set.

### GetMetadata

`func (o *ExternalNotifyRequest) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *ExternalNotifyRequest) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *ExternalNotifyRequest) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *ExternalNotifyRequest) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


