# ProjectDestinationsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NotifyKey** | Pointer to **string** |  | [optional] 
**Destinations** | Pointer to [**[]DestinationResponse**](DestinationResponse.md) |  | [optional] 

## Methods

### NewProjectDestinationsResponse

`func NewProjectDestinationsResponse() *ProjectDestinationsResponse`

NewProjectDestinationsResponse instantiates a new ProjectDestinationsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectDestinationsResponseWithDefaults

`func NewProjectDestinationsResponseWithDefaults() *ProjectDestinationsResponse`

NewProjectDestinationsResponseWithDefaults instantiates a new ProjectDestinationsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNotifyKey

`func (o *ProjectDestinationsResponse) GetNotifyKey() string`

GetNotifyKey returns the NotifyKey field if non-nil, zero value otherwise.

### GetNotifyKeyOk

`func (o *ProjectDestinationsResponse) GetNotifyKeyOk() (*string, bool)`

GetNotifyKeyOk returns a tuple with the NotifyKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifyKey

`func (o *ProjectDestinationsResponse) SetNotifyKey(v string)`

SetNotifyKey sets NotifyKey field to given value.

### HasNotifyKey

`func (o *ProjectDestinationsResponse) HasNotifyKey() bool`

HasNotifyKey returns a boolean if a field has been set.

### GetDestinations

`func (o *ProjectDestinationsResponse) GetDestinations() []DestinationResponse`

GetDestinations returns the Destinations field if non-nil, zero value otherwise.

### GetDestinationsOk

`func (o *ProjectDestinationsResponse) GetDestinationsOk() (*[]DestinationResponse, bool)`

GetDestinationsOk returns a tuple with the Destinations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinations

`func (o *ProjectDestinationsResponse) SetDestinations(v []DestinationResponse)`

SetDestinations sets Destinations field to given value.

### HasDestinations

`func (o *ProjectDestinationsResponse) HasDestinations() bool`

HasDestinations returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


