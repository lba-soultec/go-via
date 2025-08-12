import { TestBed } from '@angular/core/testing';
import { AuthConfig, OAuthService } from 'angular-oauth2-oidc';
import { AuthenticationService } from './authentication.service';
import { StatehandlerService } from './statehandler.service';
import { StorageService } from './storage.service';

describe('AuthenticationService', () => {
  let service: AuthenticationService;
  let mockOAuthService: jasmine.SpyObj<OAuthService>;
  let mockStorageService: jasmine.SpyObj<StorageService>;
  let mockStatehandlerService: jasmine.SpyObj<StatehandlerService>;
  let mockAuthConfig: AuthConfig;

  beforeEach(() => {
    const oauthSpy = jasmine.createSpyObj('OAuthService', [
      'hasValidAccessToken',
      'hasValidIdToken',
      'configure',
      'setupAutomaticSilentRefresh',
      'loadDiscoveryDocumentAndTryLogin',
      'initCodeFlow',
      'logOut'
    ]);
    
    const storageSpy = jasmine.createSpyObj('StorageService', ['getItem', 'setItem', 'removeItem']);
    const statehandlerSpy = jasmine.createSpyObj('StatehandlerService', ['createState']);
    
    mockAuthConfig = {
      issuer: 'https://test.example.com',
      clientId: 'test-client',
      scope: 'openid profile',
      responseType: 'code',
      oidc: true
    };

    TestBed.configureTestingModule({
      providers: [
        AuthenticationService,
        { provide: OAuthService, useValue: oauthSpy },
        { provide: StorageService, useValue: storageSpy },
        { provide: StatehandlerService, useValue: statehandlerSpy },
        { provide: AuthConfig, useValue: mockAuthConfig }
      ]
    });

    service = TestBed.inject(AuthenticationService);
    mockOAuthService = TestBed.inject(OAuthService) as jasmine.SpyObj<OAuthService>;
    mockStorageService = TestBed.inject(StorageService) as jasmine.SpyObj<StorageService>;
    mockStatehandlerService = TestBed.inject(StatehandlerService) as jasmine.SpyObj<StatehandlerService>;
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should restore authentication state when tokens exist in storage', () => {
    // Arrange
    mockStorageService.getItem.and.callFake((key: string) => {
      if (key === 'access_token') return 'mock-access-token';
      if (key === 'id_token') return 'mock-id-token';
      return null;
    });

    // Act - create new instance to trigger initialization
    const newService = new AuthenticationService(
      mockOAuthService,
      mockAuthConfig,
      mockStatehandlerService,
      mockStorageService
    );

    // Assert
    expect(newService.authenticated).toBe(true);
  });

  it('should not be authenticated when no tokens exist in storage', () => {
    // Arrange
    mockStorageService.getItem.and.returnValue(null);

    // Act - create new instance to trigger initialization
    const newService = new AuthenticationService(
      mockOAuthService,
      mockAuthConfig,
      mockStatehandlerService,
      mockStorageService
    );

    // Assert
    expect(newService.authenticated).toBe(false);
  });

  it('should validate tokens when authenticate is called and tokens exist', async () => {
    // Arrange
    mockStorageService.getItem.and.callFake((key: string) => {
      if (key === 'access_token') return 'mock-access-token';
      if (key === 'id_token') return 'mock-id-token';
      return null;
    });
    mockOAuthService.hasValidAccessToken.and.returnValue(true);
    mockOAuthService.loadDiscoveryDocumentAndTryLogin.and.returnValue(Promise.resolve());

    // Act - create new instance and authenticate
    const newService = new AuthenticationService(
      mockOAuthService,
      mockAuthConfig,
      mockStatehandlerService,
      mockStorageService
    );
    const result = await newService.authenticate();

    // Assert
    expect(result).toBe(true);
    expect(mockOAuthService.configure).toHaveBeenCalledWith(mockAuthConfig);
    expect(mockOAuthService.hasValidAccessToken).toHaveBeenCalled();
  });

  it('should handle invalid tokens during authentication', async () => {
    // Arrange
    mockStorageService.getItem.and.callFake((key: string) => {
      if (key === 'access_token') return 'invalid-access-token';
      if (key === 'id_token') return 'invalid-id-token';
      return null;
    });
    mockOAuthService.hasValidAccessToken.and.returnValue(false);
    mockOAuthService.hasValidIdToken.and.returnValue(false);
    mockOAuthService.loadDiscoveryDocumentAndTryLogin.and.returnValue(Promise.resolve());
    mockStatehandlerService.createState.and.returnValue({ toPromise: () => Promise.resolve('mock-state') } as any);

    // Act - create new instance and authenticate
    const newService = new AuthenticationService(
      mockOAuthService,
      mockAuthConfig,
      mockStatehandlerService,
      mockStorageService
    );
    const result = await newService.authenticate();

    // Assert
    expect(newService.authenticated).toBe(false);
    expect(mockOAuthService.initCodeFlow).toHaveBeenCalled();
  });
});