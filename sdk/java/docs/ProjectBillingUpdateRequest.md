

# ProjectBillingUpdateRequest


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**paymentMethod** | [**PaymentMethodEnum**](#PaymentMethodEnum) |  |  [optional] |
|**billingCycle** | [**BillingCycleEnum**](#BillingCycleEnum) |  |  [optional] |
|**nextBillingDate** | **LocalDate** |  |  [optional] |
|**billingStatus** | [**BillingStatusEnum**](#BillingStatusEnum) |  |  [optional] |



## Enum: PaymentMethodEnum

| Name | Value |
|---- | -----|
| CREDIT_CARD | &quot;credit_card&quot; |
| BANK_TRANSFER | &quot;bank_transfer&quot; |
| PAYPAL | &quot;paypal&quot; |



## Enum: BillingCycleEnum

| Name | Value |
|---- | -----|
| MONTHLY | &quot;monthly&quot; |
| YEARLY | &quot;yearly&quot; |



## Enum: BillingStatusEnum

| Name | Value |
|---- | -----|
| ACTIVE | &quot;active&quot; |
| SUSPENDED | &quot;suspended&quot; |
| CANCELLED | &quot;cancelled&quot; |



