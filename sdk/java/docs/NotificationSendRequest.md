

# NotificationSendRequest


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**projectId** | **UUID** | Project ID |  |
|**senderId** | **UUID** | Sender user ID |  [optional] |
|**messageType** | [**MessageTypeEnum**](#MessageTypeEnum) | Type of message |  |
|**content** | **String** | Message content |  |
|**mentions** | **List&lt;String&gt;** | List of user mentions |  [optional] |
|**priority** | [**PriorityEnum**](#PriorityEnum) | Message priority |  [optional] |
|**targets** | **List&lt;String&gt;** | List of target destination IDs |  |
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



