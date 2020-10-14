import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { DefaultKraken, Kraken, Release } from 'src/app/services/kraken/kraken';
import { Chart, ChartVersion, DefaultPoseidon, Poseidon, Repo } from 'src/app/services/poseidon/poseidon';
import { FormControl, FormGroup } from '@angular/forms';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-update-dialog',
  templateUrl: './update-dialog.component.html',
  styleUrls: ['./update-dialog.component.sass']
})
export class UpdateDialogComponent implements OnInit {
  private kraken: Kraken;
  private poseidon: Poseidon;

  public name: string;
  public chart: string;
  public release: Release;

  public repos: Repo[];
  public charts: Chart[];
  public versions: ChartVersion[];

  public upgradeForm: FormGroup = new FormGroup({
    repo: new FormControl(''),
    chart: new FormControl(''),
    chartVersion: new FormControl(''),
    appVersion: new FormControl(''),
    additionalValues: new FormControl('')
  });

  constructor(public dialog: MatDialogRef<UpdateDialogComponent>, @Inject(MAT_DIALOG_DATA) data) {
    this.name = data['name'];
    this.chart = data['chart'];
    this.release = data['release'];
    this.repos = data['repos'];
    this.kraken = new DefaultKraken(`${location.protocol}//${location.host}/api`, fetch['default']);
    this.poseidon = new DefaultPoseidon(`${location.protocol}//${location.host}/api`, fetch['default']);
  }

  ngOnInit() : void {
    this.poseidon.repos({}).then(value => {
      this.repos = value.repos;
      for (let i = 0; i < this.repos.length; i++) {
        if (this.repos[i].ready) {
          this.upgradeForm.controls['repo'].setValue(this.repos[i].name)
        }
      }
    });

    this.upgradeForm.controls['chart'].disable();
    this.upgradeForm.controls['chartVersion'].disable();
    this.upgradeForm.controls['appVersion'].disable();
    this.upgradeForm.controls['additionalValues'].disable();
  }

  chartUpdated(chart: Chart) {
    this.upgradeForm.controls['chartVersion'].enable();
    this.versions = chart.versions;
    this.upgradeForm.controls['chartVersion'].setValue(this.versions[0].chartVersion);
  }

  chartVersionUpdated() {
    this.upgradeForm.controls['appVersion'].enable();
    this.upgradeForm.controls['additionalValues'].enable();
  }

  repoUpdated(repo: Repo) {
    this.poseidon.charts(repo).then(value => {
      this.charts = value.charts;
      this.upgradeForm.controls['chart'].setValue(this.chart);
      this.upgradeForm.controls['chart'].enable();
    })
  }
}
