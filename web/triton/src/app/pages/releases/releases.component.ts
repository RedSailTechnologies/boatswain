import { Component, OnInit } from '@angular/core';
import { DefaultKraken, Kraken, Releases } from 'src/app/services/kraken/service';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-releases',
  templateUrl: './releases.component.html',
  styleUrls: ['./releases.component.sass']
})
export class ReleasesComponent implements OnInit {
  private client: Kraken;
  public releasesList: Releases[];

  constructor() {
    this.client = new DefaultKraken(`${location.protocol}//${location.host}/api`, fetch['default']);
  }

  ngOnInit(): void {
    this.client.clusters({}).then(clusters => {
      this.client.releases({"clusters": clusters.clusters}).then(releases => {
        this.releasesList = releases.releaseLists;
      })
    });
  }
}
