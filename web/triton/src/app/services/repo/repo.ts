
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface CreateRepo {
    name: string;
    endpoint: string;
    type: string;
    token: string;
    username: string;
    password: string;
    helmOci: boolean;
    
}

interface CreateRepoJSON {
    name: string;
    endpoint: string;
    type: string;
    token: string;
    username: string;
    password: string;
    helm_oci: boolean;
    
}


const CreateRepoToJSON = (m: CreateRepo): CreateRepoJSON => {
    return {
        name: m.name,
        endpoint: m.endpoint,
        type: m.type,
        token: m.token,
        username: m.username,
        password: m.password,
        helm_oci: m.helmOci,
        
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
    type: string;
    token: string;
    username: string;
    password: string;
    helmOci: boolean;
    
}

interface UpdateRepoJSON {
    uuid: string;
    name: string;
    endpoint: string;
    type: string;
    token: string;
    username: string;
    password: string;
    helm_oci: boolean;
    
}


const UpdateRepoToJSON = (m: UpdateRepo): UpdateRepoJSON => {
    return {
        uuid: m.uuid,
        name: m.name,
        endpoint: m.endpoint,
        type: m.type,
        token: m.token,
        username: m.username,
        password: m.password,
        helm_oci: m.helmOci,
        
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
    type: string;
    helmOci: boolean;
    ready: boolean;
    
}

interface RepoReadJSON {
    uuid: string;
    name: string;
    endpoint: string;
    type: string;
    helm_oci: boolean;
    ready: boolean;
    
}


const JSONToRepoRead = (m: RepoRead | RepoReadJSON): RepoRead => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        endpoint: m.endpoint,
        type: m.type,
        helmOci: (((m as RepoRead).helmOci) ? (m as RepoRead).helmOci : (m as RepoReadJSON).helm_oci),
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

export interface ReadFile {
    repoId: string;
    branch: string;
    filePath: string;
    
}

interface ReadFileJSON {
    repo_id: string;
    branch: string;
    file_path: string;
    
}


const ReadFileToJSON = (m: ReadFile): ReadFileJSON => {
    return {
        repo_id: m.repoId,
        branch: m.branch,
        file_path: m.filePath,
        
    };
};

export interface FileRead {
    file: string;
    
}

interface FileReadJSON {
    file: string;
    
}


const JSONToFileRead = (m: FileRead | FileReadJSON): FileRead => {
    
    return {
        file: m.file,
        
    };
};

export interface Repo {
    create: (createRepo: CreateRepo) => Promise<RepoCreated>;
    
    update: (updateRepo: UpdateRepo) => Promise<RepoUpdated>;
    
    destroy: (destroyRepo: DestroyRepo) => Promise<RepoDestroyed>;
    
    read: (readRepo: ReadRepo) => Promise<RepoRead>;
    
    all: (readRepos: ReadRepos) => Promise<ReposRead>;
    
    file: (readFile: ReadFile) => Promise<FileRead>;
    
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
    
    file(readFile: ReadFile): Promise<FileRead> {
        const url = this.hostname + this.pathPrefix + "File";
        let body: ReadFile | ReadFileJSON = readFile;
        if (!this.writeCamelCase) {
            body = ReadFileToJSON(readFile);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToFileRead);
        });
    }
    
}

