

# ProjectCreateRequest


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**companyId** | **UUID** | Company ID |  |
|**notifyKey** | **String** | Project notify key |  |
|**description** | **String** | Project description |  |
|**dailyLimit** | **Integer** | Daily notification limit |  [optional] |
|**monthlyLimit** | **Integer** | Monthly notification limit |  [optional] |
|**priority** | [**PriorityEnum**](#PriorityEnum) | Project priority |  [optional] |
|**createdBy** | **UUID** | Creator user ID |  |



## Enum: PriorityEnum

| Name | Value |
|---- | -----|
| LOW | &quot;low&quot; |
| NORMAL | &quot;normal&quot; |
| HIGH | &quot;high&quot; |



