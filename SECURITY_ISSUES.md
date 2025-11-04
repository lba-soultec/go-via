# Critical Security Issues Found in go-via

This document outlines 10 critical security issues discovered during a comprehensive security analysis of the go-via repository. Each issue should be addressed with high priority.

## Issue 1: Hardcoded Default Credentials (CRITICAL)

**Severity:** Critical
**Component:** User Authentication
**File:** `main.go` lines 155-158

### Description
The application uses hardcoded default credentials:
- Username: `admin`
- Password: `VMware1!`

### Security Impact
- Default credentials are widely known and documented in README
- Attackers can immediately gain administrative access
- No forced password change on first login

### Recommendation
- Force password change on first login
- Generate random default password
- Implement strong password policy
- Add multi-factor authentication

### Code Location
```go
hp := api.HashAndSalt([]byte("VMware1!"))
if res := db.DB.Where(models.User{UserForm: models.UserForm{Username: "admin"}}).Attrs(models.User{UserForm: models.UserForm{Password: hp}}).FirstOrCreate(&adm); res.Error != nil {
```

---

## Issue 2: Missing Vendor Dependencies (HIGH)

**Severity:** High
**Component:** Build System
**Files:** `go.mod`, `vendor/`

### Description
Critical dependencies are missing from the vendor directory:
- `github.com/pin/tftp`
- `github.com/rakyll/statik/fs`

### Security Impact
- Application fails to build
- Potential for dependency confusion attacks
- Supply chain vulnerabilities

### Recommendation
- Run `go mod vendor` to populate vendor directory
- Implement dependency scanning
- Pin dependency versions
- Regular security updates

### Error Output
```
cannot find module providing package github.com/pin/tftp: import lookup disabled by -mod=vendor
```

---

## Issue 3: Unauthenticated Critical Endpoints (CRITICAL)

**Severity:** Critical
**Component:** API Security
**File:** `main.go` line 179

### Description
The kickstart configuration endpoint (`ks.cfg`) is exposed without authentication.

### Security Impact
- Sensitive configuration data exposed
- Potential information disclosure
- Network configuration leakage

### Recommendation
- Implement authentication for all endpoints
- Use token-based authentication for automated systems
- Add IP-based access controls

### Code Location
```go
r.GET("ks.cfg", api.Ks(key))
```

---

## Issue 4: Weak Password Hashing Configuration (HIGH)

**Severity:** High
**Component:** Cryptography
**File:** `api/users.go` line 234

### Description
Password hashing uses `bcrypt.MinCost` which provides insufficient security.

### Security Impact
- Passwords vulnerable to brute force attacks
- Fast hash computation enables rainbow table attacks
- Non-compliance with security standards

### Recommendation
- Use bcrypt cost factor of at least 12
- Implement adaptive cost based on hardware
- Regular cost factor review

### Code Location
```go
hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
```

---

## Issue 5: Information Disclosure in Authentication (MEDIUM)

**Severity:** Medium
**Component:** Authentication
**File:** `api/login.go` lines 26-34

### Description
Login error messages reveal whether usernames exist in the system.

### Security Impact
- Username enumeration attacks
- Information gathering for targeted attacks
- Privacy violations

### Recommendation
- Use generic error messages
- Implement consistent response timing
- Add logging for failed attempts

### Code Location
```go
if res := db.DB.Where("username = ?", user.Username).First(&dbUser); res.Error != nil {
    logrus.WithFields(logrus.Fields{
        "username": user.Username,
        "status":   "supplied username does not exist",
    }).Info("auth")
    c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
```

---

## Issue 6: SQL Injection Vulnerability (HIGH)

**Severity:** High
**Component:** Database Security
**File:** `api/users.go` lines 88-90

### Description
The search functionality accepts arbitrary field names and values without proper validation.

### Security Impact
- Potential SQL injection attacks
- Unauthorized data access
- Database compromise

### Recommendation
- Implement input validation and sanitization
- Use allowlisted field names
- Add parameterized queries verification

### Code Location
```go
for k, v := range form {
    query = query.Where(k, v)
}
```

---

## Issue 7: Insecure Certificate Generation (MEDIUM)

**Severity:** Medium
**Component:** PKI/TLS
**File:** `crypto/main.go` lines 19-21, 100

### Description
Certificate generation uses hardcoded and predictable values:
- Fixed serial numbers (1653, 1658)
- Hardcoded subject information
- Predictable SubjectKeyId

### Security Impact
- Certificate collisions possible
- Weak certificate validation
- Non-compliance with PKI standards

### Recommendation
- Generate random serial numbers
- Use crypto/rand for SubjectKeyId
- Implement proper certificate validation

### Code Location
```go
SerialNumber: big.NewInt(1653),
SubjectKeyId: []byte{1, 2, 3, 4, 6},
```

---

## Issue 8: DoS via Panic Statements (HIGH)

**Severity:** High
**Component:** Error Handling
**File:** `secrets/main.go` multiple locations

### Description
Encryption and decryption functions use `panic()` for error handling, causing service crashes.

### Security Impact
- Denial of Service attacks
- Service instability
- Potential for exploitation

### Recommendation
- Replace panic() with proper error handling
- Implement graceful error recovery
- Add comprehensive logging

### Code Location
```go
if err != nil {
    panic(err.Error())
}
```

---

## Issue 9: Missing TLS Security Configuration (MEDIUM)

**Severity:** Medium
**Component:** Network Security
**File:** `main.go` line 332

### Description
No TLS configuration validation or security hardening is implemented.

### Security Impact
- Weak cipher suites may be used
- Vulnerable to downgrade attacks
- Non-compliance with security standards

### Recommendation
- Implement secure TLS configuration
- Disable weak cipher suites
- Enforce minimum TLS version (1.2+)
- Add certificate validation

### Code Location
```go
err = r.RunTLS(listen, "./cert/server.crt", "./cert/server.key")
```

---

## Issue 10: No Rate Limiting Protection (HIGH)

**Severity:** High
**Component:** API Security
**File:** `api/login.go`

### Description
The login endpoint lacks rate limiting protection against brute force attacks.

### Security Impact
- Brute force password attacks
- Account enumeration
- Resource exhaustion

### Recommendation
- Implement rate limiting middleware
- Add account lockout mechanisms
- Implement CAPTCHA for repeated failures
- Add IP-based blocking

### Current Implementation
No rate limiting or brute force protection exists.

---

## Summary

These 10 critical issues represent significant security vulnerabilities that should be addressed immediately. Priority should be given to:

1. **Critical Issues (1, 3)** - Immediate fix required
2. **High Severity Issues (2, 4, 6, 8, 10)** - Fix within 1 week
3. **Medium Severity Issues (5, 7, 9)** - Fix within 1 month

Each issue includes specific recommendations and code locations to facilitate quick remediation.