
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
    name: string;
    chartVersion: string;
    appVersion: string;
    description: string;
    url: string;
    
}

interface ChartVersionJSON {
    name: string;
    chart_version: string;
    app_version: string;
    description: string;
    url: string;
    
}


const JSONToChartVersion = (m: ChartVersion | ChartVersionJSON): ChartVersion => {
    
    return {
        name: m.name,
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

export interface DownloadRequest {
    chartName: string;
    chartVersion: string;
    repoName: string;
    
}

interface DownloadRequestJSON {
    chart_name: string;
    chart_version: string;
    repo_name: string;
    
}


const DownloadRequestToJSON = (m: DownloadRequest): DownloadRequestJSON => {
    return {
        chart_name: m.chartName,
        chart_version: m.chartVersion,
        repo_name: m.repoName,
        
    };
};

export interface File {
    name: string;
    contents: string;
    
}

interface FileJSON {
    name: string;
    contents: string;
    
}


const JSONToFile = (m: File | FileJSON): File => {
    
    return {
        name: m.name,
        contents: m.contents,
        
    };
};

export interface Repo {
    uuid: string;
    name: string;
    endpoint: string;
    ready: boolean;
    
}

interface RepoJSON {
    uuid: string;
    name: string;
    endpoint: string;
    ready: boolean;
    
}


const RepoToJSON = (m: Repo): RepoJSON => {
    return {
        uuid: m.uuid,
        name: m.name,
        endpoint: m.endpoint,
        ready: m.ready,
        
    };
};

const JSONToRepo = (m: Repo | RepoJSON): Repo => {
    
    return {
        uuid: m.uuid,
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

export interface EmptyResponse {
    
}

interface EmptyResponseJSON {
    
}


const JSONToEmptyResponse = (m: EmptyResponse | EmptyResponseJSON): EmptyResponse => {
    
    return {
        
    };
};

export interface Poseidon {
    charts: (repo: Repo) => Promise<ChartsResponse>;
    
    downloadChart: (downloadRequest: DownloadRequest) => Promise<File>;
    
    addRepo: (repo: Repo) => Promise<EmptyResponse>;
    
    deleteRepo: (repo: Repo) => Promise<EmptyResponse>;
    
    editRepo: (repo: Repo) => Promise<EmptyResponse>;
    
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
    
    downloadChart(downloadRequest: DownloadRequest): Promise<File> {
        const url = this.hostname + this.pathPrefix + "DownloadChart";
        let body: DownloadRequest | DownloadRequestJSON = downloadRequest;
        if (!this.writeCamelCase) {
            body = DownloadRequestToJSON(downloadRequest);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToFile);
        });
    }
    
    addRepo(repo: Repo): Promise<EmptyResponse> {
        const url = this.hostname + this.pathPrefix + "AddRepo";
        let body: Repo | RepoJSON = repo;
        if (!this.writeCamelCase) {
            body = RepoToJSON(repo);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToEmptyResponse);
        });
    }
    
    deleteRepo(repo: Repo): Promise<EmptyResponse> {
        const url = this.hostname + this.pathPrefix + "DeleteRepo";
        let body: Repo | RepoJSON = repo;
        if (!this.writeCamelCase) {
            body = RepoToJSON(repo);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToEmptyResponse);
        });
    }
    
    editRepo(repo: Repo): Promise<EmptyResponse> {
        const url = this.hostname + this.pathPrefix + "EditRepo";
        let body: Repo | RepoJSON = repo;
        if (!this.writeCamelCase) {
            body = RepoToJSON(repo);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToEmptyResponse);
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

