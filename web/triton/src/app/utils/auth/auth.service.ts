import { EventEmitter, Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Profile, User, UserManager, UserManagerSettings } from 'oidc-client';
import * as fetch from 'isomorphic-fetch';

// TODO AdamP - refactor
const settings: UserManagerSettings = {
  // authority: 'https://login.microsoftonline.com/8c91e3f4-7f37-4334-9f70-fae3f5235c18/v2.0/',
  // client_id: '071e7d94-aa7e-42aa-8ff3-3ca84b9c9e06',
  // scope: 'openid profile api://071e7d94-aa7e-42aa-8ff3-3ca84b9c9e06/boatswain',
  authority: 'http://localhost:4011',
  client_id: 'implicit-mock-client',
  scope: 'openid profile boatswain',

  redirect_uri: 'http://localhost:4200/login',
  post_logout_redirect_uri: 'http://localhost:4200/logout',
  response_type: 'id_token token',
  filterProtocolClaims: true,
  loadUserInfo: false,
};

const actionKey: string = "loginActionInProgress";
const redirectKey: string = "loginRedirectDestination";
const userKey: string = "currentUser";

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private mgr: UserManager = new UserManager(settings);
  private user: User;

  public static IN_PROGRESS: boolean = (sessionStorage.getItem(actionKey) ?? "false") != "false";
  public events: EventEmitter<User> = new EventEmitter<User>();

  constructor(private router: Router) {
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
      // TODO - remove me
      console.log(usr.scope)
      console.log(usr.profile.roles)
      console.log(usr)

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
