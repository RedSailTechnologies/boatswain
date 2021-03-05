# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [agent.proto](#agent.proto)
    - [Action](#redsail.bosn.Action)
    - [ActionsRead](#redsail.bosn.ActionsRead)
    - [ReadActions](#redsail.bosn.ReadActions)
    - [Result](#redsail.bosn.Result)
    - [ResultReturned](#redsail.bosn.ResultReturned)
    - [ReturnResult](#redsail.bosn.ReturnResult)
  
    - [ActionType](#redsail.bosn.ActionType)
  
    - [Agent](#redsail.bosn.Agent)
    - [AgentAction](#redsail.bosn.AgentAction)
  
- [Scalar Value Types](#scalar-value-types)



<a name="agent.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## agent.proto
Agent is the service for external clusters to call into to register and receive actions.


<a name="redsail.bosn.Action"></a>

### Action



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| cluster_uuid | [string](#string) |  |  |
| cluster_token | [string](#string) |  |  |
| action_type | [ActionType](#redsail.bosn.ActionType) |  |  |
| action | [string](#string) |  |  |
| timeout_seconds | [int64](#int64) |  |  |
| args | [bytes](#bytes) |  |  |






<a name="redsail.bosn.ActionsRead"></a>

### ActionsRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| actions | [Action](#redsail.bosn.Action) | repeated |  |






<a name="redsail.bosn.ReadActions"></a>

### ReadActions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cluster_uuid | [string](#string) |  |  |
| cluster_token | [string](#string) |  |  |






<a name="redsail.bosn.Result"></a>

### Result



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  |  |
| error | [string](#string) |  |  |






<a name="redsail.bosn.ResultReturned"></a>

### ResultReturned







<a name="redsail.bosn.ReturnResult"></a>

### ReturnResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action_uuid | [string](#string) |  |  |
| cluster_uuid | [string](#string) |  |  |
| cluster_token | [string](#string) |  |  |
| result | [Result](#redsail.bosn.Result) |  |  |





 


<a name="redsail.bosn.ActionType"></a>

### ActionType


| Name | Number | Description |
| ---- | ------ | ----------- |
| HELM_ACTION | 0 |  |
| KUBE_ACTION | 1 |  |


 

 


<a name="redsail.bosn.Agent"></a>

### Agent


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Actions | [ReadActions](#redsail.bosn.ReadActions) | [ActionsRead](#redsail.bosn.ActionsRead) | gets the next action for the agent or an empty list if there&#39;s nothing to do |
| Results | [ReturnResult](#redsail.bosn.ReturnResult) | [ResultReturned](#redsail.bosn.ResultReturned) | returns a result for this agent |


<a name="redsail.bosn.AgentAction"></a>

### AgentAction


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Run | [Action](#redsail.bosn.Action) | [Result](#redsail.bosn.Result) |  |

 



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

