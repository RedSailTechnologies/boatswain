import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './pages/home/home.component';
import { ClustersComponent } from './pages/clusters/clusters.component';
import { ReposComponent } from './pages/repos/repos.component';
import { ProjectsComponent } from './pages/projects/projects.component';
import { ApplicationsComponent } from './pages/applications/applications.component';

const routes: Routes = [
    {path: '', component: HomeComponent},
    {path: 'clusters', component: ClustersComponent},
    {path: 'repos', component: ReposComponent},
    {path: 'projects', component: ProjectsComponent},
    {path: 'applications', component: ApplicationsComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
