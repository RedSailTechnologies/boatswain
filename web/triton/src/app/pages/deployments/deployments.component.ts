import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router'; 
import { ConfirmDialogComponent } from 'src/app/dialogs/confirm-dialog/confirm-dialog.component';
import { DeploymentDialogComponent } from 'src/app/dialogs/deployment-dialog/deployment-dialog.component';
import { DefaultDeployment, Deployment, DeploymentReadSummary } from 'src/app/services/deployment/deployment';
import { AuthService } from 'src/app/utils/auth/auth.service';

@Component({
  selector: 'app-deployments',
  templateUrl: './deployments.component.html',
  styleUrls: ['./deployments.component.sass']
})
export class DeploymentsComponent implements OnInit {
  private client: Deployment;
  private retries = 0;
  public deployments: DeploymentReadSummary[];

  constructor(
    private dialog: MatDialog,
    private snackBar: MatSnackBar,
    private router: Router,
    public auth: AuthService
  ) {
    this.client = new DefaultDeployment(
      `${location.protocol}//${location.host}/api`,
      auth.fetch()
    );
  }

  ngOnInit(): void {
    this.refreshDeployments();
  }

  add(): void {
    this.dialog
      .open(DeploymentDialogComponent, { // create
        minWidth: '33%',
        panelClass: 'custom-dialog-container',
        data: {
          type: 'add',
          title: 'Add Deployment',
        },
      })
      .afterClosed()
      .subscribe((_) => {
        this.refreshDeployments();
      });
  }

  edit(deployment: DeploymentReadSummary): void {
    this.dialog
      .open(DeploymentDialogComponent, {
        minWidth: '33%',
        panelClass: 'custom-dialog-container',
        data: {
          type: 'edit',
          title: `Edit ${deployment.name}`,
          uuid: deployment.uuid,
        },
      })
      .afterClosed()
      .subscribe((_) => {
        this.refreshDeployments();
      });
  }

  delete(deployment: DeploymentReadSummary): void {
    this.dialog
      .open(ConfirmDialogComponent, {
        panelClass: 'message-box',
        data: {
          reason: `Delete ${deployment.name}`,
          message: 'Do you really want to delete this deployment?',
        },
      })
      .afterClosed()
      .subscribe((result: Boolean) => {
        if (result) {
          this.client
            .destroy(deployment)
            .catch((_) => {
              this.snackBar.open(
                `${deployment.name} could not be deleted`,
                'Dismiss',
                {
                  duration: 5 * 1000,
                  panelClass: ['warn-snack'],
                }
              );
            })
            .then(() => {
              this.refreshDeployments();
            });
        }
      });
  }

  redirect(deployment: DeploymentReadSummary): void {
    this.router.navigate(['/deployment/' + deployment.uuid])
  }

  refreshDeployments(): void {
    if (this.retries < 5) {
      this.retries++;
      this.client
        .all({})
        .then((value) => {
          this.deployments = value.deployments;
        })
        .catch((_) => {
          setTimeout(() => this.refreshDeployments(), 2 * 1000);
        });
    } else {
      console.log('could not update deployments after 5 retries');
      this.retries = 0;
      this.deployments = new Array<DeploymentReadSummary>();
      this.snackBar.open(`Error getting deployments`, 'Dismiss', {
        duration: 5 * 1000,
        panelClass: ['warn-snack'],
      });
    }
  }
}
