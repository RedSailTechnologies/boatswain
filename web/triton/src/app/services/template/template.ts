
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface CreateTemplate {
    name: string;
    type: string;
    yaml: string;
    
}

interface CreateTemplateJSON {
    name: string;
    type: string;
    yaml: string;
    
}


const CreateTemplateToJSON = (m: CreateTemplate): CreateTemplateJSON => {
    return {
        name: m.name,
        type: m.type,
        yaml: m.yaml,
        
    };
};

export interface TemplateCreated {
    
}

interface TemplateCreatedJSON {
    
}


const JSONToTemplateCreated = (m: TemplateCreated | TemplateCreatedJSON): TemplateCreated => {
    
    return {
        
    };
};

export interface UpdateTemplate {
    uuid: string;
    name: string;
    type: string;
    yaml: string;
    
}

interface UpdateTemplateJSON {
    uuid: string;
    name: string;
    type: string;
    yaml: string;
    
}


const UpdateTemplateToJSON = (m: UpdateTemplate): UpdateTemplateJSON => {
    return {
        uuid: m.uuid,
        name: m.name,
        type: m.type,
        yaml: m.yaml,
        
    };
};

export interface TemplateUpdated {
    
}

interface TemplateUpdatedJSON {
    
}


const JSONToTemplateUpdated = (m: TemplateUpdated | TemplateUpdatedJSON): TemplateUpdated => {
    
    return {
        
    };
};

export interface DestroyTemplate {
    uuid: string;
    
}

interface DestroyTemplateJSON {
    uuid: string;
    
}


const DestroyTemplateToJSON = (m: DestroyTemplate): DestroyTemplateJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface TemplateDestroyed {
    
}

interface TemplateDestroyedJSON {
    
}


const JSONToTemplateDestroyed = (m: TemplateDestroyed | TemplateDestroyedJSON): TemplateDestroyed => {
    
    return {
        
    };
};

export interface ReadTemplate {
    uuid: string;
    
}

interface ReadTemplateJSON {
    uuid: string;
    
}


const ReadTemplateToJSON = (m: ReadTemplate): ReadTemplateJSON => {
    return {
        uuid: m.uuid,
        
    };
};

export interface TemplateRead {
    uuid: string;
    name: string;
    type: string;
    yaml: string;
    
}

interface TemplateReadJSON {
    uuid: string;
    name: string;
    type: string;
    yaml: string;
    
}


const JSONToTemplateRead = (m: TemplateRead | TemplateReadJSON): TemplateRead => {
    
    return {
        uuid: m.uuid,
        name: m.name,
        type: m.type,
        yaml: m.yaml,
        
    };
};

export interface ReadTemplates {
    
}

interface ReadTemplatesJSON {
    
}


const ReadTemplatesToJSON = (m: ReadTemplates): ReadTemplatesJSON => {
    return {
        
    };
};

export interface TemplatesRead {
    templates: TemplateRead[];
    
}

interface TemplatesReadJSON {
    templates: TemplateReadJSON[];
    
}


const JSONToTemplatesRead = (m: TemplatesRead | TemplatesReadJSON): TemplatesRead => {
    
    return {
        templates: (m.templates as (TemplateRead | TemplateReadJSON)[]).map(JSONToTemplateRead),
        
    };
};

export interface Template {
    create: (createTemplate: CreateTemplate) => Promise<TemplateCreated>;
    
    update: (updateTemplate: UpdateTemplate) => Promise<TemplateUpdated>;
    
    destroy: (destroyTemplate: DestroyTemplate) => Promise<TemplateDestroyed>;
    
    read: (readTemplate: ReadTemplate) => Promise<TemplateRead>;
    
    all: (readTemplates: ReadTemplates) => Promise<TemplatesRead>;
    
}

export class DefaultTemplate implements Template {
    private hostname: string;
    private fetch: Fetch;
    private writeCamelCase: boolean;
    private pathPrefix = "/redsail.bosn.Template/";

    constructor(hostname: string, fetch: Fetch, writeCamelCase = false) {
        this.hostname = hostname;
        this.fetch = fetch;
        this.writeCamelCase = writeCamelCase;
    }
    create(createTemplate: CreateTemplate): Promise<TemplateCreated> {
        const url = this.hostname + this.pathPrefix + "Create";
        let body: CreateTemplate | CreateTemplateJSON = createTemplate;
        if (!this.writeCamelCase) {
            body = CreateTemplateToJSON(createTemplate);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToTemplateCreated);
        });
    }
    
    update(updateTemplate: UpdateTemplate): Promise<TemplateUpdated> {
        const url = this.hostname + this.pathPrefix + "Update";
        let body: UpdateTemplate | UpdateTemplateJSON = updateTemplate;
        if (!this.writeCamelCase) {
            body = UpdateTemplateToJSON(updateTemplate);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToTemplateUpdated);
        });
    }
    
    destroy(destroyTemplate: DestroyTemplate): Promise<TemplateDestroyed> {
        const url = this.hostname + this.pathPrefix + "Destroy";
        let body: DestroyTemplate | DestroyTemplateJSON = destroyTemplate;
        if (!this.writeCamelCase) {
            body = DestroyTemplateToJSON(destroyTemplate);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToTemplateDestroyed);
        });
    }
    
    read(readTemplate: ReadTemplate): Promise<TemplateRead> {
        const url = this.hostname + this.pathPrefix + "Read";
        let body: ReadTemplate | ReadTemplateJSON = readTemplate;
        if (!this.writeCamelCase) {
            body = ReadTemplateToJSON(readTemplate);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToTemplateRead);
        });
    }
    
    all(readTemplates: ReadTemplates): Promise<TemplatesRead> {
        const url = this.hostname + this.pathPrefix + "All";
        let body: ReadTemplates | ReadTemplatesJSON = readTemplates;
        if (!this.writeCamelCase) {
            body = ReadTemplatesToJSON(readTemplates);
        }
        return this.fetch(createTwirpRequest(url, body)).then((resp) => {
            if (!resp.ok) {
                return throwTwirpError(resp);
            }

            return resp.json().then(JSONToTemplatesRead);
        });
    }
    
}

