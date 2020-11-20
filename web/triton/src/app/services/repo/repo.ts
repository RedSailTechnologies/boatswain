
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface CreateRepo {
    name: string;
    endpoint: string;
    
}

interface CreateRepoJSON {
    name: string;
    endpoint: string;
    
}


const CreateRepoToJSON = (m: CreateRepo): CreateRepoJSON => {
    return {
        name: m.name,
        endpoint: m.endpoint,
        
    };
};

export interface RepoCreated {
    
}

interface RepoCreatedJSON {
    
}


const JSONToRepoCreated = (m: RepoCreated | RepoCreatedJSON): RepoCreated => {
    
    return {
        
    };
};

export interface UpdateRepo {
    uuid: string;
    name: string;
    endpoint: string;
    
}

interface UpdateRepoJSON {
    uuid: string;
    name: string;
    endpoint: string;
    
}


const UpdateRepoToJSON = (m: UpdateRepo): UpdateRepoJSON => {
    return {
        uuid: m.uuid,
        name: m.name,
        endpoint: m.endpoint,
        
    };
};

export interface RepoUpdated {
    
}

interface RepoUpdatedJSON {
    
}


const JSONToRepoUpdated = (m: RepoUpdated | RepoUpdatedJSON): RepoUpdated => {
    
    return {
        
    };
};

export interface DestroyRepo {
    uuid: string;
    
}

interface DestroyRepoJSON {
    uuid: string;
    
}


const DestroyRepoToJSON = (m: DestroyRepo): DestroyRepoJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface RepoDestroyed {
    
}

interface RepoDestroyedJSON {
    
}


const JSONToRepoDestroyed = (m: RepoDestroyed | RepoDestroyedJSON): RepoDestroyed => {
    
    return {
        
    };
};

export interface ReadRepo {
    uuid: string;
    
}

interface ReadRepoJSON {
    uuid: string;
    
}


const ReadRepoToJSON = (m: ReadRepo): ReadRepoJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface RepoRead {
    uuid: string;
    name: string;
    endpoint: string;
    ready: boolean;
    
}

interface RepoReadJSON {
    uuid: string;
    name: string;
    endpoint: string;
    ready: boolean;
    
}


const JSONToRepoRead = (m: RepoRead | RepoReadJSON): RepoRead => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        endpoint: m.endpoint,
        ready: m.ready,
        
    };
};

export interface ReadRepos {
    
}

interface ReadReposJSON {
    
}


const ReadReposToJSON = (m: ReadRepos): ReadReposJSON => {
    return {
        
    };
};

export interface ReposRead {
    repos: RepoRead[];
    
}

interface ReposReadJSON {
    repos: RepoReadJSON[];
    
}


const JSONToReposRead = (m: ReposRead | ReposReadJSON): ReposRead => {
    
    return {
        repos: (m.repos as (RepoRead | RepoReadJSON)[]).map(JSONToRepoRead),
        
    };
};

export interface ChartRead {
    name: string;
    versions: VersionRead[];
    
}

interface ChartReadJSON {
    name: string;
    versions: VersionReadJSON[];
    
}


const JSONToChartRead = (m: ChartRead | ChartReadJSON): ChartRead => {
    
    return {
        name: m.name,
        versions: (m.versions as (VersionRead | VersionReadJSON)[]).map(JSONToVersionRead),
        
    };
};

export interface VersionRead {
    name: string;
    chartVersion: string;
    appVersion: string;
    description: string;
    url: string;
    
}

interface VersionReadJSON {
    name: string;
    chart_version: string;
    app_version: string;
    description: string;
    url: string;
    
}


const JSONToVersionRead = (m: VersionRead | VersionReadJSON): VersionRead => {
    
    return {
        name: m.name,
        chartVersion: (((m as VersionRead).chartVersion) ? (m as VersionRead).chartVersion : (m as VersionReadJSON).chart_version),
        appVersion: (((m as VersionRead).appVersion) ? (m as VersionRead).appVersion : (m as VersionReadJSON).app_version),
        description: m.description,
        url: m.url,
        
    };
};

export interface ChartsRead {
    charts: ChartRead[];
    
}

interface ChartsReadJSON {
    charts: ChartReadJSON[];
    
}


const JSONToChartsRead = (m: ChartsRead | ChartsReadJSON): ChartsRead => {
    
    return {
        charts: (m.charts as (ChartRead | ChartReadJSON)[]).map(JSONToChartRead),
        
    };
};

export interface Repo {
    create: (createRepo: CreateRepo) => Promise<RepoCreated>;
    
    update: (updateRepo: UpdateRepo) => Promise<RepoUpdated>;
    
    destroy: (destroyRepo: DestroyRepo) => Promise<RepoDestroyed>;
    
    read: (readRepo: ReadRepo) => Promise<RepoRead>;
    
    all: (readRepos: ReadRepos) => Promise<ReposRead>;
    
    charts: (readRepo: ReadRepo) => Promise<ChartsRead>;
    
}

export class DefaultRepo implements Repo {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Repo/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    create(createRepo: CreateRepo): Promise<RepoCreated> {
        const url = this.hostname + this.pathPrefix + "Create";
        let body: CreateRepo | CreateRepoJSON = createRepo;
        if (!this.writeCamelCase) {
            body = CreateRepoToJSON(createRepo);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToRepoCreated);
        });
    }
    
    update(updateRepo: UpdateRepo): Promise<RepoUpdated> {
        const url = this.hostname + this.pathPrefix + "Update";
        let body: UpdateRepo | UpdateRepoJSON = updateRepo;
        if (!this.writeCamelCase) {
            body = UpdateRepoToJSON(updateRepo);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToRepoUpdated);
        });
    }
    
    destroy(destroyRepo: DestroyRepo): Promise<RepoDestroyed> {
        const url = this.hostname + this.pathPrefix + "Destroy";
        let body: DestroyRepo | DestroyRepoJSON = destroyRepo;
        if (!this.writeCamelCase) {
            body = DestroyRepoToJSON(destroyRepo);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToRepoDestroyed);
        });
    }
    
    read(readRepo: ReadRepo): Promise<RepoRead> {
        const url = this.hostname + this.pathPrefix + "Read";
        let body: ReadRepo | ReadRepoJSON = readRepo;
        if (!this.writeCamelCase) {
            body = ReadRepoToJSON(readRepo);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToRepoRead);
        });
    }
    
    all(readRepos: ReadRepos): Promise<ReposRead> {
        const url = this.hostname + this.pathPrefix + "All";
        let body: ReadRepos | ReadReposJSON = readRepos;
        if (!this.writeCamelCase) {
            body = ReadReposToJSON(readRepos);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToReposRead);
        });
    }
    
    charts(readRepo: ReadRepo): Promise<ChartsRead> {
        const url = this.hostname + this.pathPrefix + "Charts";
        let body: ReadRepo | ReadRepoJSON = readRepo;
        if (!this.writeCamelCase) {
            body = ReadRepoToJSON(readRepo);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToChartsRead);
        });
    }
    
}

