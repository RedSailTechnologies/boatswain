// package: 
// file: kraken/service.proto

import * as jspb from "google-protobuf";

export class Cluster extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Cluster.AsObject;
  static toObject(includeInstance: boolean, msg: Cluster): Cluster.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Cluster, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Cluster;
  static deserializeBinaryFromReader(message: Cluster, reader: jspb.BinaryReader): Cluster;
}

export namespace Cluster {
  export type AsObject = {
    name: string,
  }
}

export class ClusterRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClusterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ClusterRequest): ClusterRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ClusterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClusterRequest;
  static deserializeBinaryFromReader(message: ClusterRequest, reader: jspb.BinaryReader): ClusterRequest;
}

export namespace ClusterRequest {
  export type AsObject = {
  }
}

export class ClusterResponse extends jspb.Message {
  clearClustersList(): void;
  getClustersList(): Array<Cluster>;
  setClustersList(value: Array<Cluster>): void;
  addClusters(value?: Cluster, index?: number): Cluster;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClusterResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ClusterResponse): ClusterResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ClusterResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClusterResponse;
  static deserializeBinaryFromReader(message: ClusterResponse, reader: jspb.BinaryReader): ClusterResponse;
}

export namespace ClusterResponse {
  export type AsObject = {
    clustersList: Array<Cluster.AsObject>,
  }
}

export class Deployment extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getNamespace(): string;
  setNamespace(value: string): void;

  getVersion(): string;
  setVersion(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Deployment.AsObject;
  static toObject(includeInstance: boolean, msg: Deployment): Deployment.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Deployment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Deployment;
  static deserializeBinaryFromReader(message: Deployment, reader: jspb.BinaryReader): Deployment;
}

export namespace Deployment {
  export type AsObject = {
    name: string,
    namespace: string,
    version: string,
  }
}

export class DeploymentRequest extends jspb.Message {
  getCluster(): string;
  setCluster(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentRequest): DeploymentRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeploymentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentRequest;
  static deserializeBinaryFromReader(message: DeploymentRequest, reader: jspb.BinaryReader): DeploymentRequest;
}

export namespace DeploymentRequest {
  export type AsObject = {
    cluster: string,
  }
}

export class DeploymentResponse extends jspb.Message {
  clearDeploymentsList(): void;
  getDeploymentsList(): Array<Deployment>;
  setDeploymentsList(value: Array<Deployment>): void;
  addDeployments(value?: Deployment, index?: number): Deployment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentResponse): DeploymentResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DeploymentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentResponse;
  static deserializeBinaryFromReader(message: DeploymentResponse, reader: jspb.BinaryReader): DeploymentResponse;
}

export namespace DeploymentResponse {
  export type AsObject = {
    deploymentsList: Array<Deployment.AsObject>,
  }
}

