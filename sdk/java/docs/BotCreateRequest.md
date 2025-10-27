

# BotCreateRequest


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**name** | **String** | Bot name |  |
|**appId** | **String** | Microsoft App ID |  |
|**appPassword** | **String** | Microsoft App Password |  |
|**capabilities** | **List&lt;String&gt;** | Bot capabilities |  |
|**status** | [**StatusEnum**](#StatusEnum) |  |  [optional] |
|**description** | **String** | Bot description |  [optional] |



## Enum: StatusEnum

| Name | Value |
|---- | -----|
| ACTIVE | &quot;active&quot; |
| INACTIVE | &quot;inactive&quot; |
| ERROR | &quot;error&quot; |



