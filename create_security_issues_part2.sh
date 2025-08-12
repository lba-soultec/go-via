#!/bin/bash

# Part 2 of GitHub Issues Creation Script for remaining security issues

REPO="lba-soultec/go-via"

# Issue 7: Insecure Certificate Generation
gh issue create \
  --repo "$REPO" \
  --title "🟡 MEDIUM: Insecure Certificate Generation with Predictable Values" \
  --label "security,medium,cryptography,certificates" \
  --body "## Security Vulnerability: Insecure Certificate Generation

**Severity:** Medium  
**CWE:** CWE-330 (Use of Insufficiently Random Values)  
**Component:** PKI/TLS  
**File:** \`crypto/main.go\` lines 19-21, 100

### Description
Certificate generation uses hardcoded and predictable values that compromise certificate security.

### Security Impact
- ⚠️ Certificate collisions possible
- ⚠️ Weak certificate validation
- ⚠️ Non-compliance with PKI standards
- ⚠️ Potential for certificate spoofing

### Affected Code
\`\`\`go
ca := &x509.Certificate{
    SerialNumber: big.NewInt(1653),  // Hardcoded
    // ...
}

cert := &x509.Certificate{
    SerialNumber: big.NewInt(1658),  // Hardcoded
    SubjectKeyId: []byte{1, 2, 3, 4, 6},  // Predictable
    // ...
}
\`\`\`

### Issues
1. Fixed serial numbers (1653, 1658)
2. Hardcoded subject information
3. Predictable SubjectKeyId values

### Recommendations
1. 🔧 Generate random serial numbers using crypto/rand
2. 🔧 Use crypto/rand for SubjectKeyId generation
3. 🔧 Implement proper certificate validation
4. 🔧 Add certificate authority best practices

### Code Fix
\`\`\`go
serialNumber, _ := rand.Int(rand.Reader, big.NewInt(1000000))
subjectKeyId := make([]byte, 20)
rand.Read(subjectKeyId)
\`\`\`"

# Issue 8: DoS via Panic Statements
gh issue create \
  --repo "$REPO" \
  --title "🟠 HIGH: Denial of Service via Panic Statements in Crypto Functions" \
  --label "security,high,availability,error-handling" \
  --body "## Security Vulnerability: DoS via Panic Statements

**Severity:** High  
**CWE:** CWE-248 (Uncaught Exception)  
**Component:** Error Handling  
**File:** \`secrets/main.go\` multiple locations

### Description
Encryption and decryption functions use \`panic()\` for error handling, causing immediate service crashes when errors occur.

### Security Impact
- ⚠️ Denial of Service attacks possible
- ⚠️ Service instability and crashes
- ⚠️ Potential for exploitation by malformed inputs
- ⚠️ Poor user experience and availability

### Affected Code
\`\`\`go
func Encrypt(stringToEncrypt string, keyString string) (encryptedString string) {
    if err != nil {
        panic(err.Error())  // Causes service crash
    }
}

func Decrypt(encryptedString string, keyString string) (decryptedString string) {
    if err != nil {
        panic(err.Error())  // Causes service crash
    }
}
\`\`\`

### Attack Scenario
1. Attacker sends malformed encrypted data
2. Decrypt function encounters error
3. panic() is called, crashing the service
4. Service becomes unavailable

### Recommendations
1. 🔧 Replace panic() with proper error handling
2. 🔧 Implement graceful error recovery
3. 🔧 Add comprehensive logging without crashing
4. 🔧 Validate inputs before processing
5. 🔧 Add error boundaries and recovery mechanisms

### Code Fix
\`\`\`go
func Decrypt(encryptedString string, keyString string) (string, error) {
    // ... crypto operations ...
    if err != nil {
        logrus.WithError(err).Error(\"Decryption failed\")
        return \"\", fmt.Errorf(\"decryption failed: %w\", err)
    }
    return string(plaintext), nil
}
\`\`\`"

# Issue 9: Missing TLS Security Configuration
gh issue create \
  --repo "$REPO" \
  --title "🟡 MEDIUM: Missing TLS Security Configuration" \
  --label "security,medium,tls,network" \
  --body "## Security Vulnerability: Missing TLS Security Configuration

**Severity:** Medium  
**CWE:** CWE-326 (Inadequate Encryption Strength)  
**Component:** Network Security  
**File:** \`main.go\` line 332

### Description
The TLS server configuration lacks security hardening, potentially allowing weak cipher suites and protocols.

### Security Impact
- ⚠️ Weak cipher suites may be used
- ⚠️ Vulnerable to downgrade attacks
- ⚠️ Non-compliance with security standards (TLS 1.3)
- ⚠️ Potential for man-in-the-middle attacks

### Affected Code
\`\`\`go
err = r.RunTLS(listen, \"./cert/server.crt\", \"./cert/server.key\")
\`\`\`

### Missing Security Features
1. No minimum TLS version enforcement
2. No cipher suite restrictions
3. No certificate validation requirements
4. No HSTS headers
5. No perfect forward secrecy enforcement

### Recommendations
1. 🔧 Implement secure TLS configuration
2. 🔧 Disable weak cipher suites
3. 🔧 Enforce minimum TLS version (1.2+, prefer 1.3)
4. 🔧 Add certificate validation
5. 🔧 Implement HSTS headers
6. 🔧 Add perfect forward secrecy

### Secure Configuration Example
\`\`\`go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    CurvePreferences: []tls.CurveID{
        tls.CurveP256,
        tls.X25519,
    },
    PreferServerCipherSuites: true,
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
    },
}
\`\`\`"

# Issue 10: No Rate Limiting Protection
gh issue create \
  --repo "$REPO" \
  --title "🟠 HIGH: Missing Rate Limiting on Authentication Endpoints" \
  --label "security,high,rate-limiting,authentication" \
  --body "## Security Vulnerability: No Rate Limiting Protection

**Severity:** High  
**CWE:** CWE-307 (Improper Restriction of Excessive Authentication Attempts)  
**Component:** API Security  
**File:** \`api/login.go\`

### Description
The login endpoint lacks rate limiting protection, making it vulnerable to brute force attacks.

### Security Impact
- ⚠️ Brute force password attacks possible
- ⚠️ Account enumeration through repeated attempts
- ⚠️ Resource exhaustion attacks
- ⚠️ Potential for account compromise
- ⚠️ No protection against automated attacks

### Current Implementation
No rate limiting, account lockout, or brute force protection exists on the login endpoint.

### Attack Scenario
1. Attacker identifies login endpoint
2. Launches automated brute force attack
3. Attempts thousands of password combinations
4. Successfully compromises weak passwords
5. Gains unauthorized access

### Recommendations
1. 🔧 Implement rate limiting middleware (e.g., 5 attempts per minute)
2. 🔧 Add account lockout mechanisms after failed attempts
3. 🔧 Implement CAPTCHA for repeated failures
4. 🔧 Add IP-based blocking for suspicious activity
5. 🔧 Implement progressive delays for failed attempts
6. 🔧 Add comprehensive audit logging

### Implementation Example
\`\`\`go
// Rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(time.Minute), 5)
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                \"error\": \"Too many requests\",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
\`\`\`

### Additional Security Measures
- Account lockout after 5 failed attempts
- Progressive delay: 1s, 2s, 4s, 8s, 16s
- CAPTCHA after 3 failed attempts
- Email notification on repeated failures"

echo "All security issues prepared for creation. Execute commands manually with GitHub CLI."