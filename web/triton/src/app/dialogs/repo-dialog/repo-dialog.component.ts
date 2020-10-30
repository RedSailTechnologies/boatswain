import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import * as fetch from 'isomorphic-fetch';
import { DefaultPoseidon, Poseidon, Repo } from 'src/app/services/poseidon/poseidon';
import { BusyComponent } from '../busy/busy.component';
import { MessageDialogComponent } from '../message-dialog/message-dialog.component';

@Component({
  selector: 'app-repo-dialog',
  templateUrl: './repo-dialog.component.html',
  styleUrls: ['./repo-dialog.component.sass']
})
export class RepoDialogComponent implements OnInit {
  private client: Poseidon;
  private repo: Repo;
  public repoForm: FormGroup = new FormGroup({
    name: new FormControl(''),
    endpoint: new FormControl(''),
  });
  public isAdd: boolean;
  public title: string;

  constructor(public dialog: MatDialogRef<RepoDialogComponent>, 
              @Inject(MAT_DIALOG_DATA) data,
              private spinner: MatDialog,
              private error: MatDialog) {
    this.title = data["title"];
    this.isAdd = data["type"] == "add";
    if (!this.isAdd) {
      this.repo = data["repo"];
      this.repoForm.controls["name"].setValue(this.repo.name);
      this.repoForm.controls["endpoint"].setValue(this.repo.endpoint);
    }
    this.client = new DefaultPoseidon(`${location.protocol}//${location.host}/api`, fetch["default"]);
  }

  ngOnInit(): void {
  }

  enter($event): void {
    if ($event.keyCode == 13 && this.repoForm.valid) {
      this.submit()
    };
  }

  submit(): void {
    var spinnerRef: MatDialogRef<BusyComponent> = this.spinner.open(BusyComponent, {
      panelClass: 'transparent',
      disableClose: true
    });
    var cluster = <Repo>{
      "uuid": this.repo != null ? this.repo.uuid : null,
      "name": this.repoForm.controls["name"].value,
      "endpoint": this.repoForm.controls["endpoint"].value,
      "ready": false
    };

    var promise: Promise<any>
    if (this.isAdd) {
      promise = this.client.addRepo(cluster);
    } else {
      promise = this.client.editRepo(cluster);
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
}
