# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [deployment.proto](#deployment.proto)
    - [CreateDeployment](#redsail.bosn.CreateDeployment)
    - [DeploymentCreated](#redsail.bosn.DeploymentCreated)
    - [DeploymentDestroyed](#redsail.bosn.DeploymentDestroyed)
    - [DeploymentRead](#redsail.bosn.DeploymentRead)
    - [DeploymentReadSummary](#redsail.bosn.DeploymentReadSummary)
    - [DeploymentUpdated](#redsail.bosn.DeploymentUpdated)
    - [DeploymentsRead](#redsail.bosn.DeploymentsRead)
    - [DestroyDeployment](#redsail.bosn.DestroyDeployment)
    - [ReadDeployment](#redsail.bosn.ReadDeployment)
    - [ReadDeployments](#redsail.bosn.ReadDeployments)
    - [UpdateDeployment](#redsail.bosn.UpdateDeployment)
  
    - [Deployment](#redsail.bosn.Deployment)
  
- [Scalar Value Types](#scalar-value-types)



<a name="deployment.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## deployment.proto
Deployment is the service for creation and management of application installs/upgrades.


<a name="redsail.bosn.CreateDeployment"></a>

### CreateDeployment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | the name of this deployment |
| repo_id | [string](#string) |  | the unique id of the repo to get the deployment yaml from |
| branch | [string](#string) |  | the branch from the repo to get the file from |
| file_path | [string](#string) |  | the path to the deployment file |






<a name="redsail.bosn.DeploymentCreated"></a>

### DeploymentCreated







<a name="redsail.bosn.DeploymentDestroyed"></a>

### DeploymentDestroyed







<a name="redsail.bosn.DeploymentRead"></a>

### DeploymentRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |
| name | [string](#string) |  | the name of this deployment |
| repo_id | [string](#string) |  | the unique id of the repo to get the deployment yaml from |
| repo_name | [string](#string) |  | the name of the repo |
| branch | [string](#string) |  | the branch from the repo to get the file from |
| file_path | [string](#string) |  | the path to the deployment file |
| yaml | [string](#string) |  | the templated yaml of this deployment |






<a name="redsail.bosn.DeploymentReadSummary"></a>

### DeploymentReadSummary



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |
| name | [string](#string) |  | name of the deployment |
| repo_id | [string](#string) |  | the name of the repo |
| repo_name | [string](#string) |  | the name of the repo |
| branch | [string](#string) |  | the branch from the repo to get the file from |
| file_path | [string](#string) |  | the path to the deployment file |






<a name="redsail.bosn.DeploymentUpdated"></a>

### DeploymentUpdated







<a name="redsail.bosn.DeploymentsRead"></a>

### DeploymentsRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deployments | [DeploymentReadSummary](#redsail.bosn.DeploymentReadSummary) | repeated | the list of deployments |






<a name="redsail.bosn.DestroyDeployment"></a>

### DestroyDeployment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |






<a name="redsail.bosn.ReadDeployment"></a>

### ReadDeployment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |






<a name="redsail.bosn.ReadDeployments"></a>

### ReadDeployments







<a name="redsail.bosn.UpdateDeployment"></a>

### UpdateDeployment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |
| name | [string](#string) |  | the name of this deployment |
| repo_id | [string](#string) |  | the unique id of the repo to get the deployment yaml from |
| branch | [string](#string) |  | the branch from the repo to get the file from |
| file_path | [string](#string) |  | the path to the deployment file |





 

 

 


<a name="redsail.bosn.Deployment"></a>

### Deployment


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [CreateDeployment](#redsail.bosn.CreateDeployment) | [DeploymentCreated](#redsail.bosn.DeploymentCreated) | creates a new delivery |
| Update | [UpdateDeployment](#redsail.bosn.UpdateDeployment) | [DeploymentUpdated](#redsail.bosn.DeploymentUpdated) | edits an already existing deployment |
| Destroy | [DestroyDeployment](#redsail.bosn.DestroyDeployment) | [DeploymentDestroyed](#redsail.bosn.DeploymentDestroyed) | removes a deployment from the list of configurations |
| Read | [ReadDeployment](#redsail.bosn.ReadDeployment) | [DeploymentRead](#redsail.bosn.DeploymentRead) | reads out a deployment |
| All | [ReadDeployments](#redsail.bosn.ReadDeployments) | [DeploymentsRead](#redsail.bosn.DeploymentsRead) | gets all deployments currently configured and their status |

 



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

