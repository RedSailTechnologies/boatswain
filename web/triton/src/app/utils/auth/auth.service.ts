import { EventEmitter, Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Profile, User, UserManager, UserManagerSettings } from 'oidc-client';

// TODO AdamP - refactor
const settings: UserManagerSettings = {
  authority: 'http://localhost:4011/',
  client_id: 'implicit-mock-client',
  redirect_uri: 'http://localhost:4200/login',
  post_logout_redirect_uri: 'http://localhost:4200/logout',
  response_type: 'id_token token',
  scope: 'openid profile',
  filterProtocolClaims: true,
  loadUserInfo: true,
};

const key: string = "loginRedirectDestination";
const action: string = "loginActionInProgress";

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private mgr: UserManager = new UserManager(settings);
  private user: User;

  public static IN_PROGRESS: boolean = (sessionStorage.getItem(action) ?? "false") != "false";
  public events: EventEmitter<User> = new EventEmitter<User>();

  constructor(private router: Router) {}

  authHeader(): string {
    if (this.user == null) {
      return null;
    }
    return `${this.user.token_type} ${this.user.access_token}`;
  }

  getProfile(): Profile {
    if (this.user == null) {
      return null;
    }
    return this.user.profile;
  }

  getScopes(): string[] {
    if (this.user == null) {
      return null;
    }
    return this.user.scopes;
  }

  loggedIn(): boolean {
    return this.user != null && !this.user.expired;
  }

  startLogin(dest: string = ''): Promise<void> {
    sessionStorage.setItem(action, "true")
    sessionStorage.setItem(key, dest);
    return this.mgr.signinRedirect();
  }

  completeLogin(): Promise<void> {
    sessionStorage.setItem(action, "false")
    return this.mgr.signinRedirectCallback().then(usr => {
      if (usr != null ) {
        this.user = usr;
        this.events.emit(this.user);
        this.mgr.startSilentRenew();
      }
    }).finally(() => {
      sessionStorage.setItem(key, '');
      this.router.navigate([sessionStorage.getItem(key)]);
    });
  }

  startLogout(): Promise<void> {
    sessionStorage.setItem(action, "true")
    return this.mgr.signoutRedirect();
  }

  completeLogout(): Promise<void> {
    sessionStorage.setItem(action, "false")
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
