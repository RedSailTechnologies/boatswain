# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [gyarados.proto](#gyarados.proto)
    - [Delivery](#redsail.bosn.Delivery)
    - [Delivery.Application](#redsail.bosn.Delivery.Application)
    - [Deployment](#redsail.bosn.Deployment)
    - [Deployment.Docker](#redsail.bosn.Deployment.Docker)
    - [Deployment.Helm](#redsail.bosn.Deployment.Helm)
    - [Step](#redsail.bosn.Step)
    - [Step.StepAction](#redsail.bosn.Step.StepAction)
    - [Step.StepAction.Docker](#redsail.bosn.Step.StepAction.Docker)
    - [Step.StepAction.Helm](#redsail.bosn.Step.StepAction.Helm)
    - [Template](#redsail.bosn.Template)
    - [Trigger](#redsail.bosn.Trigger)
    - [Trigger.Approval](#redsail.bosn.Trigger.Approval)
    - [Trigger.Delivery](#redsail.bosn.Trigger.Delivery)
    - [Trigger.Web](#redsail.bosn.Trigger.Web)
  
    - [Step.StepAction.Helm.HelmAction](#redsail.bosn.Step.StepAction.Helm.HelmAction)
    - [Template.TemplateType](#redsail.bosn.Template.TemplateType)
  
- [Scalar Value Types](#scalar-value-types)



<a name="gyarados.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gyarados.proto
TODO


<a name="redsail.bosn.Delivery"></a>

### Delivery



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| version | [string](#string) |  |  |
| application | [Delivery.Application](#redsail.bosn.Delivery.Application) |  |  |
| clusters | [string](#string) | repeated |  |
| deployments | [Deployment](#redsail.bosn.Deployment) | repeated |  |
| tests | [Deployment](#redsail.bosn.Deployment) | repeated |  |
| triggers | [Trigger](#redsail.bosn.Trigger) | repeated |  |
| strategy | [Step](#redsail.bosn.Step) | repeated |  |






<a name="redsail.bosn.Delivery.Application"></a>

### Delivery.Application



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project | [string](#string) |  |  |






<a name="redsail.bosn.Deployment"></a>

### Deployment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| docker | [Deployment.Docker](#redsail.bosn.Deployment.Docker) |  |  |
| helm | [Deployment.Helm](#redsail.bosn.Deployment.Helm) |  |  |
| template | [string](#string) |  |  |
| arguments | [string](#string) |  |  |






<a name="redsail.bosn.Deployment.Docker"></a>

### Deployment.Docker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| image | [string](#string) |  |  |
| tag | [string](#string) |  |  |






<a name="redsail.bosn.Deployment.Helm"></a>

### Deployment.Helm



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chart | [string](#string) |  |  |
| repo | [string](#string) |  |  |
| version | [string](#string) |  |  |






<a name="redsail.bosn.Step"></a>

### Step



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| displayName | [string](#string) |  |  |
| success | [Step.StepAction](#redsail.bosn.Step.StepAction) | repeated |  |
| failure | [Step.StepAction](#redsail.bosn.Step.StepAction) | repeated |  |
| any | [Step.StepAction](#redsail.bosn.Step.StepAction) | repeated |  |
| always | [Step.StepAction](#redsail.bosn.Step.StepAction) | repeated |  |
| hold | [string](#string) |  |  |
| template | [string](#string) |  |  |
| arguments | [string](#string) |  |  |






<a name="redsail.bosn.Step.StepAction"></a>

### Step.StepAction



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| deployment | [string](#string) |  |  |
| test | [string](#string) |  |  |
| docker | [Step.StepAction.Docker](#redsail.bosn.Step.StepAction.Docker) |  |  |
| helm | [Step.StepAction.Helm](#redsail.bosn.Step.StepAction.Helm) |  |  |






<a name="redsail.bosn.Step.StepAction.Docker"></a>

### Step.StepAction.Docker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entrypoint | [string](#string) |  |  |
| rm | [bool](#bool) |  |  |
| env | [string](#string) |  |  |






<a name="redsail.bosn.Step.StepAction.Helm"></a>

### Step.StepAction.Helm



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [Step.StepAction.Helm.HelmAction](#redsail.bosn.Step.StepAction.Helm.HelmAction) |  |  |
| wait | [bool](#bool) |  |  |
| test | [bool](#bool) |  |  |
| values | [string](#string) |  |  |






<a name="redsail.bosn.Template"></a>

### Template



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| type | [Template.TemplateType](#redsail.bosn.Template.TemplateType) |  |  |
| yaml | [string](#string) |  |  |






<a name="redsail.bosn.Trigger"></a>

### Trigger



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| approval | [Trigger.Approval](#redsail.bosn.Trigger.Approval) |  |  |
| delivery | [Trigger.Delivery](#redsail.bosn.Trigger.Delivery) |  |  |
| manual | [Trigger.Approval](#redsail.bosn.Trigger.Approval) |  |  |
| web | [Trigger.Web](#redsail.bosn.Trigger.Web) |  |  |
| template | [string](#string) |  |  |
| arguments | [string](#string) |  |  |






<a name="redsail.bosn.Trigger.Approval"></a>

### Trigger.Approval



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| groups | [string](#string) | repeated |  |
| users | [string](#string) | repeated |  |
| params | [string](#string) | repeated |  |






<a name="redsail.bosn.Trigger.Delivery"></a>

### Trigger.Delivery



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| trigger | [string](#string) |  |  |






<a name="redsail.bosn.Trigger.Web"></a>

### Trigger.Web



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| params | [string](#string) | repeated |  |





 


<a name="redsail.bosn.Step.StepAction.Helm.HelmAction"></a>

### Step.StepAction.Helm.HelmAction


| Name | Number | Description |
| ---- | ------ | ----------- |
| install | 0 |  |
| upgrade | 1 |  |
| rollback | 2 |  |
| uninstall | 3 |  |



<a name="redsail.bosn.Template.TemplateType"></a>

### Template.TemplateType


| Name | Number | Description |
| ---- | ------ | ----------- |
| deployment | 0 |  |
| step | 1 |  |
| trigger | 2 |  |


 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

