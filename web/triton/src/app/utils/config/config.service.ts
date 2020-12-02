import { HttpClient } from '@angular/common/http';
import { Injectable, isDevMode } from '@angular/core';
import { environment } from 'src/environments/environment';
import { IConfig, IOidcConfig } from './config.model';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {
  private static config: IConfig;

  constructor(private http: HttpClient) {}

  getOIDC(): IOidcConfig {
    return ConfigService.config.oidc;
  }

  load() {
    console.log(environment.name);
    console.log(isDevMode());
    const file = `assets/config/config.${environment.name}.json`
    return new Promise<void>((resolve, reject) => {
      this.http.get(file).toPromise().then((response: IConfig) => {
        ConfigService.config = response;
        resolve();
      }).catch((_: any) => {
        reject(`could not load ${file}`);
      })
    });
  }
}
