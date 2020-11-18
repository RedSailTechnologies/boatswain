# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [application.proto](#application.proto)
    - [ApplicationCluster](#redsail.bosn.ApplicationCluster)
    - [ApplicationRead](#redsail.bosn.ApplicationRead)
    - [ApplicationsRead](#redsail.bosn.ApplicationsRead)
    - [ReadApplications](#redsail.bosn.ReadApplications)
  
    - [Application](#redsail.bosn.Application)
  
- [Scalar Value Types](#scalar-value-types)



<a name="application.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## application.proto
Kraken is the service managing external cluster connections.
The api can be hit at /api/redsail.bosn.Kraken/&lt;Method&gt;.


<a name="redsail.bosn.ApplicationCluster"></a>

### ApplicationCluster



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cluster_name | [string](#string) |  | the cluster name |
| version | [string](#string) |  | the app version by label app.kubernetes.io/version |
| namespace | [string](#string) |  | the namespace |
| ready | [bool](#bool) |  | whether all deployment or ss pods are ready |






<a name="redsail.bosn.ApplicationRead"></a>

### ApplicationRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | the application name by label app.kubernetes.io/name |
| project | [string](#string) |  | the project by label app.kubernetes.io/part-of |
| clusters | [ApplicationCluster](#redsail.bosn.ApplicationCluster) | repeated | the list of isntances of this application by cluster |






<a name="redsail.bosn.ApplicationsRead"></a>

### ApplicationsRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| applications | [ApplicationRead](#redsail.bosn.ApplicationRead) | repeated | the list of applications |






<a name="redsail.bosn.ReadApplications"></a>

### ReadApplications






 

 

 


<a name="redsail.bosn.Application"></a>

### Application


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| All | [ReadApplications](#redsail.bosn.ReadApplications) | [ApplicationsRead](#redsail.bosn.ApplicationsRead) | gets all applications currently found in each cluster and their status |

 



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

