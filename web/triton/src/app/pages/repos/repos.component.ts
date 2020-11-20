import { Component, OnInit } from '@angular/core';
import { RepoRead, ChartRead, Repo, DefaultRepo } from 'src/app/services/repo/repo';
import * as fetch from 'isomorphic-fetch';
import { RepoDialogComponent } from 'src/app/dialogs/repo-dialog/repo-dialog.component';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';
import { ConfirmDialogComponent } from 'src/app/dialogs/confirm-dialog/confirm-dialog.component';

@Component({
  selector: 'app-repos',
  templateUrl: './repos.component.html',
  styleUrls: ['./repos.component.sass']
})
export class ReposComponent implements OnInit {
  private client: Repo;
  private retries = 0;
  public repos: RepoRead[];
  public charts: Map<RepoRead, ChartRead[]>;
  public expandedRepo: RepoRead;
  public expandedCharts: ChartRead[];

  
  constructor(private dialog: MatDialog, private snackBar: MatSnackBar) {
    this.client = new DefaultRepo(`${location.protocol}//${location.host}/api`, fetch['default']);
  }
  
  ngOnInit(): void {
    this.refreshRepos()
  }
  
  add() : void {
    this.dialog.open(RepoDialogComponent, {
      minWidth: "33%",
      panelClass: "custom-dialog-container",
      data: {
        "type": "add",
        "title": "Add Repo"
      }
    }).afterClosed().subscribe(_ => {
      this.refreshRepos();
    });
  }

  edit(repo: RepoRead) : void {
    this.dialog.open(RepoDialogComponent, {
      minWidth: "33%",
      panelClass: "custom-dialog-container",
      data: {
        "type": "edit",
        "title": `Edit ${repo.name}`,
        "repo": repo
      }
    }).afterClosed().subscribe(_ => {
      this.refreshRepos();
    });
  }

  delete(repo: RepoRead) : void {
    this.dialog.open(ConfirmDialogComponent, {
      panelClass: 'message-box',
      data: {
        "reason": `Delete ${repo.name}`,
        "message": "Do you really want to delete this repo?"
      }
    }).afterClosed().subscribe((result: Boolean) => {
      if (result) {
        this.client.destroy(repo).catch(_ => {
          this.snackBar.open(`${repo.name} could not be deleted`, "Dismiss", {
            duration: 5 * 1000,
            panelClass: ["warn-snack"]
          })
        }).then(() => {
          this.refreshRepos()
        });
      }
    });
  }

  refreshRepos() : void {
    if (this.retries < 5) {
      this.retries++;
      this.charts = new Map<RepoRead, ChartRead[]>();
      this.client.all({}).then(value => {
        this.repos = value.repos;
        value.repos.forEach(repo => {
          if (repo.ready) {
            this.client.charts(repo).then(results => {
              this.charts.set(repo, results.charts);
            }).catch(_ => {
              setTimeout(() => this.refreshRepos(), 2 * 1000)
            });
          }
        });
      }).catch(_ => {
        setTimeout(() => this.refreshRepos(), 2 * 1000)
      });
    } else {
      console.log("could not update repos after 5 retries");
      this.retries = 0;
      this.repos = new Array<RepoRead>();
      this.snackBar.open(`Error getting repos`, "Dismiss", {
        duration: 5 * 1000,
        panelClass: ["warn-snack"]
      });
    }
  }
}
