
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface RegisterAgent {
    clusterUuid: string;
    
}

interface RegisterAgentJSON {
    cluster_uuid: string;
    
}


const RegisterAgentToJSON = (m: RegisterAgent): RegisterAgentJSON => {
    return {
        cluster_uuid: m.clusterUuid,
        
    };
};

export interface AgentRegistered {
    clusterToken: string;
    
}

interface AgentRegisteredJSON {
    cluster_token: string;
    
}


const JSONToAgentRegistered = (m: AgentRegistered | AgentRegisteredJSON): AgentRegistered => {
    
    return {
        clusterToken: (((m as AgentRegistered).clusterToken) ? (m as AgentRegistered).clusterToken : (m as AgentRegisteredJSON).cluster_token),
        
    };
};

export interface ReadActions {
    clusterUuid: string;
    clusterToken: string;
    
}

interface ReadActionsJSON {
    cluster_uuid: string;
    cluster_token: string;
    
}


const ReadActionsToJSON = (m: ReadActions): ReadActionsJSON => {
    return {
        cluster_uuid: m.clusterUuid,
        cluster_token: m.clusterToken,
        
    };
};

export interface ActionsRead {
    actions: Action[];
    
}

interface ActionsReadJSON {
    actions: ActionJSON[];
    
}


const JSONToActionsRead = (m: ActionsRead | ActionsReadJSON): ActionsRead => {
    
    return {
        actions: (m.actions as (Action | ActionJSON)[]).map(JSONToAction),
        
    };
};

export interface ReturnResult {
    actionUuid: string;
    clusterUuid: string;
    clusterToken: string;
    result: Result;
    
}

interface ReturnResultJSON {
    action_uuid: string;
    cluster_uuid: string;
    cluster_token: string;
    result: ResultJSON;
    
}


const ReturnResultToJSON = (m: ReturnResult): ReturnResultJSON => {
    return {
        action_uuid: m.actionUuid,
        cluster_uuid: m.clusterUuid,
        cluster_token: m.clusterToken,
        result: ResultToJSON(m.result),
        
    };
};

export interface ResultReturned {
    
}

interface ResultReturnedJSON {
    
}


const JSONToResultReturned = (m: ResultReturned | ResultReturnedJSON): ResultReturned => {
    
    return {
        
    };
};

export interface Action {
    uuid: string;
    clusterUuid: string;
    clusterToken: string;
    actionType: string;
    action: string;
    timeoutSeconds: number;
    args: string;
    
}

interface ActionJSON {
    uuid: string;
    cluster_uuid: string;
    cluster_token: string;
    action_type: string;
    action: string;
    timeout_seconds: number;
    args: string;
    
}


const ActionToJSON = (m: Action): ActionJSON => {
    return {
        uuid: m.uuid,
        cluster_uuid: m.clusterUuid,
        cluster_token: m.clusterToken,
        action_type: m.actionType,
        action: m.action,
        timeout_seconds: m.timeoutSeconds,
        args: m.args,
        
    };
};

const JSONToAction = (m: Action | ActionJSON): Action => {
    
    return {
        uuid: m.uuid,
        clusterUuid: (((m as Action).clusterUuid) ? (m as Action).clusterUuid : (m as ActionJSON).cluster_uuid),
        clusterToken: (((m as Action).clusterToken) ? (m as Action).clusterToken : (m as ActionJSON).cluster_token),
        actionType: (((m as Action).actionType) ? (m as Action).actionType : (m as ActionJSON).action_type),
        action: m.action,
        timeoutSeconds: (((m as Action).timeoutSeconds) ? (m as Action).timeoutSeconds : (m as ActionJSON).timeout_seconds),
        args: m.args,
        
    };
};

export interface Result {
    data: string;
    error: string;
    
}

interface ResultJSON {
    data: string;
    error: string;
    
}


const ResultToJSON = (m: Result): ResultJSON => {
    return {
        data: m.data,
        error: m.error,
        
    };
};

const JSONToResult = (m: Result | ResultJSON): Result => {
    
    return {
        data: m.data,
        error: m.error,
        
    };
};

export interface Agent {
    register: (registerAgent: RegisterAgent) => Promise<AgentRegistered>;
    
    actions: (readActions: ReadActions) => Promise<ActionsRead>;
    
    results: (returnResult: ReturnResult) => Promise<ResultReturned>;
    
}

export class DefaultAgent implements Agent {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Agent/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    register(registerAgent: RegisterAgent): Promise<AgentRegistered> {
        const url = this.hostname + this.pathPrefix + "Register";
        let body: RegisterAgent | RegisterAgentJSON = registerAgent;
        if (!this.writeCamelCase) {
            body = RegisterAgentToJSON(registerAgent);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToAgentRegistered);
        });
    }
    
    actions(readActions: ReadActions): Promise<ActionsRead> {
        const url = this.hostname + this.pathPrefix + "Actions";
        let body: ReadActions | ReadActionsJSON = readActions;
        if (!this.writeCamelCase) {
            body = ReadActionsToJSON(readActions);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToActionsRead);
        });
    }
    
    results(returnResult: ReturnResult): Promise<ResultReturned> {
        const url = this.hostname + this.pathPrefix + "Results";
        let body: ReturnResult | ReturnResultJSON = returnResult;
        if (!this.writeCamelCase) {
            body = ReturnResultToJSON(returnResult);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToResultReturned);
        });
    }
    
}

export interface AgentAction {
    run: (action: Action) => Promise<Result>;
    
}

export class DefaultAgentAction implements AgentAction {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.AgentAction/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    run(action: Action): Promise<Result> {
        const url = this.hostname + this.pathPrefix + "Run";
        let body: Action | ActionJSON = action;
        if (!this.writeCamelCase) {
            body = ActionToJSON(action);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToResult);
        });
    }
    
}

