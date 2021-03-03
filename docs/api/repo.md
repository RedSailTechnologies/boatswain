# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [repo.proto](#repo.proto)
    - [CreateRepo](#redsail.bosn.CreateRepo)
    - [DestroyRepo](#redsail.bosn.DestroyRepo)
    - [FileRead](#redsail.bosn.FileRead)
    - [ReadFile](#redsail.bosn.ReadFile)
    - [ReadRepo](#redsail.bosn.ReadRepo)
    - [ReadRepos](#redsail.bosn.ReadRepos)
    - [RepoCreated](#redsail.bosn.RepoCreated)
    - [RepoDestroyed](#redsail.bosn.RepoDestroyed)
    - [RepoRead](#redsail.bosn.RepoRead)
    - [RepoUpdated](#redsail.bosn.RepoUpdated)
    - [ReposRead](#redsail.bosn.ReposRead)
    - [UpdateRepo](#redsail.bosn.UpdateRepo)
  
    - [AuthType](#redsail.bosn.AuthType)
    - [RepoType](#redsail.bosn.RepoType)
  
    - [Repo](#redsail.bosn.Repo)
  
- [Scalar Value Types](#scalar-value-types)



<a name="repo.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## repo.proto
Repo is the service managing external repositories, such as helm.


<a name="redsail.bosn.CreateRepo"></a>

### CreateRepo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of the repo |
| endpoint | [string](#string) |  | repo endpoint |
| type | [RepoType](#redsail.bosn.RepoType) |  | type of repo |
| token | [string](#string) |  | an auth token (if required for this repo) |
| username | [string](#string) |  | a username to login with |
| password | [string](#string) |  | a password to login with |
| helm_oci | [bool](#bool) |  | whether or not this is an OCI registry rather than a typical helm registry |






<a name="redsail.bosn.DestroyRepo"></a>

### DestroyRepo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the repo |






<a name="redsail.bosn.FileRead"></a>

### FileRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| file | [bytes](#bytes) |  | the contents of the file read |






<a name="redsail.bosn.ReadFile"></a>

### ReadFile



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| repo_id | [string](#string) |  | unique id of the repo |
| branch | [string](#string) |  | the branch to read the file from |
| file_path | [string](#string) |  | relative path to the file |






<a name="redsail.bosn.ReadRepo"></a>

### ReadRepo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the repo |






<a name="redsail.bosn.ReadRepos"></a>

### ReadRepos







<a name="redsail.bosn.RepoCreated"></a>

### RepoCreated







<a name="redsail.bosn.RepoDestroyed"></a>

### RepoDestroyed







<a name="redsail.bosn.RepoRead"></a>

### RepoRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the repo |
| name | [string](#string) |  | name of the repo |
| endpoint | [string](#string) |  | repo endpoint |
| type | [RepoType](#redsail.bosn.RepoType) |  | type of repo |
| helm_oci | [bool](#bool) |  | whether or not this is an OCI registry rather than a typical helm registry |
| ready | [bool](#bool) |  | repo ready status, based on whether index.yaml can be fetched |






<a name="redsail.bosn.RepoUpdated"></a>

### RepoUpdated







<a name="redsail.bosn.ReposRead"></a>

### ReposRead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| repos | [RepoRead](#redsail.bosn.RepoRead) | repeated | repos read |






<a name="redsail.bosn.UpdateRepo"></a>

### UpdateRepo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  | unique id of the repo |
| name | [string](#string) |  | name of the repo |
| endpoint | [string](#string) |  | repo endpoint |
| type | [RepoType](#redsail.bosn.RepoType) |  | type of repo |
| token | [string](#string) |  | an auth token (if required for this repo) |
| username | [string](#string) |  | a username to login with |
| password | [string](#string) |  | a password to login with |
| helm_oci | [bool](#bool) |  | whether or not this is an OCI registry rather than a typical helm registry |





 


<a name="redsail.bosn.AuthType"></a>

### AuthType


| Name | Number | Description |
| ---- | ------ | ----------- |
| CREDENTIALS | 0 |  |
| TOKEN | 1 |  |



<a name="redsail.bosn.RepoType"></a>

### RepoType


| Name | Number | Description |
| ---- | ------ | ----------- |
| HELM_REPO | 0 |  |
| GIT_REPO | 1 |  |


 

 


<a name="redsail.bosn.Repo"></a>

### Repo


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [CreateRepo](#redsail.bosn.CreateRepo) | [RepoCreated](#redsail.bosn.RepoCreated) | adds a repo to the list of configurations |
| Update | [UpdateRepo](#redsail.bosn.UpdateRepo) | [RepoUpdated](#redsail.bosn.RepoUpdated) | edits an already existing repo |
| Destroy | [DestroyRepo](#redsail.bosn.DestroyRepo) | [RepoDestroyed](#redsail.bosn.RepoDestroyed) | removes a repo from the list of configurations |
| Read | [ReadRepo](#redsail.bosn.ReadRepo) | [RepoRead](#redsail.bosn.RepoRead) | reads out a repo |
| All | [ReadRepos](#redsail.bosn.ReadRepos) | [ReposRead](#redsail.bosn.ReposRead) | gets all repos currently configured and their status |
| File | [ReadFile](#redsail.bosn.ReadFile) | [FileRead](#redsail.bosn.FileRead) | gets the contents of a file from this git repository |

 



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

