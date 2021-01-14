import { Component, OnInit } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { BusyComponent } from 'src/app/dialogs/busy/busy.component';
import { DefaultDeployment, Deployment, DeploymentRead, ReadDeployment, ReadRuns, RunReadSummary, TemplateDeployment, TriggerDeployment } from 'src/app/services/deployment/deployment';
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
              private spinner: MatDialog,
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
          this.runs = value.runs;
        })
      });
    })
  }

  run() {
    var runTrigger = this.client.trigger(<TriggerDeployment>{
      uuid: this.deployment.uuid,
      type: "MANUAL",
      // TODO - args
    });
    var spinnerRef: MatDialogRef<BusyComponent> = this.spinner.open(BusyComponent, {
      panelClass: 'transparent',
      disableClose: true
    });
    runTrigger.finally(() => {
      spinnerRef.close();
    });
  }

  redirect(run: RunReadSummary) {
    this.router.navigate(['/run/' + run.uuid]);
  }
}
