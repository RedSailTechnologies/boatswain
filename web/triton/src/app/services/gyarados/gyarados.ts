
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface Delivery {
    uuid: string;
    name: string;
    version: string;
    application: Application;
    clusters: string[];
    deployments: Deployment[];
    tests: Deployment[];
    triggers: Trigger[];
    strategy: Step[];
    
}

interface DeliveryJSON {
    uuid: string;
    name: string;
    version: string;
    application: ApplicationJSON;
    clusters: string[];
    deployments: DeploymentJSON[];
    tests: DeploymentJSON[];
    triggers: TriggerJSON[];
    strategy: StepJSON[];
    
}


export interface Deployment {
    uuid: string;
    name: string;
    docker: Docker;
    helm: Helm;
    template: string;
    arguments: string;
    
}

interface DeploymentJSON {
    uuid: string;
    name: string;
    docker: DockerJSON;
    helm: HelmJSON;
    template: string;
    arguments: string;
    
}


export interface Trigger {
    uuid: string;
    name: string;
    approval: Approval;
    delivery: Delivery;
    manual: Approval;
    web: Web;
    template: string;
    arguments: string;
    
}

interface TriggerJSON {
    uuid: string;
    name: string;
    approval: ApprovalJSON;
    delivery: DeliveryJSON;
    manual: ApprovalJSON;
    web: WebJSON;
    template: string;
    arguments: string;
    
}


export interface Step {
    uuid: string;
    name: string;
    displayname: string;
    success: StepAction[];
    failure: StepAction[];
    any: StepAction[];
    always: StepAction[];
    hold: string;
    template: string;
    arguments: string;
    
}

interface StepJSON {
    uuid: string;
    name: string;
    displayName: string;
    success: StepActionJSON[];
    failure: StepActionJSON[];
    any: StepActionJSON[];
    always: StepActionJSON[];
    hold: string;
    template: string;
    arguments: string;
    
}


export interface Template {
    uuid: string;
    name: string;
    type: string;
    yaml: string;
    
}

interface TemplateJSON {
    uuid: string;
    name: string;
    type: string;
    yaml: string;
    
}


