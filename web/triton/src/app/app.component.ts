import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { AuthService } from './utils/auth/auth.service';
import { Profile } from 'oidc-client';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass']
})
export class AppComponent {
  private auth: AuthService;
  public page = "";
  public sideNav = false;
  public userIcon: string;
  public userProfile: Profile;

  constructor(router: Router, title: Title, auth: AuthService) {
      router.events.subscribe(_ => {
        this.page = this.getPage(router.url);
        title.setTitle('Boatswain - ' + this.getPage(router.url));
      });

      this.auth = auth;
      auth.events.subscribe(_ => {
        this.updateAuth();
      })
      this.updateAuth();
  }

  public clickLoginButton() {
    if (this.auth.loggedIn()) {
      this.logout();
    } else {
      this.login();
    }
  }
  
  public toggleSideNav() {
    this.sideNav = !this.sideNav;
  }
  
  private getPage(url: string) {
    if (url == '/') return 'Home';
    return url[1].toUpperCase() + url.slice(2);
  }
  
  private login() {
    this.auth.startLogin();
  }

  private logout() {
    this.auth.startLogout();
  }

  private updateAuth() {
    this.userIcon = this.auth.loggedIn() ? 'logout' : 'login';
    if (this.auth.loggedIn()) {
      this.userProfile = this.auth.getProfile();
    }
  }
}
