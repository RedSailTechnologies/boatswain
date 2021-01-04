
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface CheckLive {
    
}

interface CheckLiveJSON {
    
}


const CheckLiveToJSON = (m: CheckLive): CheckLiveJSON => {
    return {
        
    };
};

export interface LiveCheck {
    
}

interface LiveCheckJSON {
    
}


const JSONToLiveCheck = (m: LiveCheck | LiveCheckJSON): LiveCheck => {
    
    return {
        
    };
};

export interface CheckReady {
    
}

interface CheckReadyJSON {
    
}


const CheckReadyToJSON = (m: CheckReady): CheckReadyJSON => {
    return {
        
    };
};

export interface ReadyCheck {
    
}

interface ReadyCheckJSON {
    
}


const JSONToReadyCheck = (m: ReadyCheck | ReadyCheckJSON): ReadyCheck => {
    
    return {
        
    };
};

export interface Health {
    live: (checkLive: CheckLive) => Promise<LiveCheck>;
    
    ready: (checkReady: CheckReady) => Promise<ReadyCheck>;
    
}

export class DefaultHealth implements Health {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Health/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    live(checkLive: CheckLive): Promise<LiveCheck> {
        const url = this.hostname + this.pathPrefix + "Live";
        let body: CheckLive | CheckLiveJSON = checkLive;
        if (!this.writeCamelCase) {
            body = CheckLiveToJSON(checkLive);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToLiveCheck);
        });
    }
    
    ready(checkReady: CheckReady): Promise<ReadyCheck> {
        const url = this.hostname + this.pathPrefix + "Ready";
        let body: CheckReady | CheckReadyJSON = checkReady;
        if (!this.writeCamelCase) {
            body = CheckReadyToJSON(checkReady);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToReadyCheck);
        });
    }
    
}

