# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [trigger.proto](#trigger.proto)
    - [ManualTriggered](#redsail.bosn.ManualTriggered)
    - [ReadStatus](#redsail.bosn.ReadStatus)
    - [StatusRead](#redsail.bosn.StatusRead)
    - [TriggerManual](#redsail.bosn.TriggerManual)
    - [TriggerWeb](#redsail.bosn.TriggerWeb)
    - [WebTriggered](#redsail.bosn.WebTriggered)
  
    - [TriggerStatus](#redsail.bosn.TriggerStatus)
  
    - [Trigger](#redsail.bosn.Trigger)
  
- [Scalar Value Types](#scalar-value-types)



<a name="trigger.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## trigger.proto
Trigger is the service for creating triggers to start deployments.


<a name="redsail.bosn.ManualTriggered"></a>

### ManualTriggered



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| run_uuid | [string](#string) |  |  |






<a name="redsail.bosn.ReadStatus"></a>

### ReadStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deployment_uuid | [string](#string) |  |  |
| deployment_token | [string](#string) |  |  |
| run_uuid | [string](#string) |  |  |






<a name="redsail.bosn.StatusRead"></a>

### StatusRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [TriggerStatus](#redsail.bosn.TriggerStatus) |  |  |






<a name="redsail.bosn.TriggerManual"></a>

### TriggerManual



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| args | [string](#string) |  |  |






<a name="redsail.bosn.TriggerWeb"></a>

### TriggerWeb



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| token | [string](#string) |  |  |
| args | [string](#string) |  |  |






<a name="redsail.bosn.WebTriggered"></a>

### WebTriggered



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| run_uuid | [string](#string) |  |  |





 


<a name="redsail.bosn.TriggerStatus"></a>

### TriggerStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| NOT_STARTED | 0 |  |
| IN_PROGRESS | 1 |  |
| AWAITING_APPROVAL | 2 |  |
| FAILED | 3 |  |
| SUCCEEDED | 4 |  |
| SKIPPED | 5 |  |


 

 


<a name="redsail.bosn.Trigger"></a>

### Trigger


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Manual | [TriggerManual](#redsail.bosn.TriggerManual) | [ManualTriggered](#redsail.bosn.ManualTriggered) | triggers a deployment manually |
| Web | [TriggerWeb](#redsail.bosn.TriggerWeb) | [WebTriggered](#redsail.bosn.WebTriggered) | triggers a deployment from a web call |
| Status | [ReadStatus](#redsail.bosn.ReadStatus) | [StatusRead](#redsail.bosn.StatusRead) | gets the status of a run |

 



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

