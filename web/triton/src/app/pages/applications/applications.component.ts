import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ApplicationRead, DefaultApplication, Application } from 'src/app/services/application/application';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-applications',
  templateUrl: './applications.component.html',
  styleUrls: ['./applications.component.sass']
})
export class ApplicationsComponent implements OnInit {
  private client: Application;
  private retries = 0;
  public applications: ApplicationRead[];

  constructor(private snackBar: MatSnackBar) {
    this.client = new DefaultApplication(`${location.protocol}//${location.host}/api`, fetch['default']);
  }

  ngOnInit(): void {
    this.refresh();
  }

  refresh(): void {
    if (this.retries < 5) {
      this.retries++;
      this.client.all({}).then(value => {
        this.applications = value.applications;
      }).catch(_ => {
        setTimeout(() => this.refresh(), 2 * 1000)
      });
    } else {
      console.log("could not update applications after 5 retries");
      this.retries = 0;
      this.applications = new Array<ApplicationRead>();
      this.snackBar.open(`Error getting applications`, "Dismiss", {
        duration: 5 * 1000,
        panelClass: ["warn-snack"]
      });
    }
  }
}
