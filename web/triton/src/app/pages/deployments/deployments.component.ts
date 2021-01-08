import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
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
          title: 'Add Repo',
        },
      })
      .afterClosed()
      .subscribe((_) => {
        this.refreshDeployments();
      });
  }

  edit(repo: DeploymentReadSummary): void {
    this.dialog
      .open(DeploymentDialogComponent, {
        minWidth: '33%',
        panelClass: 'custom-dialog-container',
        data: {
          type: 'edit',
          title: `Edit ${repo.name}`,
          repo: repo,
        },
      })
      .afterClosed()
      .subscribe((_) => {
        this.refreshDeployments();
      });
  }

  delete(repo: DeploymentReadSummary): void {
    this.dialog
      .open(DeploymentDialogComponent, {
        panelClass: 'message-box',
        data: {
          reason: `Delete ${repo.name}`,
          message: 'Do you really want to delete this repo?',
        },
      })
      .afterClosed()
      .subscribe((result: Boolean) => {
        if (result) {
          this.client
            .destroy(repo)
            .catch((_) => {
              this.snackBar.open(
                `${repo.name} could not be deleted`,
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
      console.log('could not update repos after 5 retries');
      this.retries = 0;
      this.deployments = new Array<DeploymentReadSummary>();
      this.snackBar.open(`Error getting repos`, 'Dismiss', {
        duration: 5 * 1000,
        panelClass: ['warn-snack'],
      });
    }
  }
}
