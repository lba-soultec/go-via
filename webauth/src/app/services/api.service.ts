import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment'; // Import environment
import { AuthenticationService } from './authentication.service'; // Import AuthenticationService
import { ConfigService } from './config.service';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  billingApiEndpoint: string;

  constructor(
    private httpClient: HttpClient,
    private authService: AuthenticationService, // Inject AuthenticationService
    private configService: ConfigService // Inject ConfigService
  ) {
    this.billingApiEndpoint = this.configService.getEnv('billingApiEndpoint'); // Read billingApiEndpoint from ConfigService
  }

  public getCustomers(): Observable<any> {
    const authToken = this.authService['oauthService'].getAccessToken(); // Use OAuthService to get the token
    return this.httpClient.get(this.billingApiEndpoint + '/api/getCustomers', {
      headers: { 'Authorization': `Bearer ${authToken}` }
    });
  }

  public getRoutes(): Observable<any> {
    const authToken = this.authService['oauthService'].getAccessToken(); // Use OAuthService to get the token
    return this.httpClient.get(this.billingApiEndpoint + '/api/getRoutes', {
      headers: { 'Authorization': `Bearer ${authToken}` }
    });
  }

  public generateCharts(payload: { year: string; month: string }): Observable<any> {
    const authToken = this.authService['oauthService'].getAccessToken(); // Use OAuthService to get the token
    return this.httpClient.post(this.billingApiEndpoint + '/api/generateCharts', payload, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken}`
      }
    });
  }

  public calculateOvercommitment(payload: { year: string; month: string; customerID: string; limit: string }): Observable<any> {
    const authToken = this.authService['oauthService'].getAccessToken(); // Use OAuthService to get the token
    return this.httpClient.post(this.billingApiEndpoint + '/api/calculateOvercommitment', payload, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken}`
      }
    });
  }
}