import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Application, DefaultKraken, Kraken } from 'src/app/services/kraken/kraken';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-projects',
  templateUrl: './projects.component.html',
  styleUrls: ['./projects.component.sass']
})
export class ProjectsComponent implements OnInit {
  private client: Kraken;
  private retries = 0;
  public applications: Application[];
  projects: Project[];

  constructor(private snackBar: MatSnackBar) {
    this.client = new DefaultKraken(`${location.protocol}//${location.host}/api`, fetch['default']);
  }

  ngOnInit(): void {
    this.refresh();
  }

  refresh(): void {
    if (this.retries < 5) {
      this.retries++;
      this.client.applications({}).then(value => {
        this.applications = value.applications;
        this.populateProjects();
      }).catch(_ => {
        setTimeout(() => this.refresh(), 2 * 1000)
      });
    } else {
      console.log("could not update applications after 5 retries");
      this.retries = 0;
      this.applications = new Array<Application>();
      this.snackBar.open(`Error getting projects`, "Dismiss", {
        duration: 5 * 1000,
        panelClass: ["warn-snack"]
      });
    }
  }

  populateProjects(): void {
    this.applications.forEach(app => {
      if (this.projects.find(x => x.name == app.project) == undefined) {
        var project = new Project();
        project.name = app.project;
        project.applications = new Array<ProjectApp>();
        this.projects.push(project);
      }

      var pApp = new ProjectApp();
      pApp.name = app.name;
      pApp.clusters = new Array<string>();
      app.clusters.forEach(x => pApp.clusters.push(x.clusterName));
      this.projects.forEach(x => {
        if (x.name == app.project) {
          x.applications.push(pApp);
        }
      });
    });
  }

}

class Project {
  public name: string;
  public applications: ProjectApp[]
}

class ProjectApp {
  public name: string;
  public clusters: string[];
}
