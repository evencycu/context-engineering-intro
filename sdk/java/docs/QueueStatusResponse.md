

# QueueStatusResponse


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**status** | [**StatusEnum**](#StatusEnum) |  |  [optional] |
|**workers** | **Integer** |  |  [optional] |
|**activeWorkers** | **Integer** |  |  [optional] |
|**queueSize** | **Integer** |  |  [optional] |
|**lastProcessed** | **OffsetDateTime** |  |  [optional] |



## Enum: StatusEnum

| Name | Value |
|---- | -----|
| RUNNING | &quot;running&quot; |
| PAUSED | &quot;paused&quot; |
| STOPPED | &quot;stopped&quot; |
| ERROR | &quot;error&quot; |



