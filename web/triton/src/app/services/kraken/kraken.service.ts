import { Injectable } from '@angular/core';
import { grpc } from '@improbable-eng/grpc-web';
import { ClusterRequest, ClusterResponse, DeploymentRequest, DeploymentResponse } from './service_pb';
import { Kraken } from './service_pb_service'

@Injectable({
  providedIn: 'root'
})
export class KrakenService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = `${location.protocol}//${location.host}`;
  }

  public getClusters() : Promise<ClusterResponse.AsObject> {
    const request = new ClusterRequest();
    return new Promise((resolve, reject) => {
      grpc.unary<ClusterRequest, ClusterResponse, typeof Kraken.Clusters>(Kraken.Clusters, {
        request: request,
        host: this.baseUrl,
        onEnd: response => {
          const { status, statusMessage, headers, message, trailers } = response;
          if (status == grpc.Code.OK && message) {
            resolve(message.toObject());
          } else {
            reject(statusMessage)
          }
        }
      })
    });
  }

  public getDeployments(cluster: string) : Promise<DeploymentResponse.AsObject> {
    const request = new DeploymentRequest();
    request.setCluster(cluster);
    return new Promise((resolve, reject) => {
      grpc.unary<DeploymentRequest, DeploymentResponse, typeof Kraken.Deployments>(Kraken.Deployments, {
        request: request,
        host: this.baseUrl,
        onEnd: response => {
          const { status, statusMessage, headers, message, trailers } = response;
          if (status == grpc.Code.OK && message) {
            resolve(message.toObject())
          } else {
            reject(statusMessage)
          }
        }
      })
    })
  }
}
