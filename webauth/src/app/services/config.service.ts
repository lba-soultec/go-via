import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class ConfigService {
    private config: any;

    constructor(private http: HttpClient) {}
  
  getEnv(key: string): string {
    return (window as any).__env[key] || '';
  }
  


  loadConfig(): Promise<any> {
    return this.http
      .get('/assets/config/config.json')
      .toPromise()
      .then((config) => {
        this.config = config;
        return config;
      });
  }

  getConfig(): any {
    return this.config;
  }
}