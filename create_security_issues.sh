#!/bin/bash

# GitHub Issues Creation Script
# This script contains the commands to create GitHub issues for the critical security issues found
# Since direct issue creation is not available, these commands can be executed manually

REPO="lba-soultec/go-via"

echo "Creating GitHub issues for critical security vulnerabilities..."

# Issue 1: Hardcoded Default Credentials
gh issue create \
  --repo "$REPO" \
  --title "🔴 CRITICAL: Hardcoded Default Credentials Security Vulnerability" \
  --label "security,critical,authentication" \
  --body "## Security Vulnerability: Hardcoded Default Credentials

**Severity:** Critical  
**CWE:** CWE-798 (Use of Hard-coded Credentials)  
**Component:** User Authentication  
**File:** \`main.go\` lines 155-158

### Description
The application uses hardcoded default credentials that are publicly documented:
- Username: \`admin\`
- Password: \`VMware1!\`

### Security Impact
- ⚠️ Default credentials are widely known and documented in README
- ⚠️ Attackers can immediately gain administrative access
- ⚠️ No forced password change on first login
- ⚠️ Potential for immediate system compromise

### Affected Code
\`\`\`go
hp := api.HashAndSalt([]byte(\"VMware1!\"))
if res := db.DB.Where(models.User{UserForm: models.UserForm{Username: \"admin\"}}).Attrs(models.User{UserForm: models.UserForm{Password: hp}}).FirstOrCreate(&adm); res.Error != nil {
\`\`\`

### Recommendations
1. 🔧 Force password change on first login
2. 🔧 Generate random default password
3. 🔧 Implement strong password policy
4. 🔧 Add multi-factor authentication
5. 🔧 Remove credentials from documentation

### Priority
This issue requires **immediate attention** as it poses a critical security risk."

# Issue 2: Missing Vendor Dependencies
gh issue create \
  --repo "$REPO" \
  --title "🟠 HIGH: Missing Vendor Dependencies Causing Build Failures" \
  --label "build,dependencies,high" \
  --body "## Build Issue: Missing Vendor Dependencies

**Severity:** High  
**Component:** Build System  
**Files:** \`go.mod\`, \`vendor/\`

### Description
Critical dependencies are missing from the vendor directory causing build failures:
- \`github.com/pin/tftp\`
- \`github.com/rakyll/statik/fs\`

### Security Impact
- 🔧 Application fails to build
- ⚠️ Potential for dependency confusion attacks
- ⚠️ Supply chain vulnerabilities
- ⚠️ Inability to deploy secure updates

### Error Output
\`\`\`
cannot find module providing package github.com/pin/tftp: import lookup disabled by -mod=vendor
cannot find module providing package github.com/rakyll/statik/fs: import lookup disabled by -mod=vendor
\`\`\`

### Recommendations
1. Run \`go mod vendor\` to populate vendor directory
2. Implement dependency scanning in CI/CD
3. Pin dependency versions
4. Regular security updates schedule
5. Add dependency vulnerability scanning

### Steps to Reproduce
1. Clone repository
2. Run \`go build\`
3. Observe build failure"

# Issue 3: Unauthenticated Critical Endpoints
gh issue create \
  --repo "$REPO" \
  --title "🔴 CRITICAL: Unauthenticated Access to Sensitive Configuration Endpoint" \
  --label "security,critical,authentication,api" \
  --body "## Security Vulnerability: Unauthenticated Critical Endpoints

**Severity:** Critical  
**CWE:** CWE-306 (Missing Authentication for Critical Function)  
**Component:** API Security  
**File:** \`main.go\` line 179

### Description
The kickstart configuration endpoint (\`ks.cfg\`) is exposed without any authentication mechanism.

### Security Impact
- ⚠️ Sensitive configuration data exposed to unauthorized users
- ⚠️ Potential information disclosure of network configurations
- ⚠️ Infrastructure details accessible to attackers
- ⚠️ Violation of principle of least privilege

### Affected Code
\`\`\`go
r.GET(\"ks.cfg\", api.Ks(key))
\`\`\`

### Attack Scenario
1. Attacker discovers the service
2. Accesses \`/ks.cfg\` endpoint without authentication
3. Obtains sensitive configuration information
4. Uses information for further attacks

### Recommendations
1. 🔧 Implement authentication for all endpoints
2. 🔧 Use token-based authentication for automated systems
3. 🔧 Add IP-based access controls
4. 🔧 Implement proper authorization checks
5. 🔧 Add audit logging for configuration access

### Priority
This issue requires **immediate attention** as it exposes sensitive data."

# Issue 4: Weak Password Hashing
gh issue create \
  --repo "$REPO" \
  --title "🟠 HIGH: Weak Password Hashing Configuration" \
  --label "security,high,cryptography,passwords" \
  --body "## Security Vulnerability: Weak Password Hashing Configuration

**Severity:** High  
**CWE:** CWE-327 (Use of a Broken or Risky Cryptographic Algorithm)  
**Component:** Cryptography  
**File:** \`api/users.go\` line 234

### Description
Password hashing uses \`bcrypt.MinCost\` (cost factor 4) which provides insufficient security against modern attacks.

### Security Impact
- ⚠️ Passwords vulnerable to brute force attacks
- ⚠️ Fast hash computation enables rainbow table attacks
- ⚠️ Non-compliance with OWASP security standards
- ⚠️ Inadequate protection for user credentials

### Affected Code
\`\`\`go
hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
\`\`\`

### Technical Details
- Current cost factor: 4 (MinCost)
- Recommended minimum: 12
- Time to hash with cost 4: ~1ms
- Time to hash with cost 12: ~250ms

### Recommendations
1. 🔧 Use bcrypt cost factor of at least 12
2. 🔧 Implement adaptive cost based on hardware capabilities
3. 🔧 Regular cost factor review and updates
4. 🔧 Consider migration to Argon2id for new implementations

### Code Fix
\`\`\`go
const BCRYPT_COST = 12
hash, err := bcrypt.GenerateFromPassword(pwd, BCRYPT_COST)
\`\`\`"

# Issue 5: Information Disclosure in Authentication
gh issue create \
  --repo "$REPO" \
  --title "🟡 MEDIUM: Username Enumeration via Authentication Errors" \
  --label "security,medium,authentication,information-disclosure" \
  --body "## Security Vulnerability: Information Disclosure in Authentication

**Severity:** Medium  
**CWE:** CWE-204 (Observable Response Discrepancy)  
**Component:** Authentication  
**File:** \`api/login.go\` lines 26-34

### Description
Login error messages reveal whether usernames exist in the system, enabling username enumeration attacks.

### Security Impact
- ⚠️ Username enumeration attacks possible
- ⚠️ Information gathering for targeted attacks
- ⚠️ Privacy violations
- ⚠️ Facilitates social engineering attacks

### Affected Code
\`\`\`go
if res := db.DB.Where(\"username = ?\", user.Username).First(&dbUser); res.Error != nil {
    logrus.WithFields(logrus.Fields{
        \"username\": user.Username,
        \"status\":   \"supplied username does not exist\",
    }).Info(\"auth\")
    c.JSON(http.StatusUnauthorized, gin.H{\"error\": \"invalid username or password\"})
\`\`\`

### Attack Scenario
1. Attacker tries login with known username
2. System logs reveal username existence
3. Attacker builds list of valid usernames
4. Focused brute force attacks on valid accounts

### Recommendations
1. 🔧 Use generic error messages for all authentication failures
2. 🔧 Implement consistent response timing
3. 🔧 Remove username information from logs
4. 🔧 Add comprehensive audit logging without sensitive data"

# Issue 6: SQL Injection Vulnerability
gh issue create \
  --repo "$REPO" \
  --title "🟠 HIGH: Potential SQL Injection in Search Functionality" \
  --label "security,high,sql-injection,database" \
  --body "## Security Vulnerability: SQL Injection in Search Functionality

**Severity:** High  
**CWE:** CWE-89 (SQL Injection)  
**Component:** Database Security  
**File:** \`api/users.go\` lines 88-90

### Description
The search functionality accepts arbitrary field names and values without proper validation, potentially leading to SQL injection attacks.

### Security Impact
- ⚠️ Potential SQL injection attacks
- ⚠️ Unauthorized data access
- ⚠️ Database compromise
- ⚠️ Data integrity violations
- ⚠️ Potential for privilege escalation

### Affected Code
\`\`\`go
for k, v := range form {
    query = query.Where(k, v)
}
\`\`\`

### Attack Scenario
1. Attacker sends malicious field names/values
2. ORM constructs unsafe queries
3. Malicious SQL executed
4. Database compromise achieved

### Recommendations
1. 🔧 Implement input validation and sanitization
2. 🔧 Use allowlisted field names only
3. 🔧 Add parameterized queries verification
4. 🔧 Implement proper access controls
5. 🔧 Add SQL injection detection and prevention

### Code Fix
\`\`\`go
allowedFields := map[string]bool{
    \"username\": true,
    \"email\": true,
}
for k, v := range form {
    if allowedFields[k] {
        query = query.Where(k, v)
    }
}
\`\`\`"

echo "Issue creation commands prepared. Execute each 'gh issue create' command manually or run this script with gh CLI installed."