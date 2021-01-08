
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface CreateDeployment {
    name: string;
    repoId: string;
    branch: string;
    filePath: string;
    
}

interface CreateDeploymentJSON {
    name: string;
    repo_id: string;
    branch: string;
    file_path: string;
    
}


const CreateDeploymentToJSON = (m: CreateDeployment): CreateDeploymentJSON => {
    return {
        name: m.name,
        repo_id: m.repoId,
        branch: m.branch,
        file_path: m.filePath,
        
    };
};

export interface DeploymentCreated {
    
}

interface DeploymentCreatedJSON {
    
}


const JSONToDeploymentCreated = (m: DeploymentCreated | DeploymentCreatedJSON): DeploymentCreated => {
    
    return {
        
    };
};

export interface UpdateDeployment {
    uuid: string;
    name: string;
    repoId: string;
    branch: string;
    filePath: string;
    
}

interface UpdateDeploymentJSON {
    uuid: string;
    name: string;
    repo_id: string;
    branch: string;
    file_path: string;
    
}


const UpdateDeploymentToJSON = (m: UpdateDeployment): UpdateDeploymentJSON => {
    return {
        uuid: m.uuid,
        name: m.name,
        repo_id: m.repoId,
        branch: m.branch,
        file_path: m.filePath,
        
    };
};

export interface DeploymentUpdated {
    
}

interface DeploymentUpdatedJSON {
    
}


const JSONToDeploymentUpdated = (m: DeploymentUpdated | DeploymentUpdatedJSON): DeploymentUpdated => {
    
    return {
        
    };
};

export interface DestroyDeployment {
    uuid: string;
    
}

interface DestroyDeploymentJSON {
    uuid: string;
    
}


const DestroyDeploymentToJSON = (m: DestroyDeployment): DestroyDeploymentJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface DeploymentDestroyed {
    
}

interface DeploymentDestroyedJSON {
    
}


const JSONToDeploymentDestroyed = (m: DeploymentDestroyed | DeploymentDestroyedJSON): DeploymentDestroyed => {
    
    return {
        
    };
};

export interface ReadDeployment {
    uuid: string;
    
}

interface ReadDeploymentJSON {
    uuid: string;
    
}


const ReadDeploymentToJSON = (m: ReadDeployment): ReadDeploymentJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface DeploymentRead {
    uuid: string;
    name: string;
    repoId: string;
    repoName: string;
    branch: string;
    filePath: string;
    yaml: string;
    
}

interface DeploymentReadJSON {
    uuid: string;
    name: string;
    repo_id: string;
    repo_name: string;
    branch: string;
    file_path: string;
    yaml: string;
    
}


const JSONToDeploymentRead = (m: DeploymentRead | DeploymentReadJSON): DeploymentRead => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        repoId: (((m as DeploymentRead).repoId) ? (m as DeploymentRead).repoId : (m as DeploymentReadJSON).repo_id),
        repoName: (((m as DeploymentRead).repoName) ? (m as DeploymentRead).repoName : (m as DeploymentReadJSON).repo_name),
        branch: m.branch,
        filePath: (((m as DeploymentRead).filePath) ? (m as DeploymentRead).filePath : (m as DeploymentReadJSON).file_path),
        yaml: m.yaml,
        
    };
};

export interface DeploymentReadSummary {
    uuid: string;
    name: string;
    repoId: string;
    repoName: string;
    branch: string;
    filePath: string;
    
}

interface DeploymentReadSummaryJSON {
    uuid: string;
    name: string;
    repo_id: string;
    repo_name: string;
    branch: string;
    file_path: string;
    
}


const JSONToDeploymentReadSummary = (m: DeploymentReadSummary | DeploymentReadSummaryJSON): DeploymentReadSummary => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        repoId: (((m as DeploymentReadSummary).repoId) ? (m as DeploymentReadSummary).repoId : (m as DeploymentReadSummaryJSON).repo_id),
        repoName: (((m as DeploymentReadSummary).repoName) ? (m as DeploymentReadSummary).repoName : (m as DeploymentReadSummaryJSON).repo_name),
        branch: m.branch,
        filePath: (((m as DeploymentReadSummary).filePath) ? (m as DeploymentReadSummary).filePath : (m as DeploymentReadSummaryJSON).file_path),
        
    };
};

export interface ReadDeployments {
    
}

interface ReadDeploymentsJSON {
    
}


const ReadDeploymentsToJSON = (m: ReadDeployments): ReadDeploymentsJSON => {
    return {
        
    };
};

export interface DeploymentsRead {
    deployments: DeploymentReadSummary[];
    
}

interface DeploymentsReadJSON {
    deployments: DeploymentReadSummaryJSON[];
    
}


const JSONToDeploymentsRead = (m: DeploymentsRead | DeploymentsReadJSON): DeploymentsRead => {
    
    return {
        deployments: (m.deployments as (DeploymentReadSummary | DeploymentReadSummaryJSON)[]).map(JSONToDeploymentReadSummary),
        
    };
};

export interface Deployment {
    create: (createDeployment: CreateDeployment) => Promise<DeploymentCreated>;
    
    update: (updateDeployment: UpdateDeployment) => Promise<DeploymentUpdated>;
    
    destroy: (destroyDeployment: DestroyDeployment) => Promise<DeploymentDestroyed>;
    
    read: (readDeployment: ReadDeployment) => Promise<DeploymentRead>;
    
    all: (readDeployments: ReadDeployments) => Promise<DeploymentsRead>;
    
}

export class DefaultDeployment implements Deployment {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Deployment/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    create(createDeployment: CreateDeployment): Promise<DeploymentCreated> {
        const url = this.hostname + this.pathPrefix + "Create";
        let body: CreateDeployment | CreateDeploymentJSON = createDeployment;
        if (!this.writeCamelCase) {
            body = CreateDeploymentToJSON(createDeployment);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeploymentCreated);
        });
    }
    
    update(updateDeployment: UpdateDeployment): Promise<DeploymentUpdated> {
        const url = this.hostname + this.pathPrefix + "Update";
        let body: UpdateDeployment | UpdateDeploymentJSON = updateDeployment;
        if (!this.writeCamelCase) {
            body = UpdateDeploymentToJSON(updateDeployment);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeploymentUpdated);
        });
    }
    
    destroy(destroyDeployment: DestroyDeployment): Promise<DeploymentDestroyed> {
        const url = this.hostname + this.pathPrefix + "Destroy";
        let body: DestroyDeployment | DestroyDeploymentJSON = destroyDeployment;
        if (!this.writeCamelCase) {
            body = DestroyDeploymentToJSON(destroyDeployment);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeploymentDestroyed);
        });
    }
    
    read(readDeployment: ReadDeployment): Promise<DeploymentRead> {
        const url = this.hostname + this.pathPrefix + "Read";
        let body: ReadDeployment | ReadDeploymentJSON = readDeployment;
        if (!this.writeCamelCase) {
            body = ReadDeploymentToJSON(readDeployment);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeploymentRead);
        });
    }
    
    all(readDeployments: ReadDeployments): Promise<DeploymentsRead> {
        const url = this.hostname + this.pathPrefix + "All";
        let body: ReadDeployments | ReadDeploymentsJSON = readDeployments;
        if (!this.writeCamelCase) {
            body = ReadDeploymentsToJSON(readDeployments);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeploymentsRead);
        });
    }
    
}

