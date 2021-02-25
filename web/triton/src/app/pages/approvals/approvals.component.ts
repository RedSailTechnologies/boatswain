import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ConfirmDialogComponent } from 'src/app/dialogs/confirm-dialog/confirm-dialog.component';
import { ApprovalRead, ApprovalsRead, ApproveStep, DefaultDeployment, Deployment } from 'src/app/services/deployment/deployment';
import { TwirpError } from 'src/app/services/deployment/twirp';
import { AuthService } from 'src/app/utils/auth/auth.service';

@Component({
  selector: 'app-approvals',
  templateUrl: './approvals.component.html',
  styleUrls: ['./approvals.component.sass']
})
export class ApprovalsComponent implements OnInit {
  private client: Deployment;
  private retries = 0;

  public approvals: ApprovalRead[];

  constructor(
    private dialog: MatDialog,
    private snackBar: MatSnackBar,
    public auth: AuthService
  ) {
    this.client = new DefaultDeployment(
      `${location.protocol}//${location.host}/api`,
      auth.fetch()
    );
  }

  ngOnInit(): void {
    this.refresh()
  }

  approval(approval: ApprovalRead, approve: boolean, override: boolean) {
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
        message: `${message} ${approval.stepName}?`,
      },
    })
    .afterClosed()
    .subscribe((result: boolean) => {
      if (result) {
        this.client.approve(<ApproveStep>{
          runUuid: approval.runUuid,
          approve: approve,
          override: override
        }).finally(() => {
          this.refresh();
        });
      }
    });
  }

  refresh() {
    if (this.retries < 5) {
      this.retries++;
      this.client.approvals({})
      .then((value: ApprovalsRead) => {
        this.approvals = value.approvals;
        this.retries = 0;
      }).catch((err: TwirpError) => {
        if (err.code == 'Unauthorized') {
          this.snackBar.open(`Unauthorized`, 'Dismiss', {
            duration: 5 * 1000,
            panelClass: ['warn-snack'],
          });
        } else {
          // setTimeout(() => this.refreshClusters(), 2 * 1000);
        }
      });
    } else {
      console.log('could not update approvals after 5 retries');
      this.retries = 0;
      this.approvals = new Array<ApprovalRead>();
      this.snackBar.open(`Error getting approvals`, 'Dismiss', {
        duration: 5 * 1000,
        panelClass: ['warn-snack'],
      });
    }
  }
}
