
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
    
}

interface DeploymentReadJSON {
    uuid: string;
    name: string;
    repo_id: string;
    repo_name: string;
    branch: string;
    file_path: string;
    
}


const JSONToDeploymentRead = (m: DeploymentRead | DeploymentReadJSON): DeploymentRead => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        repoId: (((m as DeploymentRead).repoId) ? (m as DeploymentRead).repoId : (m as DeploymentReadJSON).repo_id),
        repoName: (((m as DeploymentRead).repoName) ? (m as DeploymentRead).repoName : (m as DeploymentReadJSON).repo_name),
        branch: m.branch,
        filePath: (((m as DeploymentRead).filePath) ? (m as DeploymentRead).filePath : (m as DeploymentReadJSON).file_path),
        
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

export interface DeploymentReadSummary {
    uuid: string;
    name: string;
    
}

interface DeploymentReadSummaryJSON {
    uuid: string;
    name: string;
    
}


const JSONToDeploymentReadSummary = (m: DeploymentReadSummary | DeploymentReadSummaryJSON): DeploymentReadSummary => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        
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

export interface TemplateDeployment {
    uuid: string;
    
}

interface TemplateDeploymentJSON {
    uuid: string;
    
}


const TemplateDeploymentToJSON = (m: TemplateDeployment): TemplateDeploymentJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface DeploymentTemplated {
    uuid: string;
    yaml: string;
    
}

interface DeploymentTemplatedJSON {
    uuid: string;
    yaml: string;
    
}


const JSONToDeploymentTemplated = (m: DeploymentTemplated | DeploymentTemplatedJSON): DeploymentTemplated => {
    
    return {
        uuid: m.uuid,
        yaml: m.yaml,
        
    };
};

export interface ReadDeploymentToken {
    uuid: string;
    
}

interface ReadDeploymentTokenJSON {
    uuid: string;
    
}


const ReadDeploymentTokenToJSON = (m: ReadDeploymentToken): ReadDeploymentTokenJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface DeploymentTokenRead {
    token: string;
    
}

interface DeploymentTokenReadJSON {
    token: string;
    
}


const JSONToDeploymentTokenRead = (m: DeploymentTokenRead | DeploymentTokenReadJSON): DeploymentTokenRead => {
    
    return {
        token: m.token,
        
    };
};

export interface ReadRun {
    deploymentUuid: string;
    
}

interface ReadRunJSON {
    deployment_uuid: string;
    
}


const ReadRunToJSON = (m: ReadRun): ReadRunJSON => {
    return {
        deployment_uuid: m.deploymentUuid,
        
    };
};

export interface StepLog {
    timestamp: number;
    level: string;
    message: string;
    
}

interface StepLogJSON {
    timestamp: number;
    level: string;
    message: string;
    
}


const JSONToStepLog = (m: StepLog | StepLogJSON): StepLog => {
    
    return {
        timestamp: m.timestamp,
        level: m.level,
        message: m.message,
        
    };
};

export interface StepRead {
    name: string;
    status: string;
    startTime: number;
    stopTime: number;
    logs: StepLog[];
    
}

interface StepReadJSON {
    name: string;
    status: string;
    start_time: number;
    stop_time: number;
    logs: StepLogJSON[];
    
}


const JSONToStepRead = (m: StepRead | StepReadJSON): StepRead => {
    
    return {
        name: m.name,
        status: m.status,
        startTime: (((m as StepRead).startTime) ? (m as StepRead).startTime : (m as StepReadJSON).start_time),
        stopTime: (((m as StepRead).stopTime) ? (m as StepRead).stopTime : (m as StepReadJSON).stop_time),
        logs: (m.logs as (StepLog | StepLogJSON)[]).map(JSONToStepLog),
        
    };
};

export interface LinkRead {
    name: string;
    url: string;
    
}

interface LinkReadJSON {
    name: string;
    url: string;
    
}


const JSONToLinkRead = (m: LinkRead | LinkReadJSON): LinkRead => {
    
    return {
        name: m.name,
        url: m.url,
        
    };
};

export interface RunRead {
    uuid: string;
    name: string;
    version: string;
    status: string;
    startTime: number;
    stopTime: number;
    links: LinkRead[];
    steps: StepRead[];
    
}

interface RunReadJSON {
    uuid: string;
    name: string;
    version: string;
    status: string;
    start_time: number;
    stop_time: number;
    links: LinkReadJSON[];
    steps: StepReadJSON[];
    
}


const JSONToRunRead = (m: RunRead | RunReadJSON): RunRead => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        version: m.version,
        status: m.status,
        startTime: (((m as RunRead).startTime) ? (m as RunRead).startTime : (m as RunReadJSON).start_time),
        stopTime: (((m as RunRead).stopTime) ? (m as RunRead).stopTime : (m as RunReadJSON).stop_time),
        links: (m.links as (LinkRead | LinkReadJSON)[]).map(JSONToLinkRead),
        steps: (m.steps as (StepRead | StepReadJSON)[]).map(JSONToStepRead),
        
    };
};

export interface ReadRuns {
    deploymentUuid: string;
    
}

interface ReadRunsJSON {
    deployment_uuid: string;
    
}


const ReadRunsToJSON = (m: ReadRuns): ReadRunsJSON => {
    return {
        deployment_uuid: m.deploymentUuid,
        
    };
};

export interface RunReadSummary {
    uuid: string;
    name: string;
    version: string;
    status: string;
    startTime: number;
    stopTime: number;
    
}

interface RunReadSummaryJSON {
    uuid: string;
    name: string;
    version: string;
    status: string;
    start_time: number;
    stop_time: number;
    
}


const JSONToRunReadSummary = (m: RunReadSummary | RunReadSummaryJSON): RunReadSummary => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        version: m.version,
        status: m.status,
        startTime: (((m as RunReadSummary).startTime) ? (m as RunReadSummary).startTime : (m as RunReadSummaryJSON).start_time),
        stopTime: (((m as RunReadSummary).stopTime) ? (m as RunReadSummary).stopTime : (m as RunReadSummaryJSON).stop_time),
        
    };
};

export interface RunsRead {
    runs: RunReadSummary[];
    
}

interface RunsReadJSON {
    runs: RunReadSummaryJSON[];
    
}


const JSONToRunsRead = (m: RunsRead | RunsReadJSON): RunsRead => {
    
    return {
        runs: (m.runs as (RunReadSummary | RunReadSummaryJSON)[]).map(JSONToRunReadSummary),
        
    };
};

export interface ApproveStep {
    runUuid: string;
    approve: boolean;
    override: boolean;
    
}

interface ApproveStepJSON {
    run_uuid: string;
    approve: boolean;
    override: boolean;
    
}


const ApproveStepToJSON = (m: ApproveStep): ApproveStepJSON => {
    return {
        run_uuid: m.runUuid,
        approve: m.approve,
        override: m.override,
        
    };
};

export interface StepApproved {
    
}

interface StepApprovedJSON {
    
}


const JSONToStepApproved = (m: StepApproved | StepApprovedJSON): StepApproved => {
    
    return {
        
    };
};

export interface ReadApprovals {
    
}

interface ReadApprovalsJSON {
    
}


const ReadApprovalsToJSON = (m: ReadApprovals): ReadApprovalsJSON => {
    return {
        
    };
};

export interface ApprovalRead {
    uuid: string;
    runUuid: string;
    runName: string;
    runVersion: string;
    stepName: string;
    
}

interface ApprovalReadJSON {
    uuid: string;
    run_uuid: string;
    run_name: string;
    run_version: string;
    step_name: string;
    
}


const JSONToApprovalRead = (m: ApprovalRead | ApprovalReadJSON): ApprovalRead => {
    
    return {
        uuid: m.uuid,
        runUuid: (((m as ApprovalRead).runUuid) ? (m as ApprovalRead).runUuid : (m as ApprovalReadJSON).run_uuid),
        runName: (((m as ApprovalRead).runName) ? (m as ApprovalRead).runName : (m as ApprovalReadJSON).run_name),
        runVersion: (((m as ApprovalRead).runVersion) ? (m as ApprovalRead).runVersion : (m as ApprovalReadJSON).run_version),
        stepName: (((m as ApprovalRead).stepName) ? (m as ApprovalRead).stepName : (m as ApprovalReadJSON).step_name),
        
    };
};

export interface ApprovalsRead {
    approvals: ApprovalRead[];
    
}

interface ApprovalsReadJSON {
    approvals: ApprovalReadJSON[];
    
}


const JSONToApprovalsRead = (m: ApprovalsRead | ApprovalsReadJSON): ApprovalsRead => {
    
    return {
        approvals: (m.approvals as (ApprovalRead | ApprovalReadJSON)[]).map(JSONToApprovalRead),
        
    };
};

export interface Deployment {
    create: (createDeployment: CreateDeployment) => Promise<DeploymentCreated>;
    
    update: (updateDeployment: UpdateDeployment) => Promise<DeploymentUpdated>;
    
    destroy: (destroyDeployment: DestroyDeployment) => Promise<DeploymentDestroyed>;
    
    read: (readDeployment: ReadDeployment) => Promise<DeploymentRead>;
    
    all: (readDeployments: ReadDeployments) => Promise<DeploymentsRead>;
    
    template: (templateDeployment: TemplateDeployment) => Promise<DeploymentTemplated>;
    
    token: (readDeploymentToken: ReadDeploymentToken) => Promise<DeploymentTokenRead>;
    
    run: (readRun: ReadRun) => Promise<RunRead>;
    
    runs: (readRuns: ReadRuns) => Promise<RunsRead>;
    
    approve: (approveStep: ApproveStep) => Promise<StepApproved>;
    
    approvals: (readApprovals: ReadApprovals) => Promise<ApprovalsRead>;
    
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
    
    template(templateDeployment: TemplateDeployment): Promise<DeploymentTemplated> {
        const url = this.hostname + this.pathPrefix + "Template";
        let body: TemplateDeployment | TemplateDeploymentJSON = templateDeployment;
        if (!this.writeCamelCase) {
            body = TemplateDeploymentToJSON(templateDeployment);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeploymentTemplated);
        });
    }
    
    token(readDeploymentToken: ReadDeploymentToken): Promise<DeploymentTokenRead> {
        const url = this.hostname + this.pathPrefix + "Token";
        let body: ReadDeploymentToken | ReadDeploymentTokenJSON = readDeploymentToken;
        if (!this.writeCamelCase) {
            body = ReadDeploymentTokenToJSON(readDeploymentToken);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToDeploymentTokenRead);
        });
    }
    
    run(readRun: ReadRun): Promise<RunRead> {
        const url = this.hostname + this.pathPrefix + "Run";
        let body: ReadRun | ReadRunJSON = readRun;
        if (!this.writeCamelCase) {
            body = ReadRunToJSON(readRun);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToRunRead);
        });
    }
    
    runs(readRuns: ReadRuns): Promise<RunsRead> {
        const url = this.hostname + this.pathPrefix + "Runs";
        let body: ReadRuns | ReadRunsJSON = readRuns;
        if (!this.writeCamelCase) {
            body = ReadRunsToJSON(readRuns);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToRunsRead);
        });
    }
    
    approve(approveStep: ApproveStep): Promise<StepApproved> {
        const url = this.hostname + this.pathPrefix + "Approve";
        let body: ApproveStep | ApproveStepJSON = approveStep;
        if (!this.writeCamelCase) {
            body = ApproveStepToJSON(approveStep);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToStepApproved);
        });
    }
    
    approvals(readApprovals: ReadApprovals): Promise<ApprovalsRead> {
        const url = this.hostname + this.pathPrefix + "Approvals";
        let body: ReadApprovals | ReadApprovalsJSON = readApprovals;
        if (!this.writeCamelCase) {
            body = ReadApprovalsToJSON(readApprovals);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToApprovalsRead);
        });
    }
    
}

