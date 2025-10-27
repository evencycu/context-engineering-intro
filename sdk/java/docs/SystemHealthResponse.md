

# SystemHealthResponse


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**status** | [**StatusEnum**](#StatusEnum) |  |  [optional] |
|**timestamp** | **OffsetDateTime** |  |  [optional] |
|**services** | [**Map&lt;String, SystemHealthResponseServicesValue&gt;**](SystemHealthResponseServicesValue.md) |  |  [optional] |
|**uptime** | **String** |  |  [optional] |
|**version** | **String** |  |  [optional] |



## Enum: StatusEnum

| Name | Value |
|---- | -----|
| HEALTHY | &quot;healthy&quot; |
| DEGRADED | &quot;degraded&quot; |
| UNHEALTHY | &quot;unhealthy&quot; |



