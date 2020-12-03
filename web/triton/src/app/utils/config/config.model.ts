export interface IConfig {
    oidc: IOidcConfig
}

export interface IOidcConfig {
    authority: string
    clientId: string
    scope: string
}