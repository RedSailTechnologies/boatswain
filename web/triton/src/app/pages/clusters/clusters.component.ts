import { Component, OnInit } from "@angular/core";
import { Cluster, DefaultKraken, Kraken } from "src/app/services/kraken/kraken";
import * as fetch from "isomorphic-fetch";
import { MatDialog } from "@angular/material/dialog";
import { ClusterDialogComponent } from "src/app/dialogs/cluster-dialog/cluster-dialog.component";

@Component({
  selector: "app-clusters",
  templateUrl: "./clusters.component.html",
  styleUrls: ["./clusters.component.sass"]
})

export class ClustersComponent implements OnInit {
  private client: Kraken;
  public clusters: Cluster[];

  constructor(public dialog: MatDialog) {
    this.client = new DefaultKraken(`${location.protocol}//${location.host}/api`, fetch["default"]);
  }

  ngOnInit() : void {
    this.client.clusters({}).then(value => {
      this.clusters = value.clusters;
    });
  }

  add() : void {
    this.dialog.open(ClusterDialogComponent, {
      minWidth: "33%",
      panelClass: "custom-dialog-container",
      data: {
        "type": "add",
        "title": "Add Cluster"
      }
    });
  }

  edit(element: Cluster) : void {
    this.dialog.open(ClusterDialogComponent, {
      minWidth: "33%",
      panelClass: "custom-dialog-container",
      data: {
        "type": "edit",
        "title": `Edit ${element.name}`,
        "cluster": element
      }
    });
  }
}
