import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './pages/home/home.component';
import { ClustersComponent } from './pages/clusters/clusters.component';
import { ReposComponent } from './pages/repos/repos.component';
import { ProjectsComponent } from './pages/projects/projects.component';
import { ApplicationsComponent } from './pages/applications/applications.component';
import { LoginComponent } from './pages/login/login.component';
import { CheckService } from './utils/auth/check.service';
import { LogoutComponent } from './pages/logout/logout.component';
import { DeploymentsComponent } from './pages/deployments/deployments.component';
import { DeploymentComponent } from './pages/deployment/deployment.component';

const routes: Routes = [
  {path: '', component: HomeComponent},
  {path: 'login', component: LoginComponent},
  {path: 'logout', component: LogoutComponent},
  {path: 'clusters', component: ClustersComponent, canActivate: [CheckService]},
  {path: 'repos', component: ReposComponent, canActivate: [CheckService]},
  {path: 'projects', component: ProjectsComponent, canActivate: [CheckService]},
  {path: 'applications', component: ApplicationsComponent, canActivate: [CheckService]},
  {path: 'deployments', component: DeploymentsComponent, canActivate: [CheckService]},
  {path: 'deployment/:uuid', component: DeploymentComponent, canActivate: [CheckService]}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
