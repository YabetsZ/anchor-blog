# Forgot Password Flow

This document describes the password reset system for users who have forgotten their passwords.

## Overview

Users can request a password reset by providing their email address. They receive a reset email with a unique token that allows them to set a new password.

## Components

### 1. PasswordResetService (`internal/service/user/password_reset_service.go`)

**ForgotPassword(email)**
- Validates email format and existence
- Generates a unique, secure reset token (64 hex characters)
- Creates password reset token with 1-hour expiry
- Logs reset link to console (for development)
- Format: `http://localhost:8080/api/v1/users/reset-password?token=<token>`

**ResetPassword(token, newPassword)**
- Validates the reset token and new password
- Hashes the new password using bcrypt
- Updates user's password in database
- Marks token as used to prevent reuse

### 2. PasswordResetHandler (`api/handler/password_reset_handler.go`)

**POST /api/v1/users/forgot-password**
- Accepts email in request body
- Calls ForgotPassword service
- Returns success confirmation

**POST /api/v1/users/reset-password**
- Accepts token and new password in request body
- Calls ResetPassword service
- Returns success confirmation with user info

## API Endpoints

### Request Password Reset
```
POST /api/v1/users/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**Success Response (200):**
```json
{
  "message": "Password reset email sent successfully",
  "email": "user@example.com"
}
```

**Error Response (400):**
```json
{
  "error": "Invalid request format",
  "details": "email is required"
}
```

### Reset Password
```
POST /api/v1/users/reset-password
Content-Type: application/json

{
  "token": "abc123...",
  "new_password": "newSecurePassword123"
}
```

**Success Response (200):**
```json
{
  "message": "Password reset successfully",
  "user": {
    "id": "user-id",
    "username": "username",
    "email": "user@example.com"
  }
}
```

**Error Response (400):**
```json
{
  "error": "Failed to reset password",
  "details": "Invalid or expired reset token"
}
```

## Usage Flow

1. **Forgot Password**: User submits email → `POST /api/v1/users/forgot-password`
2. **Email**: Reset link is logged to console
3. **Reset**: User clicks link and submits new password → `POST /api/v1/users/reset-password`
4. **Verification**: Token is validated and password is updated
5. **Access**: User can now log in with new password

## Security Features

- **Unique tokens**: Cryptographically secure random tokens
- **Time-limited**: Tokens expire after 1 hour (shorter than activation)
- **Single-use**: Tokens are marked as used after password reset
- **Password hashing**: New passwords are hashed with bcrypt
- **Validation**: Email format and password strength validation
- **Minimum password length**: 6 characters required

## Validation Rules

### Email Validation
- Required field
- Must be valid email format
- Must exist in database (TODO: implement lookup)

### Password Validation
- Required field
- Minimum 6 characters
- Hashed with bcrypt before storage

## Development Notes

- Reset links are currently logged to console
- Database integration is marked with TODO comments
- Ready for email service integration
- Tokens are 64 hex characters (256-bit security)
- Shorter expiry time (1 hour) for security

## Integration Example

Frontend can use these endpoints like this:

```javascript
// Step 1: Request password reset
const forgotResponse = await fetch('/api/v1/users/forgot-password', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email: 'user@example.com' })
});

// Step 2: Reset password with token
const resetResponse = await fetch('/api/v1/users/reset-password', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ 
    token: 'token-from-email',
    new_password: 'newPassword123'
  })
});
```

## Error Handling

The system handles various error cases:
- Missing or invalid email format
- Missing or invalid token
- Expired tokens
- Already used tokens
- Weak passwords (< 6 characters)
- Database connection issues
- Password hashing failures