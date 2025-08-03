# Account Activation Flow

This document describes the account activation system implemented for user registration.

## Overview

When a user registers, they receive an activation email with a unique token. They must click the activation link to activate their account and gain full access.

## Components

### 1. ActivationService (`internal/service/user/activation_service.go`)

**SendActivationEmail(user)**
- Generates a unique, secure activation token (64 hex characters)
- Creates activation token with 24-hour expiry
- Logs activation link to console (for development)
- Format: `http://localhost:8080/api/v1/users/activate?token=<token>`

**VerifyActivation(token)**
- Validates the activation token
- Sets user `activated = true` and `role = "user"`
- Marks token as used to prevent reuse

### 2. ActivationHandler (`api/handler/activation_handler.go`)

**GET /api/v1/users/activate**
- Extracts token from query parameter
- Calls VerifyActivation service
- Returns success/error response

## API Endpoint

### Activate Account
```
GET /api/v1/users/activate?token=<activation_token>
```

**Success Response (200):**
```json
{
  "message": "Account activated successfully",
  "user": {
    "id": "user-id",
    "username": "username",
    "email": "user@example.com",
    "activated": true,
    "role": "user"
  }
}
```

**Error Response (400):**
```json
{
  "error": "Invalid or expired activation token",
  "details": "specific error message"
}
```

## Usage Flow

1. **Registration**: User registers → `SendActivationEmail()` is called
2. **Email**: Activation link is logged to console
3. **Activation**: User clicks link → `GET /api/v1/users/activate?token=...`
4. **Verification**: Token is validated and user is activated
5. **Access**: User can now log in with full permissions

## Security Features

- **Unique tokens**: Cryptographically secure random tokens
- **Time-limited**: Tokens expire after 24 hours
- **Single-use**: Tokens are marked as used after activation
- **Validation**: Comprehensive token and user validation

## Development Notes

- Activation links are currently logged to console
- Database integration is marked with TODO comments
- Ready for email service integration
- Tokens are 64 hex characters (256-bit security)

## Integration with Registration

The registration service (implemented by Dev A) should call:

```go
activationService := userservice.NewActivationService()
err := activationService.SendActivationEmail(newUser)
```

This ensures new users receive activation emails immediately after registration.