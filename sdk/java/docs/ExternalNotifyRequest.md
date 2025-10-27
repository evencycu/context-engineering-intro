

# ExternalNotifyRequest


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**notifyKey** | **String** | Project notify key |  |
|**message** | **String** | Notification message |  |
|**messageType** | [**MessageTypeEnum**](#MessageTypeEnum) | Type of message |  [optional] |
|**priority** | [**PriorityEnum**](#PriorityEnum) | Message priority |  [optional] |
|**targets** | **List&lt;String&gt;** | List of target emails, conversation IDs, or &#39;all&#39; |  [optional] |
|**mentions** | **List&lt;String&gt;** | List of user mentions |  [optional] |
|**metadata** | **Map&lt;String, Object&gt;** | Additional metadata |  [optional] |



## Enum: MessageTypeEnum

| Name | Value |
|---- | -----|
| TEXT | &quot;text&quot; |
| FILE | &quot;file&quot; |
| ADAPTIVE_CARD | &quot;adaptive_card&quot; |



## Enum: PriorityEnum

| Name | Value |
|---- | -----|
| LOW | &quot;low&quot; |
| NORMAL | &quot;normal&quot; |
| HIGH | &quot;high&quot; |



