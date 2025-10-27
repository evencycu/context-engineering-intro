# BotFrameworkActivity

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Timestamp** | Pointer to **time.Time** |  | [optional] 
**ChannelId** | Pointer to **string** |  | [optional] 
**From** | Pointer to [**ChannelAccount**](ChannelAccount.md) |  | [optional] 
**Conversation** | Pointer to [**ConversationAccount**](ConversationAccount.md) |  | [optional] 
**Recipient** | Pointer to [**ChannelAccount**](ChannelAccount.md) |  | [optional] 
**Text** | Pointer to **string** |  | [optional] 
**Attachments** | Pointer to [**[]Attachment**](Attachment.md) |  | [optional] 

## Methods

### NewBotFrameworkActivity

`func NewBotFrameworkActivity() *BotFrameworkActivity`

NewBotFrameworkActivity instantiates a new BotFrameworkActivity object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBotFrameworkActivityWithDefaults

`func NewBotFrameworkActivityWithDefaults() *BotFrameworkActivity`

NewBotFrameworkActivityWithDefaults instantiates a new BotFrameworkActivity object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *BotFrameworkActivity) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *BotFrameworkActivity) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *BotFrameworkActivity) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *BotFrameworkActivity) HasType() bool`

HasType returns a boolean if a field has been set.

### GetId

`func (o *BotFrameworkActivity) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *BotFrameworkActivity) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *BotFrameworkActivity) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *BotFrameworkActivity) HasId() bool`

HasId returns a boolean if a field has been set.

### GetTimestamp

`func (o *BotFrameworkActivity) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *BotFrameworkActivity) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *BotFrameworkActivity) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *BotFrameworkActivity) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetChannelId

`func (o *BotFrameworkActivity) GetChannelId() string`

GetChannelId returns the ChannelId field if non-nil, zero value otherwise.

### GetChannelIdOk

`func (o *BotFrameworkActivity) GetChannelIdOk() (*string, bool)`

GetChannelIdOk returns a tuple with the ChannelId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelId

`func (o *BotFrameworkActivity) SetChannelId(v string)`

SetChannelId sets ChannelId field to given value.

### HasChannelId

`func (o *BotFrameworkActivity) HasChannelId() bool`

HasChannelId returns a boolean if a field has been set.

### GetFrom

`func (o *BotFrameworkActivity) GetFrom() ChannelAccount`

GetFrom returns the From field if non-nil, zero value otherwise.

### GetFromOk

`func (o *BotFrameworkActivity) GetFromOk() (*ChannelAccount, bool)`

GetFromOk returns a tuple with the From field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrom

`func (o *BotFrameworkActivity) SetFrom(v ChannelAccount)`

SetFrom sets From field to given value.

### HasFrom

`func (o *BotFrameworkActivity) HasFrom() bool`

HasFrom returns a boolean if a field has been set.

### GetConversation

`func (o *BotFrameworkActivity) GetConversation() ConversationAccount`

GetConversation returns the Conversation field if non-nil, zero value otherwise.

### GetConversationOk

`func (o *BotFrameworkActivity) GetConversationOk() (*ConversationAccount, bool)`

GetConversationOk returns a tuple with the Conversation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConversation

`func (o *BotFrameworkActivity) SetConversation(v ConversationAccount)`

SetConversation sets Conversation field to given value.

### HasConversation

`func (o *BotFrameworkActivity) HasConversation() bool`

HasConversation returns a boolean if a field has been set.

### GetRecipient

`func (o *BotFrameworkActivity) GetRecipient() ChannelAccount`

GetRecipient returns the Recipient field if non-nil, zero value otherwise.

### GetRecipientOk

`func (o *BotFrameworkActivity) GetRecipientOk() (*ChannelAccount, bool)`

GetRecipientOk returns a tuple with the Recipient field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecipient

`func (o *BotFrameworkActivity) SetRecipient(v ChannelAccount)`

SetRecipient sets Recipient field to given value.

### HasRecipient

`func (o *BotFrameworkActivity) HasRecipient() bool`

HasRecipient returns a boolean if a field has been set.

### GetText

`func (o *BotFrameworkActivity) GetText() string`

GetText returns the Text field if non-nil, zero value otherwise.

### GetTextOk

`func (o *BotFrameworkActivity) GetTextOk() (*string, bool)`

GetTextOk returns a tuple with the Text field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetText

`func (o *BotFrameworkActivity) SetText(v string)`

SetText sets Text field to given value.

### HasText

`func (o *BotFrameworkActivity) HasText() bool`

HasText returns a boolean if a field has been set.

### GetAttachments

`func (o *BotFrameworkActivity) GetAttachments() []Attachment`

GetAttachments returns the Attachments field if non-nil, zero value otherwise.

### GetAttachmentsOk

`func (o *BotFrameworkActivity) GetAttachmentsOk() (*[]Attachment, bool)`

GetAttachmentsOk returns a tuple with the Attachments field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttachments

`func (o *BotFrameworkActivity) SetAttachments(v []Attachment)`

SetAttachments sets Attachments field to given value.

### HasAttachments

`func (o *BotFrameworkActivity) HasAttachments() bool`

HasAttachments returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


