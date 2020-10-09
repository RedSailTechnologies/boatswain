import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { MatSort } from '@angular/material/sort';
import { Cluster, Deployment } from 'src/app/services/kraken/service_pb';
import { KrakenService } from 'src/app/services/kraken/kraken.service'

@Component({
  selector: 'app-clusters',
  templateUrl: './clusters.component.html',
  styleUrls: ['./clusters.component.sass']
})

export class ClustersComponent implements AfterViewInit, OnInit {
  clusters: Cluster.AsObject[];
  deploymentData: Map<Cluster.AsObject, Deployment.AsObject[]>;

  constructor(private kraken: KrakenService) {  }

  @ViewChild(MatSort) sort: MatSort;

  ngOnInit() {
    this.deploymentData = new Map<Cluster.AsObject, Deployment.AsObject[]>();
    this.kraken.getClusters().then(value => {
      this.clusters = value.clustersList;
      this.clusters.forEach(cluster => {
        this.kraken.getDeployments(cluster.name).then(value => {
          this.deploymentData[cluster.name] = value.deploymentsList;
        });
      });
    });
  }
  
  ngAfterViewInit() {
    // this.dataSource.sort = this.sort;
  }
}