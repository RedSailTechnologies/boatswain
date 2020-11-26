
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface CreateCluster {
    name: string;
    endpoint: string;
    token: string;
    cert: string;
    
}

interface CreateClusterJSON {
    name: string;
    endpoint: string;
    token: string;
    cert: string;
    
}


const CreateClusterToJSON = (m: CreateCluster): CreateClusterJSON => {
    return {
        name: m.name,
        endpoint: m.endpoint,
        token: m.token,
        cert: m.cert,
        
    };
};

export interface ClusterCreated {
    
}

interface ClusterCreatedJSON {
    
}


const JSONToClusterCreated = (m: ClusterCreated | ClusterCreatedJSON): ClusterCreated => {
    
    return {
        
    };
};

export interface UpdateCluster {
    uuid: string;
    name: string;
    endpoint: string;
    token: string;
    cert: string;
    
}

interface UpdateClusterJSON {
    uuid: string;
    name: string;
    endpoint: string;
    token: string;
    cert: string;
    
}


const UpdateClusterToJSON = (m: UpdateCluster): UpdateClusterJSON => {
    return {
        uuid: m.uuid,
        name: m.name,
        endpoint: m.endpoint,
        token: m.token,
        cert: m.cert,
        
    };
};

export interface ClusterUpdated {
    
}

interface ClusterUpdatedJSON {
    
}


const JSONToClusterUpdated = (m: ClusterUpdated | ClusterUpdatedJSON): ClusterUpdated => {
    
    return {
        
    };
};

export interface DestroyCluster {
    uuid: string;
    
}

interface DestroyClusterJSON {
    uuid: string;
    
}


const DestroyClusterToJSON = (m: DestroyCluster): DestroyClusterJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface ClusterDestroyed {
    
}

interface ClusterDestroyedJSON {
    
}


const JSONToClusterDestroyed = (m: ClusterDestroyed | ClusterDestroyedJSON): ClusterDestroyed => {
    
    return {
        
    };
};

export interface ReadCluster {
    uuid: string;
    
}

interface ReadClusterJSON {
    uuid: string;
    
}


const ReadClusterToJSON = (m: ReadCluster): ReadClusterJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface ClusterRead {
    uuid: string;
    name: string;
    endpoint: string;
    token: string;
    cert: string;
    ready: boolean;
    
}

interface ClusterReadJSON {
    uuid: string;
    name: string;
    endpoint: string;
    token: string;
    cert: string;
    ready: boolean;
    
}


const JSONToClusterRead = (m: ClusterRead | ClusterReadJSON): ClusterRead => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        endpoint: m.endpoint,
        token: m.token,
        cert: m.cert,
        ready: m.ready,
        
    };
};

export interface ReadClusters {
    
}

interface ReadClustersJSON {
    
}


const ReadClustersToJSON = (m: ReadClusters): ReadClustersJSON => {
    return {
        
    };
};

export interface ClustersRead {
    clusters: ClusterRead[];
    
}

interface ClustersReadJSON {
    clusters: ClusterReadJSON[];
    
}


const JSONToClustersRead = (m: ClustersRead | ClustersReadJSON): ClustersRead => {
    
    return {
        clusters: (m.clusters as (ClusterRead | ClusterReadJSON)[]).map(JSONToClusterRead),
        
    };
};

export interface Cluster {
    create: (createCluster: CreateCluster) => Promise<ClusterCreated>;
    
    update: (updateCluster: UpdateCluster) => Promise<ClusterUpdated>;
    
    destroy: (destroyCluster: DestroyCluster) => Promise<ClusterDestroyed>;
    
    read: (readCluster: ReadCluster) => Promise<ClusterRead>;
    
    all: (readClusters: ReadClusters) => Promise<ClustersRead>;
    
}

export class DefaultCluster implements Cluster {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Cluster/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    create(createCluster: CreateCluster): Promise<ClusterCreated> {
        const url = this.hostname + this.pathPrefix + "Create";
        let body: CreateCluster | CreateClusterJSON = createCluster;
        if (!this.writeCamelCase) {
            body = CreateClusterToJSON(createCluster);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToClusterCreated);
        });
    }
    
    update(updateCluster: UpdateCluster): Promise<ClusterUpdated> {
        const url = this.hostname + this.pathPrefix + "Update";
        let body: UpdateCluster | UpdateClusterJSON = updateCluster;
        if (!this.writeCamelCase) {
            body = UpdateClusterToJSON(updateCluster);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToClusterUpdated);
        });
    }
    
    destroy(destroyCluster: DestroyCluster): Promise<ClusterDestroyed> {
        const url = this.hostname + this.pathPrefix + "Destroy";
        let body: DestroyCluster | DestroyClusterJSON = destroyCluster;
        if (!this.writeCamelCase) {
            body = DestroyClusterToJSON(destroyCluster);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToClusterDestroyed);
        });
    }
    
    read(readCluster: ReadCluster): Promise<ClusterRead> {
        const url = this.hostname + this.pathPrefix + "Read";
        let body: ReadCluster | ReadClusterJSON = readCluster;
        if (!this.writeCamelCase) {
            body = ReadClusterToJSON(readCluster);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToClusterRead);
        });
    }
    
    all(readClusters: ReadClusters): Promise<ClustersRead> {
        const url = this.hostname + this.pathPrefix + "All";
        let body: ReadClusters | ReadClustersJSON = readClusters;
        if (!this.writeCamelCase) {
            body = ReadClustersToJSON(readClusters);
        }
        console.log("calling fetch")
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToClustersRead);
        });
    }
    
}

