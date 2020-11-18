import { Component, OnInit } from "@angular/core";
import { ClusterRead, Cluster, DefaultCluster } from "src/app/services/cluster/cluster";
import * as fetch from "isomorphic-fetch";
import { MatDialog } from "@angular/material/dialog";
import { ClusterDialogComponent } from "src/app/dialogs/cluster-dialog/cluster-dialog.component";
import { ConfirmDialogComponent } from 'src/app/dialogs/confirm-dialog/confirm-dialog.component';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: "app-clusters",
  templateUrl: "./clusters.component.html",
  styleUrls: ["./clusters.component.sass"]
})

export class ClustersComponent implements OnInit {
  private client: Cluster;
  public clusters: ClusterRead[];
  private retries = 0;

  constructor(private dialog: MatDialog, private snackBar: MatSnackBar) {
    this.client = new DefaultCluster(`${location.protocol}//${location.host}/api`, fetch["default"]);
  }

  ngOnInit() : void {
    this.refreshClusters()
  }

  add() : void {
    this.dialog.open(ClusterDialogComponent, {
      minWidth: "33%",
      panelClass: "custom-dialog-container",
      data: {
        "type": "add",
        "title": "Add Cluster"
      }
    }).afterClosed().subscribe(_ => {
      this.refreshClusters();
    });
  }

  edit(element: ClusterRead) : void {
    this.dialog.open(ClusterDialogComponent, {
      minWidth: "33%",
      panelClass: "custom-dialog-container",
      data: {
        "type": "edit",
        "title": `Edit ${element.name}`,
        "cluster": element
      }
    }).afterClosed().subscribe(_ => {
      this.refreshClusters()
    });
  }

  delete(element: ClusterRead) : void {
    this.dialog.open(ConfirmDialogComponent, {
      panelClass: 'message-box',
      data: {
        "reason": `Delete ${element.name}`,
        "message": "Do you really want to delete this cluster?"
      }
    }).afterClosed().subscribe((result: Boolean) => {
      if (result) {
        this.client.destroy({"uuid": element.uuid}).catch(_ => {
          this.snackBar.open(`${element.name} could not be deleted`, "Dismiss", {
            duration: 5 * 1000,
            panelClass: ["warn-snack"]
          })
        }).then(() => {
          this.refreshClusters()
        });
      }
    });
  }

  refreshClusters() : void {
    if (this.retries < 5) {
      this.retries++;
      this.client.all({}).then(value => {
        this.clusters = value.clusters;
        this.retries = 0;
      }).catch(_ => {
        setTimeout(() => this.refreshClusters(), 2 * 1000)
      });
    } else {
      console.log("could not update clusters after 5 retries")
      this.retries = 0;
      this.clusters = new Array<ClusterRead>();
      this.snackBar.open(`Error getting clusters`, "Dismiss", {
        duration: 5 * 1000,
        panelClass: ["warn-snack"]
      })
    }
  }
}
