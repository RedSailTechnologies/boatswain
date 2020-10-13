import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './pages/home/home.component';
import { ClustersComponent } from './pages/clusters/clusters.component';
import { ReleasesComponent } from './pages/releases/releases.component';
import { ReposComponent } from './pages/repos/repos.component';

const routes: Routes = [
    {path: '', component: HomeComponent},
    {path: 'clusters', component: ClustersComponent},
    {path: 'repos', component: ReposComponent},
    {path: 'releases', component: ReleasesComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
