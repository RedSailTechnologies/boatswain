# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [kraken/service.proto](#kraken/service.proto)
    - [Cluster](#redsail.bosn.Cluster)
    - [ClustersRequest](#redsail.bosn.ClustersRequest)
    - [ClustersResponse](#redsail.bosn.ClustersResponse)
    - [Deployment](#redsail.bosn.Deployment)
    - [DeploymentsRequest](#redsail.bosn.DeploymentsRequest)
    - [DeploymentsResponse](#redsail.bosn.DeploymentsResponse)
  
    - [Kraken](#redsail.bosn.Kraken)
  
- [Scalar Value Types](#scalar-value-types)



<a name="kraken/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## kraken/service.proto
Kraken is the service managing external cluster connections.
The api can be hit at /api/redsail.bosn.Kraken/&lt;Method&gt; if
using a json client.


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






<a name="redsail.bosn.Deployment"></a>

### Deployment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of the deployment |
| namespace | [string](#string) |  | namespace of the deployment |
| ready | [bool](#bool) |  | if the deployment is ready |
| version | [string](#string) |  | version of the deployment |
| cluster | [Cluster](#redsail.bosn.Cluster) |  | the the cluster deployed to |






<a name="redsail.bosn.DeploymentsRequest"></a>

### DeploymentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cluster | [Cluster](#redsail.bosn.Cluster) |  | the cluster to get deployments for |






<a name="redsail.bosn.DeploymentsResponse"></a>

### DeploymentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deployments | [Deployment](#redsail.bosn.Deployment) | repeated | the list of deployments |





 

 

 


<a name="redsail.bosn.Kraken"></a>

### Kraken


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Clusters | [ClustersRequest](#redsail.bosn.ClustersRequest) | [ClustersResponse](#redsail.bosn.ClustersResponse) | gets all clusters currently configured and their status |
| ClusterStatus | [Cluster](#redsail.bosn.Cluster) | [Cluster](#redsail.bosn.Cluster) | gets the status for a single cluster |
| Deployments | [DeploymentsRequest](#redsail.bosn.DeploymentsRequest) | [DeploymentsResponse](#redsail.bosn.DeploymentsResponse) | gets all the deployments for a particular cluster |
| DeploymentStatus | [Deployment](#redsail.bosn.Deployment) | [Deployment](#redsail.bosn.Deployment) | gets a |

 



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

