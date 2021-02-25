import { STEPPER_GLOBAL_OPTIONS } from '@angular/cdk/stepper';
import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { ConfirmDialogComponent } from 'src/app/dialogs/confirm-dialog/confirm-dialog.component';
import { ApproveStep, DefaultDeployment, Deployment, ReadRun, RunRead, StepRead } from 'src/app/services/deployment/deployment';
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

  constructor(private dialog: MatDialog, private route: ActivatedRoute, public auth: AuthService) {
    this.client = new DefaultDeployment(
      `${location.protocol}//${location.host}/api`,
      auth.fetch()
    );
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      this.id = params['uuid']
      this.refresh();
    }).unsubscribe()
  }

  refresh(): void {
    this.client.run(<ReadRun>{
      deploymentUuid: this.id
    }).then(value => {
      this.run = value;
      this.start = this.formatDate(this.run.startTime);
      this.stop = this.formatDate(this.run.stopTime);

      if (this.run.status == "IN_PROGRESS") {
        setTimeout(() => this.refresh(), 3 * 1000);
      }
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

  approval(step: StepRead, approve: boolean, override: boolean) {
    var message = "Approve";
    if (!approve) {
      message = "Reject";
    }
    if (override) {
      message = "Override";
    }
    this.dialog.open(ConfirmDialogComponent, {
      panelClass: 'message-box',
      data: {
        reason: message,
        message: `${message} ${step.name}?`,
      },
    })
    .afterClosed()
    .subscribe((result: boolean) => {
      if (result) {
        this.client.approve(<ApproveStep>{
          runUuid: this.run.uuid,
          approve: approve,
          override: override
        });
      }
    });
  }
}
