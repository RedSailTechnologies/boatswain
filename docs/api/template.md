# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [template.proto](#template.proto)
    - [CreateTemplate](#redsail.bosn.CreateTemplate)
    - [DestroyTemplate](#redsail.bosn.DestroyTemplate)
    - [ReadTemplate](#redsail.bosn.ReadTemplate)
    - [ReadTemplates](#redsail.bosn.ReadTemplates)
    - [TemplateCreated](#redsail.bosn.TemplateCreated)
    - [TemplateDestroyed](#redsail.bosn.TemplateDestroyed)
    - [TemplateRead](#redsail.bosn.TemplateRead)
    - [TemplateUpdated](#redsail.bosn.TemplateUpdated)
    - [TemplatesRead](#redsail.bosn.TemplatesRead)
    - [UpdateTemplate](#redsail.bosn.UpdateTemplate)
  
    - [TemplateType](#redsail.bosn.TemplateType)
  
    - [Template](#redsail.bosn.Template)
  
- [Scalar Value Types](#scalar-value-types)



<a name="template.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## template.proto
Template is the service for dealing with yaml templates.


<a name="redsail.bosn.CreateTemplate"></a>

### CreateTemplate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of the template |
| type | [TemplateType](#redsail.bosn.TemplateType) |  | the type of the template |
| yaml | [string](#string) |  | the template |






<a name="redsail.bosn.DestroyTemplate"></a>

### DestroyTemplate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the template |






<a name="redsail.bosn.ReadTemplate"></a>

### ReadTemplate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the template |






<a name="redsail.bosn.ReadTemplates"></a>

### ReadTemplates







<a name="redsail.bosn.TemplateCreated"></a>

### TemplateCreated







<a name="redsail.bosn.TemplateDestroyed"></a>

### TemplateDestroyed







<a name="redsail.bosn.TemplateRead"></a>

### TemplateRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the template |
| name | [string](#string) |  | name of the template |
| type | [TemplateType](#redsail.bosn.TemplateType) |  | the type of the template |
| yaml | [string](#string) |  | the template |






<a name="redsail.bosn.TemplateUpdated"></a>

### TemplateUpdated







<a name="redsail.bosn.TemplatesRead"></a>

### TemplatesRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| templates | [TemplateRead](#redsail.bosn.TemplateRead) | repeated | templates read |






<a name="redsail.bosn.UpdateTemplate"></a>

### UpdateTemplate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the template |
| name | [string](#string) |  | name of the template |
| type | [TemplateType](#redsail.bosn.TemplateType) |  | the type of the template |
| yaml | [string](#string) |  | the template |





 


<a name="redsail.bosn.TemplateType"></a>

### TemplateType


| Name | Number | Description |
| ---- | ------ | ----------- |
| DEPLOYMENT | 0 |  |
| STEP | 1 |  |
| TRIGGER | 2 |  |


 

 


<a name="redsail.bosn.Template"></a>

### Template


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [CreateTemplate](#redsail.bosn.CreateTemplate) | [TemplateCreated](#redsail.bosn.TemplateCreated) | adds a template |
| Update | [UpdateTemplate](#redsail.bosn.UpdateTemplate) | [TemplateUpdated](#redsail.bosn.TemplateUpdated) | edits an existing template |
| Destroy | [DestroyTemplate](#redsail.bosn.DestroyTemplate) | [TemplateDestroyed](#redsail.bosn.TemplateDestroyed) | removes a template |
| Read | [ReadTemplate](#redsail.bosn.ReadTemplate) | [TemplateRead](#redsail.bosn.TemplateRead) | reads a template |
| All | [ReadTemplates](#redsail.bosn.ReadTemplates) | [TemplatesRead](#redsail.bosn.TemplatesRead) | gets all templates |

 



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

