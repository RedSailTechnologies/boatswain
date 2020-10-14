# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [poseidon.proto](#poseidon.proto)
    - [Chart](#redsail.bosn.Chart)
    - [ChartVersion](#redsail.bosn.ChartVersion)
    - [ChartsResponse](#redsail.bosn.ChartsResponse)
    - [Repo](#redsail.bosn.Repo)
    - [ReposRequest](#redsail.bosn.ReposRequest)
    - [ReposResponse](#redsail.bosn.ReposResponse)
  
    - [Poseidon](#redsail.bosn.Poseidon)
  
- [Scalar Value Types](#scalar-value-types)



<a name="poseidon.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## poseidon.proto
Poseidon is the service concerned with managing helm repos and charts.
The api can be hit at /api/redsail.bosn.Poseidon/&lt;Method&gt;.


<a name="redsail.bosn.Chart"></a>

### Chart



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | the chart name |
| versions | [ChartVersion](#redsail.bosn.ChartVersion) | repeated | the versions available for this chart |






<a name="redsail.bosn.ChartVersion"></a>

### ChartVersion



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chart_version | [string](#string) |  | the chart version |
| app_version | [string](#string) |  | the chart&#39;s default app version |
| description | [string](#string) |  | description of the chart |
| url | [string](#string) |  | the url for this specific version of the chart |






<a name="redsail.bosn.ChartsResponse"></a>

### ChartsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| charts | [Chart](#redsail.bosn.Chart) | repeated | the list of charts |






<a name="redsail.bosn.Repo"></a>

### Repo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of the repo |
| endpoint | [string](#string) |  | the endpoint for the repo |
| ready | [bool](#bool) |  | if the repo is ready (checking by getting index.yaml) |






<a name="redsail.bosn.ReposRequest"></a>

### ReposRequest







<a name="redsail.bosn.ReposResponse"></a>

### ReposResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| repos | [Repo](#redsail.bosn.Repo) | repeated | the list of currently configured repositories |





 

 

 


<a name="redsail.bosn.Poseidon"></a>

### Poseidon


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Charts | [Repo](#redsail.bosn.Repo) | [ChartsResponse](#redsail.bosn.ChartsResponse) | gets all the charts for this repository |
| Repos | [ReposRequest](#redsail.bosn.ReposRequest) | [ReposResponse](#redsail.bosn.ReposResponse) | gets all the currently configured repositories |

 



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

