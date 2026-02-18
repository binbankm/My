# Username Whitespace Trimming Fix

## Problem Statement

Users reported experiencing "Invalid credentials" errors when attempting to log in, even with correct credentials. The issue was described as "出现改写登录，为什么一直invalid credentials" (Login rewrite appeared, why always invalid credentials).

## Root Cause

The login function in `internal/api/auth.go` did not sanitize username input before authentication. When users accidentally included leading or trailing whitespace in the username field (commonly occurring when copy-pasting credentials from password managers or documentation), the authentication would fail because:

1. Database query used exact match: `WHERE username = ?`
2. Username stored in database: `"admin"`
3. Username with spaces: `" admin "` or `" admin"` or `"admin "`
4. Query result: no match found → "Invalid credentials" error

## Solution

Added input sanitization by trimming whitespace from the username field:

```go
// Trim whitespace from username to prevent login failures due to accidental spaces
// Note: We don't trim password as users may intentionally use spaces in their password
req.Username = strings.TrimSpace(req.Username)
```

### Why Only Username?

- **Username**: Standard practice to trim usernames as they typically don't contain intentional spaces
- **Password**: Preserved as-is because users may intentionally include spaces in their passwords for security

## Changes Made

**File**: `internal/api/auth.go`

1. Added `strings` import
2. Added `strings.TrimSpace()` call for username sanitization
3. Added comment explaining why password is not trimmed

Total: 5 lines added (minimal change)

## Testing

All scenarios tested and verified:

| Test Case | Username | Password | Expected | Result |
|-----------|----------|----------|----------|--------|
| Normal login | `admin` | `admin123` | Success | ✓ PASS |
| Leading space in username | ` admin` | `admin123` | Success (trimmed) | ✓ PASS |
| Trailing space in username | `admin ` | `admin123` | Success (trimmed) | ✓ PASS |
| Both spaces in username | ` admin ` | `admin123` | Success (trimmed) | ✓ PASS |
| Wrong password | `admin` | `wrongpass` | Reject | ✓ PASS |
| Non-existent user | `hacker` | `admin123` | Reject | ✓ PASS |

## Security Considerations

- ✅ No security vulnerabilities introduced (verified with CodeQL)
- ✅ Password integrity preserved (not trimmed)
- ✅ Username validation still works correctly
- ✅ Invalid credentials still properly rejected

## Impact

- **User Experience**: Significantly improved - users can successfully login even if they accidentally copy/paste spaces
- **Security**: No negative impact - maintains same security level
- **Backward Compatibility**: Fully backward compatible - existing users unaffected
- **Code Complexity**: Minimal - only 5 lines added

## Related Issues

This fix addresses a common UX issue in web applications. It complements the previous CORS middleware fix (PR #12) to provide a more robust login experience.

## Version

- **Fixed in**: This PR
- **Affects**: All previous versions
- **Breaking Changes**: None
