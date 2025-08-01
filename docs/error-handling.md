# ❗ Error Handling – `internal/errors`

This module provides centralized, domain-specific error definitions for consistent and type-safe error handling across the application.

## 📦 Defined Errors

### `ErrUserNotFound`

* **Type:** `error`
* **Description:** Returned when a requested user does not exist in the database.

### `ErrEmailAlreadyExists`

* **Type:** `error`
* **Description:** Raised when attempting to register or update with an email that is already in use.

### `ErrUsernameTaken`

* **Type:** `error`
* **Description:** Indicates the username is already taken by another user.

### `ErrInvalidCredentials`

* **Type:** `error`
* **Description:** Returned on failed login due to incorrect email or password.

### `ErrUnauthorized`

* **Type:** `error`
* **Description:** Returned when an action is attempted without proper authentication.

### `ErrForbidden`

* **Type:** `error`
* **Description:** Indicates the user is authenticated but not authorized to perform the requested action.

### `ErrValidationFailed`

* **Type:** `error`
* **Description:** Used when input validation fails (e.g., missing required fields, invalid format).

### `ErrInternalServer`

* **Type:** `error`
* **Description:** Represents a generic internal server error, typically used as a fallback for unexpected failures.