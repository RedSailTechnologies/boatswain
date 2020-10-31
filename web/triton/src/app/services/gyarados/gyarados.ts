
import {createTwirpRequest, throwTwirpError, Fetch} from './twirp';


export interface Deployment {
    name: string;
    docker: DockerDeployment;
    helm: HelmDeployment;
    
}

interface DeploymentJSON {
    name: string;
    docker: DockerDeploymentJSON;
    helm: HelmDeploymentJSON;
    
}


export interface DockerDeployment {
    image: string;
    tag: string;
    
}

interface DockerDeploymentJSON {
    image: string;
    tag: string;
    
}


export interface HelmDeployment {
    chart: string;
    repo: string;
    version: string;
    
}

interface HelmDeploymentJSON {
    chart: string;
    repo: string;
    version: string;
    
}


export interface Template {
    template: string;
    arguments: string;
    
}

interface TemplateJSON {
    template: string;
    arguments: string;
    
}


