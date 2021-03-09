
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface TriggerManual {
    uuid: string;
    name: string;
    args: string;
    
}

interface TriggerManualJSON {
    uuid: string;
    name: string;
    args: string;
    
}


const TriggerManualToJSON = (m: TriggerManual): TriggerManualJSON => {
    return {
        uuid: m.uuid,
        name: m.name,
        args: m.args,
        
    };
};

export interface ManualTriggered {
    runUuid: string;
    
}

interface ManualTriggeredJSON {
    run_uuid: string;
    
}


const JSONToManualTriggered = (m: ManualTriggered | ManualTriggeredJSON): ManualTriggered => {
    
    return {
        runUuid: (((m as ManualTriggered).runUuid) ? (m as ManualTriggered).runUuid : (m as ManualTriggeredJSON).run_uuid),
        
    };
};

export interface TriggerWeb {
    uuid: string;
    name: string;
    token: string;
    args: string;
    
}

interface TriggerWebJSON {
    uuid: string;
    name: string;
    token: string;
    args: string;
    
}


const TriggerWebToJSON = (m: TriggerWeb): TriggerWebJSON => {
    return {
        uuid: m.uuid,
        name: m.name,
        token: m.token,
        args: m.args,
        
    };
};

export interface WebTriggered {
    runUuid: string;
    
}

interface WebTriggeredJSON {
    run_uuid: string;
    
}


const JSONToWebTriggered = (m: WebTriggered | WebTriggeredJSON): WebTriggered => {
    
    return {
        runUuid: (((m as WebTriggered).runUuid) ? (m as WebTriggered).runUuid : (m as WebTriggeredJSON).run_uuid),
        
    };
};

export interface ReadStatus {
    deploymentUuid: string;
    deploymentToken: string;
    runUuid: string;
    
}

interface ReadStatusJSON {
    deployment_uuid: string;
    deployment_token: string;
    run_uuid: string;
    
}


const ReadStatusToJSON = (m: ReadStatus): ReadStatusJSON => {
    return {
        deployment_uuid: m.deploymentUuid,
        deployment_token: m.deploymentToken,
        run_uuid: m.runUuid,
        
    };
};

export interface StatusRead {
    status: string;
    
}

interface StatusReadJSON {
    status: string;
    
}


const JSONToStatusRead = (m: StatusRead | StatusReadJSON): StatusRead => {
    
    return {
        status: m.status,
        
    };
};

export interface Trigger {
    manual: (triggerManual: TriggerManual) => Promise<ManualTriggered>;
    
    web: (triggerWeb: TriggerWeb) => Promise<WebTriggered>;
    
    status: (readStatus: ReadStatus) => Promise<StatusRead>;
    
}

export class DefaultTrigger implements Trigger {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Trigger/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    manual(triggerManual: TriggerManual): Promise<ManualTriggered> {
        const url = this.hostname + this.pathPrefix + "Manual";
        let body: TriggerManual | TriggerManualJSON = triggerManual;
        if (!this.writeCamelCase) {
            body = TriggerManualToJSON(triggerManual);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToManualTriggered);
        });
    }
    
    web(triggerWeb: TriggerWeb): Promise<WebTriggered> {
        const url = this.hostname + this.pathPrefix + "Web";
        let body: TriggerWeb | TriggerWebJSON = triggerWeb;
        if (!this.writeCamelCase) {
            body = TriggerWebToJSON(triggerWeb);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToWebTriggered);
        });
    }
    
    status(readStatus: ReadStatus): Promise<StatusRead> {
        const url = this.hostname + this.pathPrefix + "Status";
        let body: ReadStatus | ReadStatusJSON = readStatus;
        if (!this.writeCamelCase) {
            body = ReadStatusToJSON(readStatus);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToStatusRead);
        });
    }
    
}

