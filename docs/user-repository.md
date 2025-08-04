# 📦 `userrepo` Package Documentation

## Overview

The `userrepo` package provides a MongoDB-based implementation of the `IUserRepository` interface defined in the `entities` package. It supports all necessary CRUD operations and business logic required to manage user data in a scalable and maintainable way.

This implementation adheres to **Clean Architecture** principles by separating domain logic (`entities`) from infrastructure logic (`userrepo`), ensuring testability and flexibility.

---

## 🧩 Interface Implementations

### 🔍 Read Operations (`IUserReaderRepository`)

#### `GetUserByID(ctx, id)`

Fetches a user by MongoDB ObjectID. Returns a domain `*User`.

#### `GetUserByEmail(ctx, email)`

Returns a `*User` by email. Useful for authentication or validation.

#### `GetUserByUsername(ctx, username)`

Looks up user by username. Validates uniqueness during registration.

#### `GetUsers(ctx, limit, offset)`

Paginates through users using MongoDB's `Find()` with `limit` and `skip`.

#### `CountUsersByRole(ctx, role)`

Returns the number of users for a specific role (e.g., "admin", "user").

#### `CountAllUsers(ctx)`

Returns the total number of users in the system.

#### `CountActiveUsers(ctx)`

Returns count of users where `activated == true`.

#### `CountInactiveUsers(ctx)`

Returns count of users where `activated == false`.

#### `GetUserRoleByID(ctx, id)`

Returns the role string of a user by ID.

---

### 🛠 Write Operations (`IUserWriterRepository`)

#### `CreateUser(ctx, user)`

Converts a domain-level `User` into a MongoDB document and inserts it.

#### `EditUserByID(ctx, id, user)`

Updates user fields only if non-zero values are provided. Uses `$set` to apply changes.

#### `DeleteUserByID(ctx, id)`

Removes a user document by its MongoDB ID.

#### `SetLastSeen(ctx, id, timestamp)`

Updates a user's `LastSeen` timestamp field.

---

### 🔐 Auth Operations (`IUserAuthRepository`)

#### `CheckEmail(ctx, email)`

Checks if an email is already registered. Returns `true` if exists.

#### `CheckUsername(ctx, username)`

Checks if a username is already taken. Returns `true` if exists.

#### `ChangePassword(ctx, id, newHashedPassword)`

Updates a user's `PasswordHash` field.

#### `ChangeEmail(ctx, email, newEmail)`

Replaces an existing email address with a new one.

---

### 🛡 Admin Operations (`IUserAdminRepository`)

#### `SetRole(ctx, id, role)`

Changes a user's role. E.g., promoting a "user" to "admin".

#### `ActivateUserByID(ctx, id)`

Marks `activated = true` for a user.

#### `DeactivateUserByID(ctx, id)`

Marks `activated = false` for a user.

---

## 🧱 MongoDB Schema Assumptions

Each `User` document is expected to have:

* `_id` as `ObjectID`
* Fields: `username`, `email`, `passwordHash`, `role`, `activated`, `lastSeen`
* Embedded `profile` with:

  * `bio`
  * `pictureURL`
  * `socialLinks` (slice of `{platform, url}`)

---

## ⚙️ Utility Functions

### `IsUserProfileEmpty(profile UserProfile)`

Checks if the `UserProfile` struct contains any meaningful data. Prevents overwriting with empty values during updates.

---

## 🧪 Error Handling

All errors are logged using `log.Printf` and returned as centralized custom errors from the `internal/errors` package (e.g., `ErrInternalServer`), ensuring consistency in error response and debugging.

---