import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Application, DefaultKraken, Kraken } from 'src/app/services/kraken/kraken';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-applications',
  templateUrl: './applications.component.html',
  styleUrls: ['./applications.component.sass']
})
export class ApplicationsComponent implements OnInit {
  private client: Kraken;
  private retries = 0;
  public applications: Application[];

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
      }).catch(_ => {
        setTimeout(() => this.refresh(), 2 * 1000)
      });
    } else {
      console.log("could not update applications after 5 retries");
      this.retries = 0;
      this.applications = new Array<Application>();
      this.snackBar.open(`Error getting repos`, "Dismiss", {
        duration: 5 * 1000,
        panelClass: ["warn-snack"]
      });
    }
  }
}
