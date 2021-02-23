import { EventEmitter, Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Profile, User, UserManager, UserManagerSettings } from 'oidc-client';
import * as fetch from 'isomorphic-fetch';
import { ConfigService } from '../config/config.service';

const actionKey: string = "loginActionInProgress";
const redirectKey: string = "loginRedirectDestination";
const userKey: string = "currentUser";

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private static settings: UserManagerSettings = null;
  private mgr: UserManager;
  private user: User;

  public static IN_PROGRESS: boolean = (sessionStorage.getItem(actionKey) ?? "false") != "false";
  public events: EventEmitter<User> = new EventEmitter<User>();

  constructor(private router: Router, configService: ConfigService) {
    if (AuthService.settings == null) {
      const config = configService.getOIDC();
      AuthService.settings = {
        authority: config.authority,
        client_id: config.clientId,
        scope: `${config.scope}`,
      
        redirect_uri: `${window.location.origin}/login`,
        post_logout_redirect_uri: `${window.location.origin}/logout`,
        response_type: 'id_token token',
        filterProtocolClaims: true,
        loadUserInfo: false
      };
    }
    this.mgr = new UserManager(AuthService.settings);

    var storedString = sessionStorage.getItem(userKey);
    if (storedString != null) {
      var storedUser: User = User.fromStorageString(storedString);
      if (!storedUser.expired) {
        this.user = storedUser;
        this.mgr.startSilentRenew();
      } else {
        sessionStorage.removeItem(userKey);
      }
    }
  }

  authHeader(): string {
    if (this.user == null) {
      return null;
    }
    return `${this.user.token_type} ${this.user.access_token}`;
  }

  fetch() {
    var header = this.authHeader();
    return function (
      input: RequestInfo,
      init?: RequestInit
    ): Promise<Response> {
      (<Request>input).headers.append("Authorization", header);
      return fetch['default'](input, init);
    }
  }

  getProfile(): Profile {
    if (this.user == null) {
      return null;
    }
    return this.user.profile;
  }

  isAdmin(): boolean {
    if (this.user == null) {
      return false;
    }
    for (var i = 0; i < this.user.profile.roles.length; i++) {
      if (this.user.profile.roles[i] == "Boatswain.Admin") {
        return true;
      }
    }
    return false;
  }

  isEditor(): boolean {
    if (this.user == null) {
      return false;
    }
    for (var i = 0; i < this.user.profile.roles.length; i++) {
      if (this.user.profile.roles[i] == "Boatswain.Editor") {
        return true;
      }
    }
    return this.isAdmin();
  }

  loggedIn(): boolean {
    return this.user != null && !this.user.expired;
  }

  startLogin(dest: string = ''): Promise<void> {
    sessionStorage.setItem(actionKey, "true")
    sessionStorage.setItem(redirectKey, dest);
    return this.mgr.signinRedirect();
  }

  completeLogin(): Promise<void> {
    sessionStorage.setItem(actionKey, "false")
    return this.mgr.signinRedirectCallback().then(usr => {
      if (usr != null ) {
        this.user = usr;
        sessionStorage.setItem(userKey, this.user.toStorageString());
        this.events.emit(this.user);
        this.mgr.startSilentRenew();
      }
    })
    .finally(() => {
      this.router.navigate([sessionStorage.getItem(redirectKey)]);
      sessionStorage.setItem(redirectKey, '');
    });
  }

  startLogout(): Promise<void> {
    sessionStorage.setItem(actionKey, "true");
    sessionStorage.removeItem(userKey);
    return this.mgr.signoutRedirect();
  }

  completeLogout(): Promise<void> {
    sessionStorage.setItem(actionKey, "false")
    return this.mgr.signoutRedirectCallback().then(() => {
      this.mgr.revokeAccessToken();
      this.mgr.removeUser();
      this.mgr.stopSilentRenew();
      this.user = null;
    }).finally(() => {
      this.router.navigate(['']);
    });
  }
}
