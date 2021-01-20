import { STEPPER_GLOBAL_OPTIONS } from '@angular/cdk/stepper';
import { Component, OnInit, ViewChild } from '@angular/core';
import { MatStepper } from '@angular/material/stepper';
import { ActivatedRoute } from '@angular/router';
import { DefaultDeployment, Deployment, ReadRun, RunRead } from 'src/app/services/deployment/deployment';
import { AuthService } from 'src/app/utils/auth/auth.service';

@Component({
  selector: 'app-run',
  templateUrl: './run.component.html',
  styleUrls: ['./run.component.sass'],
  providers: [{
    provide: STEPPER_GLOBAL_OPTIONS, useValue: {showError: true, displayDefaultIndicatorType: false,}
  }]
})
export class RunComponent implements OnInit {
  private client: Deployment;
  private id: string;
  public run: RunRead;
  public start: string;
  public stop: string;
  public idx: number;

  @ViewChild('stepper') stepper : MatStepper;

  constructor(private route: ActivatedRoute, auth: AuthService) {
    this.client = new DefaultDeployment(
      `${location.protocol}//${location.host}/api`,
      auth.fetch()
    );
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      this.id = params['uuid']
      this.refresh();
    })
  }

  refresh(): void {
    this.client.run(<ReadRun>{
      deploymentUuid: this.id
    }).then(value => {
      this.run = value;
      this.start = new Date(this.run.startTime * 1000).toLocaleString();
      this.stop = new Date(this.run.stopTime * 1000).toLocaleString();
      if (this.run.status == "IN_PROGRESS") {
        setTimeout(() => this.refresh(), 3 * 1000);
      }

      if (this.stepper) {
        this.stepper.next()
      }
      // FIXME set the page title and browser title somehow
    });
  }

  formatDate(date: number): string {
    if (date == 0) {
      return "";
    }
    return new Date(date * 1000).toLocaleString();
  }

  statusColor(status: string): string {
    switch (status) {
      case 'FAILED':
        return 'color:lightcoral;'
      case 'SKIPPED':
        return 'color:grey;'
      default:
        return ''
    }
  }

  stepError(status: string): boolean {
    return status == 'FAILED'
  }
}
