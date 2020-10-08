// package: 
// file: kraken/service.proto

import * as kraken_service_pb from "../kraken/service_pb";
import {grpc} from "@improbable-eng/grpc-web";

type KrakenClusters = {
  readonly methodName: string;
  readonly service: typeof Kraken;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof kraken_service_pb.ClusterRequest;
  readonly responseType: typeof kraken_service_pb.ClusterResponse;
};

type KrakenDeployments = {
  readonly methodName: string;
  readonly service: typeof Kraken;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof kraken_service_pb.DeploymentRequest;
  readonly responseType: typeof kraken_service_pb.DeploymentResponse;
};

export class Kraken {
  static readonly serviceName: string;
  static readonly Clusters: KrakenClusters;
  static readonly Deployments: KrakenDeployments;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class KrakenClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  clusters(
    requestMessage: kraken_service_pb.ClusterRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: kraken_service_pb.ClusterResponse|null) => void
  ): UnaryResponse;
  clusters(
    requestMessage: kraken_service_pb.ClusterRequest,
    callback: (error: ServiceError|null, responseMessage: kraken_service_pb.ClusterResponse|null) => void
  ): UnaryResponse;
  deployments(
    requestMessage: kraken_service_pb.DeploymentRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: kraken_service_pb.DeploymentResponse|null) => void
  ): UnaryResponse;
  deployments(
    requestMessage: kraken_service_pb.DeploymentRequest,
    callback: (error: ServiceError|null, responseMessage: kraken_service_pb.DeploymentResponse|null) => void
  ): UnaryResponse;
}

