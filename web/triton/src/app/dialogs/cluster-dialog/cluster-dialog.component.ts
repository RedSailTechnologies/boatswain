import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Cluster } from 'src/app/services/kraken/kraken';

@Component({
  selector: 'app-cluster-dialog',
  templateUrl: './cluster-dialog.component.html',
  styleUrls: ['./cluster-dialog.component.sass']
})
export class ClusterDialogComponent implements OnInit {
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
      var cluster: Cluster = data["cluster"];
      this.clusterForm.controls["name"].setValue(cluster.name);
      this.clusterForm.controls["endpoint"].setValue(cluster.endpoint);
    }
  }

  ngOnInit(): void {
  }

  submit(): void {
  }
}
