import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { MatSort } from '@angular/material/sort';
import { Cluster, DefaultKraken, Deployment, Kraken } from 'src/app/services/kraken/service';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-clusters',
  templateUrl: './clusters.component.html',
  styleUrls: ['./clusters.component.sass']
})

export class ClustersComponent implements AfterViewInit, OnInit {
  private client: Kraken;
  public clusters: Cluster[];
  public deploymentData: Map<Cluster, Deployment[]>;

  constructor() {
    this.client = new DefaultKraken(`${location.protocol}//${location.host}/api`, fetch['default']);
  }

  @ViewChild(MatSort) sort: MatSort;

  ngOnInit() {
    this.deploymentData = new Map<Cluster, Deployment[]>();
    this.client.clusters({}).then(value => {
      this.clusters = value.clusters;
      this.clusters.forEach(cluster => {
        this.client.deployments({"cluster": cluster}).then(value => {
          this.deploymentData.set(cluster, value.deployments);
        });
      });
    });
  }
  
  ngAfterViewInit() {
    // this.dataSource.sort = this.sort;
  }
}
