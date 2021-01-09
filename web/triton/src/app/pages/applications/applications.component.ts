import { Component, OnInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { MatSnackBar } from '@angular/material/snack-bar';
import {
  ApplicationRead,
  DefaultApplication,
  Application,
} from 'src/app/services/application/application';
import { AuthService } from 'src/app/utils/auth/auth.service';

@Component({
  selector: 'app-applications',
  templateUrl: './applications.component.html',
  styleUrls: ['./applications.component.sass'],
})
export class ApplicationsComponent implements OnInit {
  private client: Application;
  private retries = 0;
  private rawApplications: ApplicationRead[];
  public applications: ApplicationRead[];
  public projects: string[];
  public selectedProject: string = "";

  constructor(private snackBar: MatSnackBar, public auth: AuthService) {
    this.client = new DefaultApplication(
      `${location.protocol}//${location.host}/api`,
      auth.fetch()
    );
  }

  ngOnInit(): void {
    this.refresh();
  }

  filterByProject(event: MatSelectChange): void {
    this.selectedProject = event.value;
    if (this.selectedProject === "") {
      this.applications = this.rawApplications;
    } else {
      this.applications = this.rawApplications.filter(x => x.project == this.selectedProject);
    }
  }

  refresh(): void {
    if (this.retries < 5) {
      this.retries++;
      this.client
        .all({})
        .then((value) => {
          this.rawApplications = value.applications;
          var set = new Set(this.rawApplications.map((x => { return x.project; })));
          set.add("");
          this.projects = Array.from(set).sort();
          this.selectedProject = "";
          this.applications = this.rawApplications;
        })
        .catch((_) => {
          setTimeout(() => this.refresh(), 2 * 1000);
        });
    } else {
      console.log('could not update applications after 5 retries');
      this.retries = 0;
      this.applications = new Array<ApplicationRead>();
      this.snackBar.open(`Error getting applications`, 'Dismiss', {
        duration: 5 * 1000,
        panelClass: ['warn-snack'],
      });
    }
  }
}
