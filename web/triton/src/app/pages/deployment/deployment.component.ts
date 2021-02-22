import { Component, OnInit } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { BusyComponent } from 'src/app/dialogs/busy/busy.component';
import { MessageDialogComponent } from 'src/app/dialogs/message-dialog/message-dialog.component';
import { TriggerDialogComponent } from 'src/app/dialogs/trigger-dialog/trigger-dialog.component';
import { DefaultDeployment, Deployment, DeploymentRead, ReadDeployment, ReadRuns, ReadToken, RunReadSummary, TemplateDeployment } from 'src/app/services/deployment/deployment';
import { DefaultTrigger, ManualTriggered, TriggerManual } from 'src/app/services/trigger/trigger';
import { AuthService } from 'src/app/utils/auth/auth.service';

@Component({
  selector: 'app-deployment',
  templateUrl: './deployment.component.html',
  styleUrls: ['./deployment.component.sass']
})
export class DeploymentComponent implements OnInit {
  private client: Deployment;
  public deployment: DeploymentRead;
  public yaml: string;
  public templateError: boolean;
  public runs: RunReadSummary[];

  constructor(private route: ActivatedRoute,
              private router: Router,
              private dialog: MatDialog,
              auth: AuthService) {
    this.client = new DefaultDeployment(`${location.protocol}//${location.host}/api`, auth.fetch());
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      var id = params['uuid'];
      this.client.read(<ReadDeployment>{
        uuid: id,
      }).then(value => {
        this.deployment = value;

        // FIXME set the page title and browser title somehow
        // title.setTitle('Boatswain - ' + this.getPage(router.url));
        this.client.template(<TemplateDeployment>{
          uuid: this.deployment.uuid,
        }).then(value => {
          this.templateError = false;
          this.yaml = value.yaml;
        }).catch(reason => {
          this.templateError = true;
          this.yaml = reason;
        });

        this.client.runs(<ReadRuns>{
          deploymentUuid: id
        }).then(value => {
          value.runs = value.runs.sort((a, b) => a.startTime > b.startTime ? -1 : 1)
          this.runs = value.runs;
          this.runs.forEach(x => {
            if (x.startTime != 0) {
              x['startFormatted'] = new Date(x.startTime * 1000).toLocaleString();
            }
            if (x.stopTime != 0) {
              x['stopFormatted'] = new Date(x.stopTime * 1000).toLocaleString();
            }
          });
        })
      });
    })
  }

  run() {
    this.dialog.open(TriggerDialogComponent, {
      minWidth: '33%',
      panelClass: 'custom-dialog-container',
      data: {
        title: "Trigger Deployment",
        uuid: this.deployment.uuid
      }
    })
    .afterClosed()
    .subscribe(val => {
      this.router.navigate(['/run/' + val]);
    });
  }

  redirect(run: RunReadSummary) {
    this.router.navigate(['/run/' + run.uuid]);
  }

  token() {
    this.client.token(<ReadToken>{
      uuid: this.deployment.uuid
    }).then(val => {
      this.dialog.open(MessageDialogComponent, {
        panelClass: 'message-box',
        data: {
          "reason": "Token",
          "message": "Deployment UUID: " + this.deployment.uuid + "\nDeployment Token: " + val.token
        }
      });
    }).catch(error => {
      this.dialog.open(MessageDialogComponent, {
        panelClass: 'message-box',
        data: {
          "reason": "Error",
          "message": "An error occured.\n" + error
        }
      });
    })
  }
}
