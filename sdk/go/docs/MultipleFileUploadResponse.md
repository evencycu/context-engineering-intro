# MultipleFileUploadResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Files** | Pointer to [**[]FileResponse**](FileResponse.md) |  | [optional] 
**TotalUploaded** | Pointer to **int32** |  | [optional] 
**FailedUploads** | Pointer to **[]string** |  | [optional] 

## Methods

### NewMultipleFileUploadResponse

`func NewMultipleFileUploadResponse() *MultipleFileUploadResponse`

NewMultipleFileUploadResponse instantiates a new MultipleFileUploadResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMultipleFileUploadResponseWithDefaults

`func NewMultipleFileUploadResponseWithDefaults() *MultipleFileUploadResponse`

NewMultipleFileUploadResponseWithDefaults instantiates a new MultipleFileUploadResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFiles

`func (o *MultipleFileUploadResponse) GetFiles() []FileResponse`

GetFiles returns the Files field if non-nil, zero value otherwise.

### GetFilesOk

`func (o *MultipleFileUploadResponse) GetFilesOk() (*[]FileResponse, bool)`

GetFilesOk returns a tuple with the Files field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFiles

`func (o *MultipleFileUploadResponse) SetFiles(v []FileResponse)`

SetFiles sets Files field to given value.

### HasFiles

`func (o *MultipleFileUploadResponse) HasFiles() bool`

HasFiles returns a boolean if a field has been set.

### GetTotalUploaded

`func (o *MultipleFileUploadResponse) GetTotalUploaded() int32`

GetTotalUploaded returns the TotalUploaded field if non-nil, zero value otherwise.

### GetTotalUploadedOk

`func (o *MultipleFileUploadResponse) GetTotalUploadedOk() (*int32, bool)`

GetTotalUploadedOk returns a tuple with the TotalUploaded field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalUploaded

`func (o *MultipleFileUploadResponse) SetTotalUploaded(v int32)`

SetTotalUploaded sets TotalUploaded field to given value.

### HasTotalUploaded

`func (o *MultipleFileUploadResponse) HasTotalUploaded() bool`

HasTotalUploaded returns a boolean if a field has been set.

### GetFailedUploads

`func (o *MultipleFileUploadResponse) GetFailedUploads() []string`

GetFailedUploads returns the FailedUploads field if non-nil, zero value otherwise.

### GetFailedUploadsOk

`func (o *MultipleFileUploadResponse) GetFailedUploadsOk() (*[]string, bool)`

GetFailedUploadsOk returns a tuple with the FailedUploads field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailedUploads

`func (o *MultipleFileUploadResponse) SetFailedUploads(v []string)`

SetFailedUploads sets FailedUploads field to given value.

### HasFailedUploads

`func (o *MultipleFileUploadResponse) HasFailedUploads() bool`

HasFailedUploads returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


