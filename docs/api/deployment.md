# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [deployment.proto](#deployment.proto)
    - [ApprovalRead](#redsail.bosn.ApprovalRead)
    - [ApprovalsRead](#redsail.bosn.ApprovalsRead)
    - [ApproveStep](#redsail.bosn.ApproveStep)
    - [CreateDeployment](#redsail.bosn.CreateDeployment)
    - [DeploymentCreated](#redsail.bosn.DeploymentCreated)
    - [DeploymentDestroyed](#redsail.bosn.DeploymentDestroyed)
    - [DeploymentRead](#redsail.bosn.DeploymentRead)
    - [DeploymentReadSummary](#redsail.bosn.DeploymentReadSummary)
    - [DeploymentTemplated](#redsail.bosn.DeploymentTemplated)
    - [DeploymentTokenRead](#redsail.bosn.DeploymentTokenRead)
    - [DeploymentUpdated](#redsail.bosn.DeploymentUpdated)
    - [DeploymentsRead](#redsail.bosn.DeploymentsRead)
    - [DestroyDeployment](#redsail.bosn.DestroyDeployment)
    - [LinkRead](#redsail.bosn.LinkRead)
    - [ReadApprovals](#redsail.bosn.ReadApprovals)
    - [ReadDeployment](#redsail.bosn.ReadDeployment)
    - [ReadDeploymentToken](#redsail.bosn.ReadDeploymentToken)
    - [ReadDeployments](#redsail.bosn.ReadDeployments)
    - [ReadRun](#redsail.bosn.ReadRun)
    - [ReadRuns](#redsail.bosn.ReadRuns)
    - [RunRead](#redsail.bosn.RunRead)
    - [RunReadSummary](#redsail.bosn.RunReadSummary)
    - [RunsRead](#redsail.bosn.RunsRead)
    - [StepApproved](#redsail.bosn.StepApproved)
    - [StepLog](#redsail.bosn.StepLog)
    - [StepRead](#redsail.bosn.StepRead)
    - [TemplateDeployment](#redsail.bosn.TemplateDeployment)
    - [UpdateDeployment](#redsail.bosn.UpdateDeployment)
  
    - [LogLevel](#redsail.bosn.LogLevel)
    - [Status](#redsail.bosn.Status)
  
    - [Deployment](#redsail.bosn.Deployment)
  
- [Scalar Value Types](#scalar-value-types)



<a name="deployment.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## deployment.proto
Deployment is the service for creation and management of application installs/upgrades.


<a name="redsail.bosn.ApprovalRead"></a>

### ApprovalRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| run_uuid | [string](#string) |  |  |
| run_name | [string](#string) |  |  |
| run_version | [string](#string) |  |  |
| step_name | [string](#string) |  |  |






<a name="redsail.bosn.ApprovalsRead"></a>

### ApprovalsRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| approvals | [ApprovalRead](#redsail.bosn.ApprovalRead) | repeated |  |






<a name="redsail.bosn.ApproveStep"></a>

### ApproveStep



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| run_uuid | [string](#string) |  |  |
| approve | [bool](#bool) |  |  |
| override | [bool](#bool) |  |  |






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






<a name="redsail.bosn.DeploymentReadSummary"></a>

### DeploymentReadSummary



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |
| name | [string](#string) |  | the name of this deployment |






<a name="redsail.bosn.DeploymentTemplated"></a>

### DeploymentTemplated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |
| yaml | [string](#string) |  | the templated yaml for this deployment |






<a name="redsail.bosn.DeploymentTokenRead"></a>

### DeploymentTokenRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  | deployment token for web calls |






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






<a name="redsail.bosn.LinkRead"></a>

### LinkRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| url | [string](#string) |  |  |






<a name="redsail.bosn.ReadApprovals"></a>

### ReadApprovals







<a name="redsail.bosn.ReadDeployment"></a>

### ReadDeployment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |






<a name="redsail.bosn.ReadDeploymentToken"></a>

### ReadDeploymentToken



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |






<a name="redsail.bosn.ReadDeployments"></a>

### ReadDeployments







<a name="redsail.bosn.ReadRun"></a>

### ReadRun



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deployment_uuid | [string](#string) |  | unique id of the run |






<a name="redsail.bosn.ReadRuns"></a>

### ReadRuns



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deployment_uuid | [string](#string) |  | unique id of the deployment to get runs for |






<a name="redsail.bosn.RunRead"></a>

### RunRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| version | [string](#string) |  |  |
| status | [Status](#redsail.bosn.Status) |  |  |
| start_time | [int64](#int64) |  |  |
| stop_time | [int64](#int64) |  |  |
| links | [LinkRead](#redsail.bosn.LinkRead) | repeated |  |
| steps | [StepRead](#redsail.bosn.StepRead) | repeated |  |






<a name="redsail.bosn.RunReadSummary"></a>

### RunReadSummary



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| name | [string](#string) |  |  |
| version | [string](#string) |  |  |
| status | [Status](#redsail.bosn.Status) |  |  |
| start_time | [int64](#int64) |  |  |
| stop_time | [int64](#int64) |  |  |






<a name="redsail.bosn.RunsRead"></a>

### RunsRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| runs | [RunReadSummary](#redsail.bosn.RunReadSummary) | repeated | the runs |






<a name="redsail.bosn.StepApproved"></a>

### StepApproved







<a name="redsail.bosn.StepLog"></a>

### StepLog



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| timestamp | [int64](#int64) |  |  |
| level | [LogLevel](#redsail.bosn.LogLevel) |  |  |
| message | [string](#string) |  |  |






<a name="redsail.bosn.StepRead"></a>

### StepRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| status | [Status](#redsail.bosn.Status) |  |  |
| start_time | [int64](#int64) |  |  |
| stop_time | [int64](#int64) |  |  |
| logs | [StepLog](#redsail.bosn.StepLog) | repeated |  |






<a name="redsail.bosn.TemplateDeployment"></a>

### TemplateDeployment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |






<a name="redsail.bosn.UpdateDeployment"></a>

### UpdateDeployment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the deployment |
| name | [string](#string) |  | the name of this deployment |
| repo_id | [string](#string) |  | the unique id of the repo to get the deployment yaml from |
| branch | [string](#string) |  | the branch from the repo to get the file from |
| file_path | [string](#string) |  | the path to the deployment file |





 


<a name="redsail.bosn.LogLevel"></a>

### LogLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| DEBUG | 0 |  |
| INFO | 1 |  |
| WARN | 2 |  |
| ERROR | 3 |  |



<a name="redsail.bosn.Status"></a>

### Status


| Name | Number | Description |
| ---- | ------ | ----------- |
| NOT_STARTED | 0 |  |
| IN_PROGRESS | 1 |  |
| AWAITING_APPROVAL | 2 |  |
| FAILED | 3 |  |
| SUCCEEDED | 4 |  |
| SKIPPED | 5 |  |


 

 


<a name="redsail.bosn.Deployment"></a>

### Deployment


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [CreateDeployment](#redsail.bosn.CreateDeployment) | [DeploymentCreated](#redsail.bosn.DeploymentCreated) | creates a new deployment |
| Update | [UpdateDeployment](#redsail.bosn.UpdateDeployment) | [DeploymentUpdated](#redsail.bosn.DeploymentUpdated) | edits an already existing deployment |
| Destroy | [DestroyDeployment](#redsail.bosn.DestroyDeployment) | [DeploymentDestroyed](#redsail.bosn.DeploymentDestroyed) | removes a deployment from the list of configurations |
| Read | [ReadDeployment](#redsail.bosn.ReadDeployment) | [DeploymentRead](#redsail.bosn.DeploymentRead) | reads out a deployment |
| All | [ReadDeployments](#redsail.bosn.ReadDeployments) | [DeploymentsRead](#redsail.bosn.DeploymentsRead) | gets all deployments currently configured and their status |
| Template | [TemplateDeployment](#redsail.bosn.TemplateDeployment) | [DeploymentTemplated](#redsail.bosn.DeploymentTemplated) | get the templated version of this deployment |
| Token | [ReadDeploymentToken](#redsail.bosn.ReadDeploymentToken) | [DeploymentTokenRead](#redsail.bosn.DeploymentTokenRead) | gets the token for this deployment, for use with web calls |
| Run | [ReadRun](#redsail.bosn.ReadRun) | [RunRead](#redsail.bosn.RunRead) | read all the information about a particular run |
| Runs | [ReadRuns](#redsail.bosn.ReadRuns) | [RunsRead](#redsail.bosn.RunsRead) | read summaries of all runs for a particular deployment |
| Approve | [ApproveStep](#redsail.bosn.ApproveStep) | [StepApproved](#redsail.bosn.StepApproved) | approve a step for a run |
| Approvals | [ReadApprovals](#redsail.bosn.ReadApprovals) | [ApprovalsRead](#redsail.bosn.ApprovalsRead) | gets all approvals for the user |

 



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

