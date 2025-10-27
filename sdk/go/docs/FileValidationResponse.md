# FileValidationResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Valid** | Pointer to **bool** |  | [optional] 
**Errors** | Pointer to **[]string** |  | [optional] 
**Warnings** | Pointer to **[]string** |  | [optional] 
**FileInfo** | Pointer to [**FileValidationResponseFileInfo**](FileValidationResponseFileInfo.md) |  | [optional] 

## Methods

### NewFileValidationResponse

`func NewFileValidationResponse() *FileValidationResponse`

NewFileValidationResponse instantiates a new FileValidationResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFileValidationResponseWithDefaults

`func NewFileValidationResponseWithDefaults() *FileValidationResponse`

NewFileValidationResponseWithDefaults instantiates a new FileValidationResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetValid

`func (o *FileValidationResponse) GetValid() bool`

GetValid returns the Valid field if non-nil, zero value otherwise.

### GetValidOk

`func (o *FileValidationResponse) GetValidOk() (*bool, bool)`

GetValidOk returns a tuple with the Valid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValid

`func (o *FileValidationResponse) SetValid(v bool)`

SetValid sets Valid field to given value.

### HasValid

`func (o *FileValidationResponse) HasValid() bool`

HasValid returns a boolean if a field has been set.

### GetErrors

`func (o *FileValidationResponse) GetErrors() []string`

GetErrors returns the Errors field if non-nil, zero value otherwise.

### GetErrorsOk

`func (o *FileValidationResponse) GetErrorsOk() (*[]string, bool)`

GetErrorsOk returns a tuple with the Errors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrors

`func (o *FileValidationResponse) SetErrors(v []string)`

SetErrors sets Errors field to given value.

### HasErrors

`func (o *FileValidationResponse) HasErrors() bool`

HasErrors returns a boolean if a field has been set.

### GetWarnings

`func (o *FileValidationResponse) GetWarnings() []string`

GetWarnings returns the Warnings field if non-nil, zero value otherwise.

### GetWarningsOk

`func (o *FileValidationResponse) GetWarningsOk() (*[]string, bool)`

GetWarningsOk returns a tuple with the Warnings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarnings

`func (o *FileValidationResponse) SetWarnings(v []string)`

SetWarnings sets Warnings field to given value.

### HasWarnings

`func (o *FileValidationResponse) HasWarnings() bool`

HasWarnings returns a boolean if a field has been set.

### GetFileInfo

`func (o *FileValidationResponse) GetFileInfo() FileValidationResponseFileInfo`

GetFileInfo returns the FileInfo field if non-nil, zero value otherwise.

### GetFileInfoOk

`func (o *FileValidationResponse) GetFileInfoOk() (*FileValidationResponseFileInfo, bool)`

GetFileInfoOk returns a tuple with the FileInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileInfo

`func (o *FileValidationResponse) SetFileInfo(v FileValidationResponseFileInfo)`

SetFileInfo sets FileInfo field to given value.

### HasFileInfo

`func (o *FileValidationResponse) HasFileInfo() bool`

HasFileInfo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


