# NotificationResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**ProjectId** | Pointer to **string** |  | [optional] 
**SenderId** | Pointer to **string** |  | [optional] 
**MessageType** | Pointer to **string** |  | [optional] 
**Content** | Pointer to **string** |  | [optional] 
**Mentions** | Pointer to **[]string** |  | [optional] 
**Priority** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**ErrorMessage** | Pointer to **string** |  | [optional] 
**DestinationsSent** | Pointer to **int32** |  | [optional] 
**DestinationsFailed** | Pointer to **int32** |  | [optional] 
**TotalDestinations** | Pointer to **int32** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 
**SentAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewNotificationResponse

`func NewNotificationResponse() *NotificationResponse`

NewNotificationResponse instantiates a new NotificationResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNotificationResponseWithDefaults

`func NewNotificationResponseWithDefaults() *NotificationResponse`

NewNotificationResponseWithDefaults instantiates a new NotificationResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *NotificationResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *NotificationResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *NotificationResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *NotificationResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetProjectId

`func (o *NotificationResponse) GetProjectId() string`

GetProjectId returns the ProjectId field if non-nil, zero value otherwise.

### GetProjectIdOk

`func (o *NotificationResponse) GetProjectIdOk() (*string, bool)`

GetProjectIdOk returns a tuple with the ProjectId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectId

`func (o *NotificationResponse) SetProjectId(v string)`

SetProjectId sets ProjectId field to given value.

### HasProjectId

`func (o *NotificationResponse) HasProjectId() bool`

HasProjectId returns a boolean if a field has been set.

### GetSenderId

`func (o *NotificationResponse) GetSenderId() string`

GetSenderId returns the SenderId field if non-nil, zero value otherwise.

### GetSenderIdOk

`func (o *NotificationResponse) GetSenderIdOk() (*string, bool)`

GetSenderIdOk returns a tuple with the SenderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSenderId

`func (o *NotificationResponse) SetSenderId(v string)`

SetSenderId sets SenderId field to given value.

### HasSenderId

`func (o *NotificationResponse) HasSenderId() bool`

HasSenderId returns a boolean if a field has been set.

### GetMessageType

`func (o *NotificationResponse) GetMessageType() string`

GetMessageType returns the MessageType field if non-nil, zero value otherwise.

### GetMessageTypeOk

`func (o *NotificationResponse) GetMessageTypeOk() (*string, bool)`

GetMessageTypeOk returns a tuple with the MessageType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessageType

`func (o *NotificationResponse) SetMessageType(v string)`

SetMessageType sets MessageType field to given value.

### HasMessageType

`func (o *NotificationResponse) HasMessageType() bool`

HasMessageType returns a boolean if a field has been set.

### GetContent

`func (o *NotificationResponse) GetContent() string`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *NotificationResponse) GetContentOk() (*string, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *NotificationResponse) SetContent(v string)`

SetContent sets Content field to given value.

### HasContent

`func (o *NotificationResponse) HasContent() bool`

HasContent returns a boolean if a field has been set.

### GetMentions

`func (o *NotificationResponse) GetMentions() []string`

GetMentions returns the Mentions field if non-nil, zero value otherwise.

### GetMentionsOk

`func (o *NotificationResponse) GetMentionsOk() (*[]string, bool)`

GetMentionsOk returns a tuple with the Mentions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMentions

`func (o *NotificationResponse) SetMentions(v []string)`

SetMentions sets Mentions field to given value.

### HasMentions

`func (o *NotificationResponse) HasMentions() bool`

HasMentions returns a boolean if a field has been set.

### GetPriority

`func (o *NotificationResponse) GetPriority() string`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *NotificationResponse) GetPriorityOk() (*string, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *NotificationResponse) SetPriority(v string)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *NotificationResponse) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetStatus

`func (o *NotificationResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *NotificationResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *NotificationResponse) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *NotificationResponse) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetErrorMessage

`func (o *NotificationResponse) GetErrorMessage() string`

GetErrorMessage returns the ErrorMessage field if non-nil, zero value otherwise.

### GetErrorMessageOk

`func (o *NotificationResponse) GetErrorMessageOk() (*string, bool)`

GetErrorMessageOk returns a tuple with the ErrorMessage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorMessage

`func (o *NotificationResponse) SetErrorMessage(v string)`

SetErrorMessage sets ErrorMessage field to given value.

### HasErrorMessage

`func (o *NotificationResponse) HasErrorMessage() bool`

HasErrorMessage returns a boolean if a field has been set.

### GetDestinationsSent

`func (o *NotificationResponse) GetDestinationsSent() int32`

GetDestinationsSent returns the DestinationsSent field if non-nil, zero value otherwise.

### GetDestinationsSentOk

`func (o *NotificationResponse) GetDestinationsSentOk() (*int32, bool)`

GetDestinationsSentOk returns a tuple with the DestinationsSent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationsSent

`func (o *NotificationResponse) SetDestinationsSent(v int32)`

SetDestinationsSent sets DestinationsSent field to given value.

### HasDestinationsSent

`func (o *NotificationResponse) HasDestinationsSent() bool`

HasDestinationsSent returns a boolean if a field has been set.

### GetDestinationsFailed

`func (o *NotificationResponse) GetDestinationsFailed() int32`

GetDestinationsFailed returns the DestinationsFailed field if non-nil, zero value otherwise.

### GetDestinationsFailedOk

`func (o *NotificationResponse) GetDestinationsFailedOk() (*int32, bool)`

GetDestinationsFailedOk returns a tuple with the DestinationsFailed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationsFailed

`func (o *NotificationResponse) SetDestinationsFailed(v int32)`

SetDestinationsFailed sets DestinationsFailed field to given value.

### HasDestinationsFailed

`func (o *NotificationResponse) HasDestinationsFailed() bool`

HasDestinationsFailed returns a boolean if a field has been set.

### GetTotalDestinations

`func (o *NotificationResponse) GetTotalDestinations() int32`

GetTotalDestinations returns the TotalDestinations field if non-nil, zero value otherwise.

### GetTotalDestinationsOk

`func (o *NotificationResponse) GetTotalDestinationsOk() (*int32, bool)`

GetTotalDestinationsOk returns a tuple with the TotalDestinations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalDestinations

`func (o *NotificationResponse) SetTotalDestinations(v int32)`

SetTotalDestinations sets TotalDestinations field to given value.

### HasTotalDestinations

`func (o *NotificationResponse) HasTotalDestinations() bool`

HasTotalDestinations returns a boolean if a field has been set.

### GetCreatedAt

`func (o *NotificationResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *NotificationResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *NotificationResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *NotificationResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *NotificationResponse) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *NotificationResponse) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *NotificationResponse) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *NotificationResponse) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetSentAt

`func (o *NotificationResponse) GetSentAt() time.Time`

GetSentAt returns the SentAt field if non-nil, zero value otherwise.

### GetSentAtOk

`func (o *NotificationResponse) GetSentAtOk() (*time.Time, bool)`

GetSentAtOk returns a tuple with the SentAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSentAt

`func (o *NotificationResponse) SetSentAt(v time.Time)`

SetSentAt sets SentAt field to given value.

### HasSentAt

`func (o *NotificationResponse) HasSentAt() bool`

HasSentAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


