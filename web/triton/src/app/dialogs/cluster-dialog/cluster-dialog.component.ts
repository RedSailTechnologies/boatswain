import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ClusterRead, DefaultCluster, Cluster, UpdateCluster, CreateCluster } from 'src/app/services/cluster/cluster';
import * as fetch from 'isomorphic-fetch';
import { BusyComponent } from '../busy/busy.component';
import { MessageDialogComponent } from '../message-dialog/message-dialog.component';
import { AuthService } from 'src/app/utils/auth/auth.service';

@Component({
  selector: 'app-cluster-dialog',
  templateUrl: './cluster-dialog.component.html',
  styleUrls: ['./cluster-dialog.component.sass']
})
export class ClusterDialogComponent implements OnInit {
  private client: Cluster;
  private cluster: ClusterRead;
  public clusterForm: FormGroup = new FormGroup({
    name: new FormControl(''),
  });
  public isAdd: boolean;
  public title: string;

  constructor(public dialog: MatDialogRef<ClusterDialogComponent>, 
              @Inject(MAT_DIALOG_DATA) data,
              private spinner: MatDialog,
              private error: MatDialog,
              private auth: AuthService) {
    this.title = data["title"];
    this.isAdd = data["type"] == "add";
    if (!this.isAdd) {
      this.cluster = data["cluster"];
      this.clusterForm.controls["name"].setValue(this.cluster.name);
    }
    this.client = new DefaultCluster(`${location.protocol}//${location.host}/api`, auth.fetch());
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
    

    var promise: Promise<any>
    var cluster
    if (this.isAdd) {
      cluster = <CreateCluster>{
        "name": this.clusterForm.controls["name"].value,
      };
      promise = this.client.create(cluster);
    } else {
      cluster = <UpdateCluster>{
        "uuid": this.cluster.uuid,
        "name": this.clusterForm.controls["name"].value,
      };
      promise = this.client.update(cluster);
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
