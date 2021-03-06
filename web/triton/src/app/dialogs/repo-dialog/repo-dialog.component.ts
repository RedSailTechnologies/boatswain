import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { RepoRead, Repo, DefaultRepo, UpdateRepo, CreateRepo } from 'src/app/services/repo/repo';
import { BusyComponent } from '../busy/busy.component';
import { MessageDialogComponent } from '../message-dialog/message-dialog.component';
import { AuthService } from 'src/app/utils/auth/auth.service';

@Component({
  selector: 'app-repo-dialog',
  templateUrl: './repo-dialog.component.html',
  styleUrls: ['./repo-dialog.component.sass']
})
export class RepoDialogComponent implements OnInit {
  private client: Repo;
  private repo: RepoRead;
  
  public repoTypes: string[] = ["HELM", "GIT"];
  public repoForm: FormGroup = new FormGroup({
    name: new FormControl(''),
    endpoint: new FormControl(''),
    token: new FormControl(''),
    username: new FormControl(''),
    password: new FormControl(''),
    type: new FormControl(''),
    oci: new FormControl(''),
  });
  public isAdd: boolean;
  public title: string;
  public authType: string;

  constructor(public dialog: MatDialogRef<RepoDialogComponent>, 
              @Inject(MAT_DIALOG_DATA) data,
              private spinner: MatDialog,
              private error: MatDialog,
              auth: AuthService) {
    this.title = data["title"];
    this.isAdd = data["type"] == "add";
    if (!this.isAdd) {
      this.repo = data["repo"];
      this.repoForm.controls["name"].setValue(this.repo.name);
      this.repoForm.controls["endpoint"].setValue(this.repo.endpoint);
      this.repoForm.controls["token"].setValue("***");
      this.repoForm.controls["username"].setValue("***");
      this.repoForm.controls["password"].setValue("***");
      this.repoForm.controls["type"].setValue(this.repo.type);
      this.repoForm.controls["oci"].setValue(this.repo.helmOci);
    }
    this.client = new DefaultRepo(`${location.protocol}//${location.host}/api`, auth.fetch());
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
    

    var promise: Promise<any>
    var repo
    if (this.isAdd) {
      repo = <CreateRepo>{
        "name": this.repoForm.controls["name"].value,
        "endpoint": this.repoForm.controls["endpoint"].value,
        "token": this.repoForm.controls["token"].value == "***" ? "" : this.repoForm.controls["token"].value,
        "username": this.repoForm.controls["username"].value == "***" ? "" : this.repoForm.controls["username"].value,
        "password": this.repoForm.controls["password"].value == "***" ? "" : this.repoForm.controls["password"].value,
        "type": <string><unknown>this.typeEnum(),
        "helmOci": this.repoForm.controls["oci"].value != "" ? this.repoForm.controls["oci"].value : false
      };
      promise = this.client.create(repo);
    } else {
      repo = <UpdateRepo>{
        "uuid": this.repo.uuid,
        "name": this.repoForm.controls["name"].value,
        "endpoint": this.repoForm.controls["endpoint"].value,
        "username": this.repoForm.controls["username"].value == "***" ? "" : this.repoForm.controls["username"].value,
        "password": this.repoForm.controls["password"].value == "***" ? "" : this.repoForm.controls["password"].value,
        "type": <string><unknown>this.typeEnum(),
        "helmOci": this.repoForm.controls["oci"].value != "" ? this.repoForm.controls["oci"].value : false
      };
      promise = this.client.update(repo);
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

  typeEnum() : number {
    switch (this.repoForm.controls["type"].value) {
      case "HELM":
        return 0;
      case "GIT":
        return 1;
    }
  }
}
