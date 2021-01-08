import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { CreateDeployment, DefaultDeployment, Deployment, DeploymentRead, UpdateDeployment } from 'src/app/services/deployment/deployment';
import { DefaultRepo, Repo, RepoRead } from 'src/app/services/repo/repo';
import { AuthService } from 'src/app/utils/auth/auth.service';
import { BusyComponent } from '../busy/busy.component';
import { MessageDialogComponent } from '../message-dialog/message-dialog.component';

@Component({
  selector: 'app-deployment-dialog',
  templateUrl: './deployment-dialog.component.html',
  styleUrls: ['./deployment-dialog.component.sass']
})
export class DeploymentDialogComponent implements OnInit {
  private client: Deployment;
  private deployment: DeploymentRead;

  public repos: RepoRead[];
  public deploymentForm: FormGroup = new FormGroup({
    name: new FormControl(''),
    repoId: new FormControl(''),
    branch: new FormControl(''),
    filePath: new FormControl(''),
  });
  public isAdd: boolean;
  public title: string;

  constructor(public dialog: MatDialogRef<DeploymentDialogComponent>, 
              @Inject(MAT_DIALOG_DATA) data,
              private spinner: MatDialog,
              private error: MatDialog,
              auth: AuthService) {
    this.title = data["title"];
    this.isAdd = data["type"] == "add";
    if (!this.isAdd) {
      this.deployment = data["deployment"];
      this.deploymentForm.controls["name"].setValue(this.deployment.name);
      this.deploymentForm.controls["repoId"].setValue(this.deployment.repoId);
      this.deploymentForm.controls["branch"].setValue(this.deployment.branch);
      this.deploymentForm.controls["filePath"].setValue(this.deployment.filePath);
    }
    this.client = new DefaultDeployment(`${location.protocol}//${location.host}/api`, auth.fetch());
    
    var repoClient = new DefaultRepo(`${location.protocol}//${location.host}/api`, auth.fetch());
    repoClient.all({}).then(value => {
      this.repos = value.repos;
    })
  }

  ngOnInit(): void {
  }

  enter($event): void {
    if ($event.keyCode == 13 && this.deploymentForm.valid) {
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
      repo = <CreateDeployment>{
        "uuid": this.deployment != null ? this.deployment.uuid : null,
        "name": this.deploymentForm.controls["name"].value,
        "repoId": this.deploymentForm.controls["repoId"].value,
        "branch": this.deploymentForm.controls["branch"].value,
        "filePath": this.deploymentForm.controls["filePath"].value,
      };
      promise = this.client.create(repo);
    } else {
      repo = <UpdateDeployment>{
        "uuid": this.deployment != null ? this.deployment.uuid : null,
        "name": this.deploymentForm.controls["name"].value,
        "repoId": this.deploymentForm.controls["repoId"].value,
        "branch": this.deploymentForm.controls["branch"].value,
        "filePath": this.deploymentForm.controls["filePath"].value,
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
}
