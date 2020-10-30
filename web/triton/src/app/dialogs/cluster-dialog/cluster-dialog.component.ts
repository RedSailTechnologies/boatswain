import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Cluster, DefaultKraken, Kraken } from 'src/app/services/kraken/kraken';
import * as fetch from 'isomorphic-fetch';
import { BusyComponent } from '../busy/busy.component';
import { MessageDialogComponent } from '../message-dialog/message-dialog.component';

@Component({
  selector: 'app-cluster-dialog',
  templateUrl: './cluster-dialog.component.html',
  styleUrls: ['./cluster-dialog.component.sass']
})
export class ClusterDialogComponent implements OnInit {
  private client: Kraken;
  private cluster: Cluster;
  public clusterForm: FormGroup = new FormGroup({
    name: new FormControl(''),
    endpoint: new FormControl(''),
    token: new FormControl(''),
    cert: new FormControl(''),
  });
  public isAdd: boolean;
  public title: string;

  constructor(public dialog: MatDialogRef<ClusterDialogComponent>, 
              @Inject(MAT_DIALOG_DATA) data,
              private spinner: MatDialog,
              private error: MatDialog) {
    this.title = data["title"];
    this.isAdd = data["type"] == "add";
    if (!this.isAdd) {
      this.cluster = data["cluster"];
      this.clusterForm.controls["name"].setValue(this.cluster.name);
      this.clusterForm.controls["endpoint"].setValue(this.cluster.endpoint);
      this.clusterForm.controls["token"].setValue("***");
      this.clusterForm.controls["cert"].setValue("***");
    }
    this.client = new DefaultKraken(`${location.protocol}//${location.host}/api`, fetch["default"]);
  }

  ngOnInit(): void {
  }

  enter($event): void {
    if ($event.keyCode == 13 && this.clusterForm.valid) {
      this.submit()
    };
  }

  submit(): void {
    var spinnerRef: MatDialogRef<BusyComponent> = this.spinner.open(BusyComponent, {
      panelClass: 'transparent',
      disableClose: true
    });
    var cluster = <Cluster>{
      "uuid": this.cluster != null ? this.cluster.uuid : null,
      "name": this.clusterForm.controls["name"].value,
      "endpoint": this.clusterForm.controls["endpoint"].value,
      "token": this.clusterForm.controls["token"].value == "***" ? this.cluster.token : this.clusterForm.controls["token"].value,
      "cert": this.clusterForm.controls["cert"].value == "***" ? this.cluster.cert : this.clusterForm.controls["cert"].value,
      "ready": false
    };

    var promise: Promise<any>
    if (this.isAdd) {
      promise = this.client.addCluster(cluster);
    } else {
      promise = this.client.editCluster(cluster);
    }

    promise.then(_ => {
      spinnerRef.close()
      this.dialog.close()
    }).catch(error => {
      spinnerRef.close();
      this.error.open(MessageDialogComponent, {
        panelClass: 'message-box',
        data: {
          "reason": "Error",
          "message": "An error occured.\n" + error
        }
      });
    });
  }

  secretFocused(control: string) : void {
    if (this.clusterForm.controls[control].value == "***") {
      this.clusterForm.controls[control].setValue("");
    }
  }
}
