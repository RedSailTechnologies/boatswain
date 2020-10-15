
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

export interface Release {
    name: string;
    chart: string;
    namespace: string;
    chartVersion: string;
    appVersion: string;
    clusterName: string;
    status: string;
    
}

interface ReleaseJSON {
    name: string;
    chart: string;
    namespace: string;
    chart_version: string;
    app_version: string;
    cluster_name: string;
    status: string;
    
}


const ReleaseToJSON = (m: Release): ReleaseJSON => {
    return {
        name: m.name,
        chart: m.chart,
        namespace: m.namespace,
        chart_version: m.chartVersion,
        app_version: m.appVersion,
        cluster_name: m.clusterName,
        status: m.status,
        
    };
};

const JSONToRelease = (m: Release | ReleaseJSON): Release => {
    
    return {
        name: m.name,
        chart: m.chart,
        namespace: m.namespace,
        chartVersion: (((m as Release).chartVersion) ? (m as Release).chartVersion : (m as ReleaseJSON).chart_version),
        appVersion: (((m as Release).appVersion) ? (m as Release).appVersion : (m as ReleaseJSON).app_version),
        clusterName: (((m as Release).clusterName) ? (m as Release).clusterName : (m as ReleaseJSON).cluster_name),
        status: m.status,
        
    };
};

export interface Releases {
    name: string;
    chart: string;
    releases: Release[];
    
}

interface ReleasesJSON {
    name: string;
    chart: string;
    releases: ReleaseJSON[];
    
}


const JSONToReleases = (m: Releases | ReleasesJSON): Releases => {
    
    return {
        name: m.name,
        chart: m.chart,
        releases: (m.releases as (Release | ReleaseJSON)[]).map(JSONToRelease),
        
    };
};

export interface ReleaseRequest {
    clusters: Cluster[];
    
}

interface ReleaseRequestJSON {
    clusters: ClusterJSON[];
    
}


const ReleaseRequestToJSON = (m: ReleaseRequest): ReleaseRequestJSON => {
    return {
        clusters: m.clusters.map(ClusterToJSON),
        
    };
};

export interface ReleaseResponse {
    releaseLists: Releases[];
    
}

interface ReleaseResponseJSON {
    release_lists: ReleasesJSON[];
    
}


const JSONToReleaseResponse = (m: ReleaseResponse | ReleaseResponseJSON): ReleaseResponse => {
    
    return {
        releaseLists: ((((m as ReleaseResponse).releaseLists) ? (m as ReleaseResponse).releaseLists : (m as ReleaseResponseJSON).release_lists) as (Releases | ReleasesJSON)[]).map(JSONToReleases),
        
    };
};

export interface UpgradeReleaseRequest {
    name: string;
    chart: string;
    namespace: string;
    chartVersion: string;
    appVersion: string;
    clusterName: string;
    repoName: string;
    values: string;
    
}

interface UpgradeReleaseRequestJSON {
    name: string;
    chart: string;
    namespace: string;
    chart_version: string;
    app_version: string;
    cluster_name: string;
    repo_name: string;
    values: string;
    
}


const UpgradeReleaseRequestToJSON = (m: UpgradeReleaseRequest): UpgradeReleaseRequestJSON => {
    return {
        name: m.name,
        chart: m.chart,
        namespace: m.namespace,
        chart_version: m.chartVersion,
        app_version: m.appVersion,
        cluster_name: m.clusterName,
        repo_name: m.repoName,
        values: m.values,
        
    };
};

export interface Kraken {
    clusters: (clustersRequest: ClustersRequest) => Promise<ClustersResponse>;
    
    clusterStatus: (cluster: Cluster) => Promise<Cluster>;
    
    releases: (releaseRequest: ReleaseRequest) => Promise<ReleaseResponse>;
    
    releaseStatus: (release: Release) => Promise<Release>;
    
    upgradeRelease: (upgradeReleaseRequest: UpgradeReleaseRequest) => Promise<Release>;
    
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
    
    releases(releaseRequest: ReleaseRequest): Promise<ReleaseResponse> {
        const url = this.hostname + this.pathPrefix + "Releases";
        let body: ReleaseRequest | ReleaseRequestJSON = releaseRequest;
        if (!this.writeCamelCase) {
            body = ReleaseRequestToJSON(releaseRequest);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToReleaseResponse);
        });
    }
    
    releaseStatus(release: Release): Promise<Release> {
        const url = this.hostname + this.pathPrefix + "ReleaseStatus";
        let body: Release | ReleaseJSON = release;
        if (!this.writeCamelCase) {
            body = ReleaseToJSON(release);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToRelease);
        });
    }
    
    upgradeRelease(upgradeReleaseRequest: UpgradeReleaseRequest): Promise<Release> {
        const url = this.hostname + this.pathPrefix + "UpgradeRelease";
        let body: UpgradeReleaseRequest | UpgradeReleaseRequestJSON = upgradeReleaseRequest;
        if (!this.writeCamelCase) {
            body = UpgradeReleaseRequestToJSON(upgradeReleaseRequest);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToRelease);
        });
    }
    
}

