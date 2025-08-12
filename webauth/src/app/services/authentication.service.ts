import { Injectable } from '@angular/core';
import { AuthConfig, OAuthService } from 'angular-oauth2-oidc';
import { BehaviorSubject, from, Observable } from 'rxjs';

import { StatehandlerService } from './statehandler.service';
import { StorageService } from './storage.service';

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
    private storageService: StorageService,
  ) {
    // Restore authentication state on service initialization
    this.initializeAuthenticationState();
  }

  private initializeAuthenticationState(): void {
    // Check if we have stored tokens using the storage service
    const accessToken = this.storageService.getItem('access_token');
    const idToken = this.storageService.getItem('id_token');
    
    if (accessToken && idToken) {
      // We have tokens, but we need to validate them properly
      // Since config might not be loaded yet, we'll assume authenticated for now
      // and validate in the authenticate method
      this._authenticated = true;
    } else {
      this._authenticated = false;
    }
    
    // Emit the initial authentication state
    this._authenticationChanged.next(this._authenticated);
  }

  public get authenticated(): boolean {
    return this._authenticated;
  }

  public get authenticationChanged(): Observable<boolean> {
    return this._authenticationChanged;
  }

  public getOIDCUser(): Observable<any> {
    return from(this.oauthService.loadUserProfile());
  }

  public async authenticate(setState: boolean = true): Promise<boolean> {
    // Configure OAuth service if not already configured
    if (this.authConfig && this.authConfig.issuer) {
      this.oauthService.configure(this.authConfig);
    }
    
    // If we think we're authenticated, validate the tokens first
    if (this._authenticated) {
      // Check if tokens are still valid
      if (this.oauthService.hasValidAccessToken()) {
        return true;
      } else {
        // Tokens are invalid, need to re-authenticate
        this._authenticated = false;
        this._authenticationChanged.next(false);
      }
    }

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
