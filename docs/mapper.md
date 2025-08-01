# Mapper Package Documentation

## Package: `mapper`

The `mapper` package provides utility functions to convert user data across three application layers:

* **Database Models (`models`)**: Structs used for MongoDB operations.
* **Domain Entities (`entities`)**: Core business logic representations.
* **DTOs (`dto`)**: Data Transfer Objects for API responses.

This mapping ensures **separation of concerns**, **clean architecture**, and **security** (e.g., hiding sensitive fields in DTOs).

---

## Function Overview

### `ModelToEntity`

```go
func ModelToEntity(model models.User) entities.User
```

**Purpose:**
Converts a MongoDB user model (`models.User`) into a domain-level entity (`entities.User`).

**Use Case:**
After retrieving user data from MongoDB, use this function to pass the result to business logic layers.

**Key Transformations:**

* `ObjectID` fields like `ID`, `UpdatedBy`, and `UserPosts` are converted to strings.
* Social links and profile data are deeply mapped.

---

### `EntityToModel`

```go
func EntityToModel(ue entities.User) (models.User, error)
```

**Purpose:**
Converts a domain entity (`entities.User`) back into a MongoDB-compatible model (`models.User`).

**Use Case:**
When saving or updating user data into the database.

**Returns:**

* `models.User`: Struct ready for MongoDB operations.
* `error`: Returns error if `ID`, `UpdatedBy`, or any `UserPosts` fail to convert from hex string to `ObjectID`.

**Key Transformations:**

* Converts all hex string IDs to `primitive.ObjectID`.
* Maps nested profile and social link structures.

---

### `EntityToDTO`

```go
func EntityToDTO(ue entities.User) dto.UserDTO
```

**Purpose:**
Converts a domain user entity (`entities.User`) into a DTO (`dto.UserDTO`) for safe API responses.

**Use Case:**
Preparing sanitized and formatted user data to return to clients via HTTP.

**Security Note:**

* Does **not** expose sensitive fields like `PasswordHash`.

**Key Transformations:**

* Nested fields like `Profile` and `SocialLinks` are mapped to their DTO counterparts.

---

## Conversion Architecture

| Layer        | Struct Name     | Description                      |
| ------------ | --------------- | -------------------------------- |
| Database     | `models.User`   | MongoDB model                    |
| Domain Logic | `entities.User` | Business logic representation    |
| Transport    | `dto.UserDTO`   | API-facing user data (sanitized) |

---

## Nested Structure Handling

Each function handles deeply nested structures:

* `UserProfile`, `UserProfileDTO`, and `UserProfile (Model)`
* `SocialLink`, `SocialLinkDTO`, and `SocialLink (Model)`

This ensures consistency across all layers.

---

## Example Usage

### 1. Converting DB User to API Response

```go
dbUser, _ := db.GetUserByID(id)
domainUser := mapper.ModelToEntity(dbUser)
responseDTO := mapper.EntityToDTO(domainUser)
// Return responseDTO as API JSON
```

---

### 2. Converting API Input to DB Model

```go
// Assume inputDTO is converted to domainUser beforehand
dbModel, err := mapper.EntityToModel(domainUser)
if err != nil {
	log.Fatal("Invalid user data")
}
db.Create(dbModel)
```

---

## Best Practices

* Always validate input before converting to `models.User`.
* Use `EntityToDTO` to ensure sensitive information is excluded from API responses.
* Handle `ObjectID` conversion errors gracefully in `EntityToModel`.

---

## Related Modules

| Package    | Description                                 |
| ---------- | ------------------------------------------- |
| `models`   | MongoDB data models                         |
| `entities` | Business logic domain structs               |
| `dto`      | Data Transfer Objects for external response |