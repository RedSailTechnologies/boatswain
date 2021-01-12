# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [cluster.proto](#cluster.proto)
    - [ClusterCreated](#redsail.bosn.ClusterCreated)
    - [ClusterDestroyed](#redsail.bosn.ClusterDestroyed)
    - [ClusterFound](#redsail.bosn.ClusterFound)
    - [ClusterRead](#redsail.bosn.ClusterRead)
    - [ClusterUpdated](#redsail.bosn.ClusterUpdated)
    - [ClustersRead](#redsail.bosn.ClustersRead)
    - [CreateCluster](#redsail.bosn.CreateCluster)
    - [DestroyCluster](#redsail.bosn.DestroyCluster)
    - [FindCluster](#redsail.bosn.FindCluster)
    - [ReadCluster](#redsail.bosn.ReadCluster)
    - [ReadClusters](#redsail.bosn.ReadClusters)
    - [UpdateCluster](#redsail.bosn.UpdateCluster)
  
    - [Cluster](#redsail.bosn.Cluster)
  
- [Scalar Value Types](#scalar-value-types)



<a name="cluster.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cluster.proto
Cluster is the service managing external clusters.


<a name="redsail.bosn.ClusterCreated"></a>

### ClusterCreated







<a name="redsail.bosn.ClusterDestroyed"></a>

### ClusterDestroyed







<a name="redsail.bosn.ClusterFound"></a>

### ClusterFound



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the cluster found |






<a name="redsail.bosn.ClusterRead"></a>

### ClusterRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the cluster |
| name | [string](#string) |  | name of the cluster |
| endpoint | [string](#string) |  | api server endpoint |
| token | [string](#string) |  | authentication token |
| cert | [string](#string) |  | server certificate |
| ready | [bool](#bool) |  | server ready status, based on node status |






<a name="redsail.bosn.ClusterUpdated"></a>

### ClusterUpdated







<a name="redsail.bosn.ClustersRead"></a>

### ClustersRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clusters | [ClusterRead](#redsail.bosn.ClusterRead) | repeated | clusters read |






<a name="redsail.bosn.CreateCluster"></a>

### CreateCluster



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of the cluster |
| endpoint | [string](#string) |  | api server endpoint |
| token | [string](#string) |  | authentication token |
| cert | [string](#string) |  | server certificate |






<a name="redsail.bosn.DestroyCluster"></a>

### DestroyCluster



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the cluster |






<a name="redsail.bosn.FindCluster"></a>

### FindCluster



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | the name of the cluster to search for |






<a name="redsail.bosn.ReadCluster"></a>

### ReadCluster



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the cluster |






<a name="redsail.bosn.ReadClusters"></a>

### ReadClusters







<a name="redsail.bosn.UpdateCluster"></a>

### UpdateCluster



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the cluster |
| name | [string](#string) |  | name of the cluster |
| endpoint | [string](#string) |  | api server endpoint |
| token | [string](#string) |  | authentication token |
| cert | [string](#string) |  | server certificate |





 

 

 


<a name="redsail.bosn.Cluster"></a>

### Cluster


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [CreateCluster](#redsail.bosn.CreateCluster) | [ClusterCreated](#redsail.bosn.ClusterCreated) | adds a cluster to the list of configurations |
| Update | [UpdateCluster](#redsail.bosn.UpdateCluster) | [ClusterUpdated](#redsail.bosn.ClusterUpdated) | edits an already existing cluster |
| Destroy | [DestroyCluster](#redsail.bosn.DestroyCluster) | [ClusterDestroyed](#redsail.bosn.ClusterDestroyed) | removes a cluster from the list of configurations |
| Read | [ReadCluster](#redsail.bosn.ReadCluster) | [ClusterRead](#redsail.bosn.ClusterRead) | reads out a cluster |
| Find | [FindCluster](#redsail.bosn.FindCluster) | [ClusterFound](#redsail.bosn.ClusterFound) | finds the cluster uuid by name |
| All | [ReadClusters](#redsail.bosn.ReadClusters) | [ClustersRead](#redsail.bosn.ClustersRead) | gets all clusters currently configured and their status |

 



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

