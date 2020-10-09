
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface Cluster {
    name: string;
    endpoint: string;
    ready: boolean;
    
}

interface ClusterJSON {
    name: string;
    endpoint: string;
    ready: boolean;
    
}


const JSONToCluster = (m: Cluster | ClusterJSON): Cluster => {
    
    return {
        name: m.name,
        endpoint: m.endpoint,
        ready: m.ready,
        
    };
};

export interface ClusterRequest {
    
}

interface ClusterRequestJSON {
    
}


const ClusterRequestToJSON = (m: ClusterRequest): ClusterRequestJSON => {
    return {
        
    };
};

export interface ClusterResponse {
    clusters: Cluster[];
    
}

interface ClusterResponseJSON {
    clusters: ClusterJSON[];
    
}


const JSONToClusterResponse = (m: ClusterResponse | ClusterResponseJSON): ClusterResponse => {
    
    return {
        clusters: (m.clusters as (Cluster | ClusterJSON)[]).map(JSONToCluster),
        
    };
};

export interface Deployment {
    name: string;
    namespace: string;
    version: string;
    
}

interface DeploymentJSON {
    name: string;
    namespace: string;
    version: string;
    
}


const JSONToDeployment = (m: Deployment | DeploymentJSON): Deployment => {
    
    return {
        name: m.name,
        namespace: m.namespace,
        version: m.version,
        
    };
};

export interface DeploymentRequest {
    cluster: string;
    
}

interface DeploymentRequestJSON {
    cluster: string;
    
}


const DeploymentRequestToJSON = (m: DeploymentRequest): DeploymentRequestJSON => {
    return {
        cluster: m.cluster,
        
    };
};

export interface DeploymentResponse {
    deployments: Deployment[];
    
}

interface DeploymentResponseJSON {
    deployments: DeploymentJSON[];
    
}


const JSONToDeploymentResponse = (m: DeploymentResponse | DeploymentResponseJSON): DeploymentResponse => {
    
    return {
        deployments: (m.deployments as (Deployment | DeploymentJSON)[]).map(JSONToDeployment),
        
    };
};

export interface Kraken {
    clusters: (clusterRequest: ClusterRequest) => Promise<ClusterResponse>;
    
    deployments: (deploymentRequest: DeploymentRequest) => Promise<DeploymentResponse>;
    
}

export class DefaultKraken implements Kraken {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Kraken/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    clusters(clusterRequest: ClusterRequest): Promise<ClusterResponse> {
        const url = this.hostname + this.pathPrefix + "Clusters";
        let body: ClusterRequest | ClusterRequestJSON = clusterRequest;
        if (!this.writeCamelCase) {
            body = ClusterRequestToJSON(clusterRequest);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToClusterResponse);
        });
    }
    
    deployments(deploymentRequest: DeploymentRequest): Promise<DeploymentResponse> {
        const url = this.hostname + this.pathPrefix + "Deployments";
        let body: DeploymentRequest | DeploymentRequestJSON = deploymentRequest;
        if (!this.writeCamelCase) {
            body = DeploymentRequestToJSON(deploymentRequest);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeploymentResponse);
        });
    }
    
}

