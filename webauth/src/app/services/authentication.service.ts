import { Injectable } from '@angular/core';
import { AuthConfig, OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, from, Observable } from 'rxjs';

import { StatehandlerService } from './statehandler.service';

@Injectable({
  providedIn: 'root',
})
export class AuthenticationService {
  private _authenticated: boolean = false;
  private readonly _authenticationChanged: BehaviorSubject<boolean> = new BehaviorSubject(this.authenticated);

  constructor(
    private oauthService: OAuthService,
    private authConfig: AuthConfig,
    private statehandler: StatehandlerService,
  ) {}

  public get authenticated(): boolean {
    return this._authenticated;
  }

  public get authenticationChanged(): Observable<boolean> {
    return this._authenticationChanged;
  }

  public getOIDCUser(): Observable<any> {
    return from(this.oauthService.loadUserProfile());
  }

  public async initializeOnStartup(): Promise<boolean> {
    // Configure OAuth service with loaded config
    this.oauthService.configure(this.authConfig);
    this.oauthService.setupAutomaticSilentRefresh();
    this.oauthService.strictDiscoveryDocumentValidation = false;
    
    // Try to restore tokens from storage without triggering login flow
    await this.oauthService.loadDiscoveryDocumentAndTryLogin();
    
    // Update authentication state based on token validity
    this._authenticated = this.oauthService.hasValidAccessToken();
    this._authenticationChanged.next(this.authenticated);
    
    return this.authenticated;
  }

  public async authenticate(setState: boolean = true): Promise<boolean> {
    this.oauthService.configure(this.authConfig);
    this.oauthService.setupAutomaticSilentRefresh();

    this.oauthService.strictDiscoveryDocumentValidation = false;
    await this.oauthService.loadDiscoveryDocumentAndTryLogin();

    this._authenticated = this.oauthService.hasValidAccessToken();

    if (!this.oauthService.hasValidIdToken() || !this.authenticated) {
      const newState = setState ? await this.statehandler.createState().toPromise() : undefined;
      this.oauthService.initCodeFlow(newState);
    }
    this._authenticationChanged.next(this.authenticated);

    return this.authenticated;
  }

  public signout(): void {
    this.oauthService.logOut();
    this._authenticated = false;
    this._authenticationChanged.next(false);
  }
}
