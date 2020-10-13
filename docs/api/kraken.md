# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [service.proto](#service.proto)
    - [Cluster](#redsail.bosn.Cluster)
    - [ClustersRequest](#redsail.bosn.ClustersRequest)
    - [ClustersResponse](#redsail.bosn.ClustersResponse)
    - [Release](#redsail.bosn.Release)
    - [ReleaseRequest](#redsail.bosn.ReleaseRequest)
    - [ReleaseResponse](#redsail.bosn.ReleaseResponse)
    - [Releases](#redsail.bosn.Releases)
  
    - [Status](#redsail.bosn.Status)
  
    - [Kraken](#redsail.bosn.Kraken)
  
- [Scalar Value Types](#scalar-value-types)



<a name="service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## service.proto
Kraken is the service managing external cluster connections.
The api can be hit at /api/redsail.bosn.Kraken/&lt;Method&gt;.


<a name="redsail.bosn.Cluster"></a>

### Cluster



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of the cluster |
| endpoint | [string](#string) |  | the cluster&#39;s api server |
| ready | [bool](#bool) |  | if the cluster is ready (checking each node for Ready status) |






<a name="redsail.bosn.ClustersRequest"></a>

### ClustersRequest







<a name="redsail.bosn.ClustersResponse"></a>

### ClustersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clusters | [Cluster](#redsail.bosn.Cluster) | repeated | the list of clusters |






<a name="redsail.bosn.Release"></a>

### Release



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| namespace | [string](#string) |  |  |
| app_version | [string](#string) |  |  |
| chart_version | [string](#string) |  |  |
| cluster_name | [string](#string) |  |  |
| status | [string](#string) |  |  |






<a name="redsail.bosn.ReleaseRequest"></a>

### ReleaseRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clusters | [Cluster](#redsail.bosn.Cluster) | repeated | the clusters to get apps for |






<a name="redsail.bosn.ReleaseResponse"></a>

### ReleaseResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| release_lists | [Releases](#redsail.bosn.Releases) | repeated |  |






<a name="redsail.bosn.Releases"></a>

### Releases



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| chart | [string](#string) |  |  |
| releases | [Release](#redsail.bosn.Release) | repeated |  |





 


<a name="redsail.bosn.Status"></a>

### Status
the helm status of the release

| Name | Number | Description |
| ---- | ------ | ----------- |
| unknown | 0 |  |
| deployed | 1 |  |
| uninstalled | 2 |  |
| superseded | 3 |  |
| failed | 4 |  |
| uninstalling | 5 |  |
| pending_install | 6 |  |
| pending_upgrade | 7 |  |
| pending_rollback | 8 |  |


 

 


<a name="redsail.bosn.Kraken"></a>

### Kraken


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Clusters | [ClustersRequest](#redsail.bosn.ClustersRequest) | [ClustersResponse](#redsail.bosn.ClustersResponse) | gets all clusters currently configured and their status |
| ClusterStatus | [Cluster](#redsail.bosn.Cluster) | [Cluster](#redsail.bosn.Cluster) | gets the status for a single cluster |
| Releases | [ReleaseRequest](#redsail.bosn.ReleaseRequest) | [ReleaseResponse](#redsail.bosn.ReleaseResponse) | gets all applications for the clusters passed |
| ReleaseStatus | [Release](#redsail.bosn.Release) | [Release](#redsail.bosn.Release) | gets the status for a single application in a single cluster |

 



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

