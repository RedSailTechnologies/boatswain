import { Component, OnInit } from '@angular/core';
import { Chart, DefaultPoseidon, Poseidon, Repo } from 'src/app/services/poseidon/poseidon';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-repos',
  templateUrl: './repos.component.html',
  styleUrls: ['./repos.component.sass']
})
export class ReposComponent implements OnInit {
  private client: Poseidon;
  public repos: Repo[];
  public charts: Map<Repo, Chart[]>;

  constructor() {
    this.client = new DefaultPoseidon(`${location.protocol}//${location.host}/api`, fetch['default']);
  }

  ngOnInit(): void {
    this.charts = new Map<Repo, Chart[]>();
    this.client.repos({}).then(value => {
      this.repos = value.repos;
      value.repos.forEach(repo => {
        if (repo.ready) {
          this.client.charts(repo).then(results => {
            this.charts.set(repo, results.charts);
          });
        }
      });
    });
  }

}
