# NotificationSendRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProjectId** | **string** | Project ID | 
**SenderId** | Pointer to **string** | Sender user ID | [optional] 
**MessageType** | **string** | Type of message | 
**Content** | **string** | Message content | 
**Mentions** | Pointer to **[]string** | List of user mentions | [optional] 
**Priority** | Pointer to **string** | Message priority | [optional] [default to "normal"]
**Targets** | **[]string** | List of target destination IDs | 
**Metadata** | Pointer to **map[string]interface{}** | Additional metadata | [optional] 

## Methods

### NewNotificationSendRequest

`func NewNotificationSendRequest(projectId string, messageType string, content string, targets []string, ) *NotificationSendRequest`

NewNotificationSendRequest instantiates a new NotificationSendRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNotificationSendRequestWithDefaults

`func NewNotificationSendRequestWithDefaults() *NotificationSendRequest`

NewNotificationSendRequestWithDefaults instantiates a new NotificationSendRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProjectId

`func (o *NotificationSendRequest) GetProjectId() string`

GetProjectId returns the ProjectId field if non-nil, zero value otherwise.

### GetProjectIdOk

`func (o *NotificationSendRequest) GetProjectIdOk() (*string, bool)`

GetProjectIdOk returns a tuple with the ProjectId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectId

`func (o *NotificationSendRequest) SetProjectId(v string)`

SetProjectId sets ProjectId field to given value.


### GetSenderId

`func (o *NotificationSendRequest) GetSenderId() string`

GetSenderId returns the SenderId field if non-nil, zero value otherwise.

### GetSenderIdOk

`func (o *NotificationSendRequest) GetSenderIdOk() (*string, bool)`

GetSenderIdOk returns a tuple with the SenderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSenderId

`func (o *NotificationSendRequest) SetSenderId(v string)`

SetSenderId sets SenderId field to given value.

### HasSenderId

`func (o *NotificationSendRequest) HasSenderId() bool`

HasSenderId returns a boolean if a field has been set.

### GetMessageType

`func (o *NotificationSendRequest) GetMessageType() string`

GetMessageType returns the MessageType field if non-nil, zero value otherwise.

### GetMessageTypeOk

`func (o *NotificationSendRequest) GetMessageTypeOk() (*string, bool)`

GetMessageTypeOk returns a tuple with the MessageType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessageType

`func (o *NotificationSendRequest) SetMessageType(v string)`

SetMessageType sets MessageType field to given value.


### GetContent

`func (o *NotificationSendRequest) GetContent() string`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *NotificationSendRequest) GetContentOk() (*string, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *NotificationSendRequest) SetContent(v string)`

SetContent sets Content field to given value.


### GetMentions

`func (o *NotificationSendRequest) GetMentions() []string`

GetMentions returns the Mentions field if non-nil, zero value otherwise.

### GetMentionsOk

`func (o *NotificationSendRequest) GetMentionsOk() (*[]string, bool)`

GetMentionsOk returns a tuple with the Mentions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMentions

`func (o *NotificationSendRequest) SetMentions(v []string)`

SetMentions sets Mentions field to given value.

### HasMentions

`func (o *NotificationSendRequest) HasMentions() bool`

HasMentions returns a boolean if a field has been set.

### GetPriority

`func (o *NotificationSendRequest) GetPriority() string`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *NotificationSendRequest) GetPriorityOk() (*string, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *NotificationSendRequest) SetPriority(v string)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *NotificationSendRequest) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetTargets

`func (o *NotificationSendRequest) GetTargets() []string`

GetTargets returns the Targets field if non-nil, zero value otherwise.

### GetTargetsOk

`func (o *NotificationSendRequest) GetTargetsOk() (*[]string, bool)`

GetTargetsOk returns a tuple with the Targets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargets

`func (o *NotificationSendRequest) SetTargets(v []string)`

SetTargets sets Targets field to given value.


### GetMetadata

`func (o *NotificationSendRequest) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *NotificationSendRequest) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *NotificationSendRequest) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *NotificationSendRequest) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


