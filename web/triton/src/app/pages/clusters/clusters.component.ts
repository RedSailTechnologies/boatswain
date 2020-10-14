import { Component, OnInit } from '@angular/core';
import { Cluster, DefaultKraken, Kraken } from 'src/app/services/kraken/kraken';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-clusters',
  templateUrl: './clusters.component.html',
  styleUrls: ['./clusters.component.sass']
})

export class ClustersComponent implements OnInit {
  private client: Kraken;
  public clusters: Cluster[];

  constructor() {
    this.client = new DefaultKraken(`${location.protocol}//${location.host}/api`, fetch['default']);
  }

  ngOnInit() : void {
    this.client.clusters({}).then(value => {
      this.clusters = value.clusters;
    });
  }
}
