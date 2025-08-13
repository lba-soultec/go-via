import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { APP_INITIALIZER, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { AuthConfig, OAuthModule, OAuthStorage } from 'angular-oauth2-oidc';
import { ClarityModule } from '@clr/angular';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './components/home/home.component';
import { SignedOutComponent } from './components/signed-out/signed-out.component';
import { UserComponent } from './components/user/user.component';
import { StatehandlerProcessorService, StatehandlerProcessorServiceImpl } from './services/statehandler-processor.service';
import { StatehandlerService, StatehandlerServiceImpl } from './services/statehandler.service';
import { StorageService } from './services/storage.service';
import { ConfigService } from './services/config.service';
import { AuthenticationService } from './services/authentication.service';

const authConfig: AuthConfig = {
    
};


let allowedUrls: string[] = [];

const loadConfig = (configService: ConfigService) => {
  return () =>
    configService.loadConfig().then((config) => {
    Object.assign(authConfig, config.authConfig); // Merge with the loaded config
    Object.assign(allowedUrls, config.allowedUrls || []);
    });
    
};

const initializeAuth = (authService: AuthenticationService, configService: ConfigService) => {
  return () => {
    // Ensure config is loaded first, then initialize authentication
    return configService.getConfig() ? authService.initializeOnStartup() : Promise.resolve(false);
  };
};

const stateHandlerFn = (stateHandler: StatehandlerService) => {
  return () => {
    return stateHandler.initStateHandler();
  };
};

@NgModule({ declarations: [AppComponent, SignedOutComponent, UserComponent, HomeComponent],
    bootstrap: [AppComponent], 
    imports: [
        ClarityModule,
        
        BrowserModule,
        AppRoutingModule,
        OAuthModule.forRoot({
            resourceServer: {
                allowedUrls:  allowedUrls,
                sendAccessToken: true,
            },
        })], providers: [
        {
            provide: APP_INITIALIZER,
            useFactory: loadConfig,
            multi: true,
            deps: [ConfigService],
        },
        {
            provide: APP_INITIALIZER,
            useFactory: initializeAuth,
            multi: true,
            deps: [AuthenticationService, ConfigService],
        },
        {
            provide: APP_INITIALIZER,
            useFactory: stateHandlerFn,
            multi: true,
            deps: [StatehandlerService],
        },
        {
            provide: AuthConfig,
            useFactory: () => authConfig,
        },
        {
            provide: StatehandlerProcessorService,
            useClass: StatehandlerProcessorServiceImpl,
        },
        {
            provide: StatehandlerService,
            useClass: StatehandlerServiceImpl,
        },
        {
            provide: OAuthStorage,
            useClass: StorageService,
        },
        provideHttpClient(withInterceptorsFromDi()),
    ] })
export class AppModule {}
