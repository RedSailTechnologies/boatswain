
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface ApplicationRead {
    name: string;
    project: string;
    clusters: ApplicationCluster[];
    
}

interface ApplicationReadJSON {
    name: string;
    project: string;
    clusters: ApplicationClusterJSON[];
    
}


const JSONToApplicationRead = (m: ApplicationRead | ApplicationReadJSON): ApplicationRead => {
    
    return {
        name: m.name,
        project: m.project,
        clusters: (m.clusters as (ApplicationCluster | ApplicationClusterJSON)[]).map(JSONToApplicationCluster),
        
    };
};

export interface ApplicationCluster {
    clusterName: string;
    version: string;
    namespace: string;
    ready: boolean;
    
}

interface ApplicationClusterJSON {
    cluster_name: string;
    version: string;
    namespace: string;
    ready: boolean;
    
}


const JSONToApplicationCluster = (m: ApplicationCluster | ApplicationClusterJSON): ApplicationCluster => {
    
    return {
        clusterName: (((m as ApplicationCluster).clusterName) ? (m as ApplicationCluster).clusterName : (m as ApplicationClusterJSON).cluster_name),
        version: m.version,
        namespace: m.namespace,
        ready: m.ready,
        
    };
};

export interface ReadApplications {
    
}

interface ReadApplicationsJSON {
    
}


const ReadApplicationsToJSON = (m: ReadApplications): ReadApplicationsJSON => {
    return {
        
    };
};

export interface ApplicationsRead {
    applications: ApplicationRead[];
    
}

interface ApplicationsReadJSON {
    applications: ApplicationReadJSON[];
    
}


const JSONToApplicationsRead = (m: ApplicationsRead | ApplicationsReadJSON): ApplicationsRead => {
    
    return {
        applications: (m.applications as (ApplicationRead | ApplicationReadJSON)[]).map(JSONToApplicationRead),
        
    };
};

export interface Application {
    all: (readApplications: ReadApplications) => Promise<ApplicationsRead>;
    
}

export class DefaultApplication implements Application {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Application/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    all(readApplications: ReadApplications): Promise<ApplicationsRead> {
        const url = this.hostname + this.pathPrefix + "All";
        let body: ReadApplications | ReadApplicationsJSON = readApplications;
        if (!this.writeCamelCase) {
            body = ReadApplicationsToJSON(readApplications);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToApplicationsRead);
        });
    }
    
}

