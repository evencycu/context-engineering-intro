

# DestinationCreateRequest


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**projectId** | **UUID** | Project ID |  |
|**botId** | **UUID** | Bot ID |  |
|**type** | [**TypeEnum**](#TypeEnum) | Destination type |  |
|**targets** | [**List&lt;DestinationTarget&gt;**](DestinationTarget.md) | List of targets |  |
|**status** | [**StatusEnum**](#StatusEnum) | Destination status |  [optional] |



## Enum: TypeEnum

| Name | Value |
|---- | -----|
| PERSONAL | &quot;personal&quot; |
| GROUP | &quot;group&quot; |
| CHANNEL | &quot;channel&quot; |



## Enum: StatusEnum

| Name | Value |
|---- | -----|
| ACTIVE | &quot;active&quot; |
| INACTIVE | &quot;inactive&quot; |



