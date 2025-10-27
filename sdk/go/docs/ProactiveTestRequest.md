# ProactiveTestRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConversationId** | **string** | Teams conversation ID | 
**Message** | **string** | Test message content | 

## Methods

### NewProactiveTestRequest

`func NewProactiveTestRequest(conversationId string, message string, ) *ProactiveTestRequest`

NewProactiveTestRequest instantiates a new ProactiveTestRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProactiveTestRequestWithDefaults

`func NewProactiveTestRequestWithDefaults() *ProactiveTestRequest`

NewProactiveTestRequestWithDefaults instantiates a new ProactiveTestRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConversationId

`func (o *ProactiveTestRequest) GetConversationId() string`

GetConversationId returns the ConversationId field if non-nil, zero value otherwise.

### GetConversationIdOk

`func (o *ProactiveTestRequest) GetConversationIdOk() (*string, bool)`

GetConversationIdOk returns a tuple with the ConversationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConversationId

`func (o *ProactiveTestRequest) SetConversationId(v string)`

SetConversationId sets ConversationId field to given value.


### GetMessage

`func (o *ProactiveTestRequest) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ProactiveTestRequest) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ProactiveTestRequest) SetMessage(v string)`

SetMessage sets Message field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


