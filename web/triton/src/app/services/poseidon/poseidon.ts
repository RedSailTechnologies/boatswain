
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface Chart {
    name: string;
    versions: ChartVersion[];
    
}

interface ChartJSON {
    name: string;
    versions: ChartVersionJSON[];
    
}


const JSONToChart = (m: Chart | ChartJSON): Chart => {
    
    return {
        name: m.name,
        versions: (m.versions as (ChartVersion | ChartVersionJSON)[]).map(JSONToChartVersion),
        
    };
};

export interface ChartVersion {
    chartVersion: string;
    appVersion: string;
    description: string;
    url: string;
    
}

interface ChartVersionJSON {
    chart_version: string;
    app_version: string;
    description: string;
    url: string;
    
}


const JSONToChartVersion = (m: ChartVersion | ChartVersionJSON): ChartVersion => {
    
    return {
        chartVersion: (((m as ChartVersion).chartVersion) ? (m as ChartVersion).chartVersion : (m as ChartVersionJSON).chart_version),
        appVersion: (((m as ChartVersion).appVersion) ? (m as ChartVersion).appVersion : (m as ChartVersionJSON).app_version),
        description: m.description,
        url: m.url,
        
    };
};

export interface ChartsResponse {
    charts: Chart[];
    
}

interface ChartsResponseJSON {
    charts: ChartJSON[];
    
}


const JSONToChartsResponse = (m: ChartsResponse | ChartsResponseJSON): ChartsResponse => {
    
    return {
        charts: (m.charts as (Chart | ChartJSON)[]).map(JSONToChart),
        
    };
};

export interface Repo {
    name: string;
    endpoint: string;
    ready: boolean;
    
}

interface RepoJSON {
    name: string;
    endpoint: string;
    ready: boolean;
    
}


const RepoToJSON = (m: Repo): RepoJSON => {
    return {
        name: m.name,
        endpoint: m.endpoint,
        ready: m.ready,
        
    };
};

const JSONToRepo = (m: Repo | RepoJSON): Repo => {
    
    return {
        name: m.name,
        endpoint: m.endpoint,
        ready: m.ready,
        
    };
};

export interface ReposRequest {
    
}

interface ReposRequestJSON {
    
}


const ReposRequestToJSON = (m: ReposRequest): ReposRequestJSON => {
    return {
        
    };
};

export interface ReposResponse {
    repos: Repo[];
    
}

interface ReposResponseJSON {
    repos: RepoJSON[];
    
}


const JSONToReposResponse = (m: ReposResponse | ReposResponseJSON): ReposResponse => {
    
    return {
        repos: (m.repos as (Repo | RepoJSON)[]).map(JSONToRepo),
        
    };
};

export interface Poseidon {
    charts: (repo: Repo) => Promise<ChartsResponse>;
    
    repos: (reposRequest: ReposRequest) => Promise<ReposResponse>;
    
}

export class DefaultPoseidon implements Poseidon {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Poseidon/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    charts(repo: Repo): Promise<ChartsResponse> {
        const url = this.hostname + this.pathPrefix + "Charts";
        let body: Repo | RepoJSON = repo;
        if (!this.writeCamelCase) {
            body = RepoToJSON(repo);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToChartsResponse);
        });
    }
    
    repos(reposRequest: ReposRequest): Promise<ReposResponse> {
        const url = this.hostname + this.pathPrefix + "Repos";
        let body: ReposRequest | ReposRequestJSON = reposRequest;
        if (!this.writeCamelCase) {
            body = ReposRequestToJSON(reposRequest);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToReposResponse);
        });
    }
    
}

