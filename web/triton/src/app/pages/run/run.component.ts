import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { DefaultDeployment, Deployment, ReadRun, RunRead } from 'src/app/services/deployment/deployment';
import { AuthService } from 'src/app/utils/auth/auth.service';

@Component({
  selector: 'app-run',
  templateUrl: './run.component.html',
  styleUrls: ['./run.component.sass']
})
export class RunComponent implements OnInit {
  private client: Deployment;
  public run: RunRead;

  constructor(private route: ActivatedRoute, auth: AuthService) {
    this.client = new DefaultDeployment(
      `${location.protocol}//${location.host}/api`,
      auth.fetch()
    );
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      var id = params['uuid']
      this.client.run(<ReadRun>{
        deploymentUuid: id
      }).then(value => {
        this.run = value;
        // FIXME set the page title and browser title somehow
      })
    })
  }
}
