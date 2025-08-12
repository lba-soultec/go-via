# Authentication Persistence Fix

## Problem Statement
Users were being logged out when refreshing the website. The requirement was to maintain authentication state across page refreshes.

## Root Cause Analysis
The issue was in the `AuthenticationService` where the `_authenticated` field was being reset to `false` on every service instantiation (page refresh), even though valid tokens might still exist in localStorage through the `StorageService`.

The authentication flow was:
1. User logs in successfully → tokens stored in localStorage via `StorageService`
2. User refreshes page → `AuthenticationService` constructor runs
3. `_authenticated` field initialized to `false` → User appears logged out
4. Valid tokens still exist in storage but are not checked during initialization

## Solution Implementation

### Changes Made

#### 1. Added `initializeAuthenticationState()` method
```typescript
private initializeAuthenticationState(): void {
  // Check if we have stored tokens using the storage service
  const accessToken = this.storageService.getItem('access_token');
  const idToken = this.storageService.getItem('id_token');
  
  if (accessToken && idToken) {
    // We have tokens, assume authenticated for now
    // Full validation happens in authenticate() method
    this._authenticated = true;
  } else {
    this._authenticated = false;
  }
  
  // Emit the initial authentication state
  this._authenticationChanged.next(this._authenticated);
}
```

#### 2. Modified Constructor
```typescript
constructor(
  private oauthService: OAuthService,
  private authConfig: AuthConfig,
  private statehandler: StatehandlerService,
  private storageService: StorageService,
) {
  // Restore authentication state on service initialization
  this.initializeAuthenticationState();
}
```

#### 3. Enhanced `authenticate()` method
```typescript
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
  
  // ... rest of authentication logic
}
```

## How the Fix Works

### Before the Fix
1. Page loads → `AuthenticationService` instantiated
2. `_authenticated = false` (hardcoded)
3. User appears logged out despite valid tokens in storage
4. Guard redirects to login even with valid tokens

### After the Fix
1. Page loads → `AuthenticationService` instantiated
2. `initializeAuthenticationState()` called
3. Check for tokens in storage using `StorageService`
4. If tokens exist → `_authenticated = true`
5. When routes are accessed, `authenticate()` validates tokens
6. If tokens valid → user stays logged in
7. If tokens invalid → proper re-authentication flow

## Key Benefits

1. **Immediate State Restoration**: Authentication state is restored immediately on page load
2. **Lazy Validation**: Full token validation only happens when needed (accessing protected routes)
3. **Graceful Degradation**: If tokens are invalid, the normal authentication flow proceeds
4. **No Breaking Changes**: Existing authentication flow remains unchanged for new logins

## User Experience Impact

### Before
- User logs in successfully
- User refreshes page
- **User appears logged out** (redirected to login)
- User has to log in again

### After
- User logs in successfully  
- User refreshes page
- **User remains logged in** (stays on protected page or can access protected routes)
- Seamless experience across page refreshes

## Technical Considerations

1. **Config Loading**: The fix handles cases where OAuth configuration might not be loaded yet during initialization
2. **Token Storage**: Uses the existing `StorageService` with proper key prefixing (`zitadel:` prefix)
3. **Observable Pattern**: Maintains the existing reactive authentication state pattern
4. **Error Handling**: Gracefully handles invalid tokens by falling back to normal authentication flow

## Testing

Added comprehensive test suite covering:
- Token restoration from storage
- Handling missing tokens
- Token validation during authentication
- Invalid token scenarios

The fix ensures that authentication persistence works reliably while maintaining all existing functionality.