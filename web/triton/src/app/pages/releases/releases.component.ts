import { Component, OnInit } from '@angular/core';
import { DefaultKraken, Kraken, Release, Releases } from 'src/app/services/kraken/kraken';
import { MatDialog } from '@angular/material/dialog';
import { UpdateDialogComponent } from 'src/app/dialogs/update-dialog/update-dialog.component';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-releases',
  templateUrl: './releases.component.html',
  styleUrls: ['./releases.component.sass']
})
export class ReleasesComponent implements OnInit {
  private client: Kraken;
  public releasesList: Releases[];

  constructor(public dialog: MatDialog) {
    this.client = new DefaultKraken(`${location.protocol}//${location.host}/api`, fetch['default']);
  }

  ngOnInit(): void {
    this.client.clusters({}).then(clusters => {
      this.client.releases({"clusters": clusters.clusters}).then(releases => {
        this.releasesList = releases.releaseLists;
      })
    });
  }

  upgradeDialog(name: string, chart: string, cluster: Release) {
    const dialog = this.dialog.open(UpdateDialogComponent, {
      minWidth: "33%",
      panelClass: 'update-dialog-container',
      data: {"name": name, "chart": chart, "release": cluster}
    })
  }
}
