import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Title } from '@angular/platform-browser';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass']
})
export class AppComponent {
  public page = "";
  public sideNav = false;

  constructor(router: Router, title: Title) {
      router.events.subscribe(_ => {
        this.page = this.getPage(router.url);
        title.setTitle('Boatswain - ' + this.getPage(router.url));
      });
  }

  public toggleSideNav() {
    this.sideNav = !this.sideNav;
  }

  private getPage(url: string) {
    if (url == "/") return "Home";
    return url[1].toUpperCase() + url.slice(2);
  }
}
