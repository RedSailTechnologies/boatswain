import { Component, OnInit } from '@angular/core';
import { DefaultPoseidon, Poseidon, Repo } from 'src/app/services/poseidon/service';
import * as fetch from 'isomorphic-fetch';

@Component({
  selector: 'app-repos',
  templateUrl: './repos.component.html',
  styleUrls: ['./repos.component.sass']
})
export class ReposComponent implements OnInit {
  private client: Poseidon;
  public repos: Repo[];

  constructor() {
    this.client = new DefaultPoseidon(`${location.protocol}//${location.host}/api`, fetch['default']);
  }

  ngOnInit(): void {
    this.client.repos({}).then(value => {
      this.repos = value.repos;
    });
  }

}
