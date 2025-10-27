# TokenCacheMetrics

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Hits** | Pointer to **int32** |  | [optional] 
**Misses** | Pointer to **int32** |  | [optional] 
**HitRate** | Pointer to **float32** |  | [optional] 
**TotalKeys** | Pointer to **int32** |  | [optional] 

## Methods

### NewTokenCacheMetrics

`func NewTokenCacheMetrics() *TokenCacheMetrics`

NewTokenCacheMetrics instantiates a new TokenCacheMetrics object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTokenCacheMetricsWithDefaults

`func NewTokenCacheMetricsWithDefaults() *TokenCacheMetrics`

NewTokenCacheMetricsWithDefaults instantiates a new TokenCacheMetrics object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHits

`func (o *TokenCacheMetrics) GetHits() int32`

GetHits returns the Hits field if non-nil, zero value otherwise.

### GetHitsOk

`func (o *TokenCacheMetrics) GetHitsOk() (*int32, bool)`

GetHitsOk returns a tuple with the Hits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHits

`func (o *TokenCacheMetrics) SetHits(v int32)`

SetHits sets Hits field to given value.

### HasHits

`func (o *TokenCacheMetrics) HasHits() bool`

HasHits returns a boolean if a field has been set.

### GetMisses

`func (o *TokenCacheMetrics) GetMisses() int32`

GetMisses returns the Misses field if non-nil, zero value otherwise.

### GetMissesOk

`func (o *TokenCacheMetrics) GetMissesOk() (*int32, bool)`

GetMissesOk returns a tuple with the Misses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMisses

`func (o *TokenCacheMetrics) SetMisses(v int32)`

SetMisses sets Misses field to given value.

### HasMisses

`func (o *TokenCacheMetrics) HasMisses() bool`

HasMisses returns a boolean if a field has been set.

### GetHitRate

`func (o *TokenCacheMetrics) GetHitRate() float32`

GetHitRate returns the HitRate field if non-nil, zero value otherwise.

### GetHitRateOk

`func (o *TokenCacheMetrics) GetHitRateOk() (*float32, bool)`

GetHitRateOk returns a tuple with the HitRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHitRate

`func (o *TokenCacheMetrics) SetHitRate(v float32)`

SetHitRate sets HitRate field to given value.

### HasHitRate

`func (o *TokenCacheMetrics) HasHitRate() bool`

HasHitRate returns a boolean if a field has been set.

### GetTotalKeys

`func (o *TokenCacheMetrics) GetTotalKeys() int32`

GetTotalKeys returns the TotalKeys field if non-nil, zero value otherwise.

### GetTotalKeysOk

`func (o *TokenCacheMetrics) GetTotalKeysOk() (*int32, bool)`

GetTotalKeysOk returns a tuple with the TotalKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalKeys

`func (o *TokenCacheMetrics) SetTotalKeys(v int32)`

SetTotalKeys sets TotalKeys field to given value.

### HasTotalKeys

`func (o *TokenCacheMetrics) HasTotalKeys() bool`

HasTotalKeys returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


