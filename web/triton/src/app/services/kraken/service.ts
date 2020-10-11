
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


const ClusterToJSON = (m: Cluster): ClusterJSON => {
    return {
        name: m.name,
        endpoint: m.endpoint,
        ready: m.ready,
        
    };
};

const JSONToCluster = (m: Cluster | ClusterJSON): Cluster => {
    
    return {
        name: m.name,
        endpoint: m.endpoint,
        ready: m.ready,
        
    };
};

export interface ClustersRequest {
    
}

interface ClustersRequestJSON {
    
}


const ClustersRequestToJSON = (m: ClustersRequest): ClustersRequestJSON => {
    return {
        
    };
};

export interface ClustersResponse {
    clusters: Cluster[];
    
}

interface ClustersResponseJSON {
    clusters: ClusterJSON[];
    
}


const JSONToClustersResponse = (m: ClustersResponse | ClustersResponseJSON): ClustersResponse => {
    
    return {
        clusters: (m.clusters as (Cluster | ClusterJSON)[]).map(JSONToCluster),
        
    };
};

export interface Deployment {
    name: string;
    namespace: string;
    ready: boolean;
    version: string;
    cluster: Cluster;
    
}

interface DeploymentJSON {
    name: string;
    namespace: string;
    ready: boolean;
    version: string;
    cluster: ClusterJSON;
    
}


const DeploymentToJSON = (m: Deployment): DeploymentJSON => {
    return {
        name: m.name,
        namespace: m.namespace,
        ready: m.ready,
        version: m.version,
        cluster: ClusterToJSON(m.cluster),
        
    };
};

const JSONToDeployment = (m: Deployment | DeploymentJSON): Deployment => {
    
    return {
        name: m.name,
        namespace: m.namespace,
        ready: m.ready,
        version: m.version,
        cluster: JSONToCluster(m.cluster),
        
    };
};

export interface DeploymentsRequest {
    cluster: Cluster;
    
}

interface DeploymentsRequestJSON {
    cluster: ClusterJSON;
    
}


const DeploymentsRequestToJSON = (m: DeploymentsRequest): DeploymentsRequestJSON => {
    return {
        cluster: ClusterToJSON(m.cluster),
        
    };
};

export interface DeploymentsResponse {
    deployments: Deployment[];
    
}

interface DeploymentsResponseJSON {
    deployments: DeploymentJSON[];
    
}


const JSONToDeploymentsResponse = (m: DeploymentsResponse | DeploymentsResponseJSON): DeploymentsResponse => {
    
    return {
        deployments: (m.deployments as (Deployment | DeploymentJSON)[]).map(JSONToDeployment),
        
    };
};

export interface Kraken {
    clusters: (clustersRequest: ClustersRequest) => Promise<ClustersResponse>;
    
    clusterStatus: (cluster: Cluster) => Promise<Cluster>;
    
    deployments: (deploymentsRequest: DeploymentsRequest) => Promise<DeploymentsResponse>;
    
    deploymentStatus: (deployment: Deployment) => Promise<Deployment>;
    
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
    clusters(clustersRequest: ClustersRequest): Promise<ClustersResponse> {
        const url = this.hostname + this.pathPrefix + "Clusters";
        let body: ClustersRequest | ClustersRequestJSON = clustersRequest;
        if (!this.writeCamelCase) {
            body = ClustersRequestToJSON(clustersRequest);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToClustersResponse);
        });
    }
    
    clusterStatus(cluster: Cluster): Promise<Cluster> {
        const url = this.hostname + this.pathPrefix + "ClusterStatus";
        let body: Cluster | ClusterJSON = cluster;
        if (!this.writeCamelCase) {
            body = ClusterToJSON(cluster);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToCluster);
        });
    }
    
    deployments(deploymentsRequest: DeploymentsRequest): Promise<DeploymentsResponse> {
        const url = this.hostname + this.pathPrefix + "Deployments";
        let body: DeploymentsRequest | DeploymentsRequestJSON = deploymentsRequest;
        if (!this.writeCamelCase) {
            body = DeploymentsRequestToJSON(deploymentsRequest);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeploymentsResponse);
        });
    }
    
    deploymentStatus(deployment: Deployment): Promise<Deployment> {
        const url = this.hostname + this.pathPrefix + "DeploymentStatus";
        let body: Deployment | DeploymentJSON = deployment;
        if (!this.writeCamelCase) {
            body = DeploymentToJSON(deployment);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeployment);
        });
    }
    
}

