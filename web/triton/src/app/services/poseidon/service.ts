
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


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

