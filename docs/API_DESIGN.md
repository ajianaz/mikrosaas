# MikroSaaS API Design Document

> **Base URL:** `/api/v1`
> **Protocol:** HTTPS (required in production)
> **Auth:** JWT Bearer tokens
> **Multi-tenancy:** tenant_id extracted from JWT claims; all resources scoped to tenant
> **Content-Type:** `application/json` for all request/response bodies

---

## Table of Contents

1. [Conventions](#conventions)
2. [Authentication](#authentication)
3. [Error Handling](#error-handling)
4. [Pagination](#pagination)
5. [Endpoints](#endpoints)
   - [Auth](#auth)
   - [Users](#users)
   - [Tenants](#tenants)
   - [Roles](#roles)
   - [RBAC / Permissions](#rbac--permissions)
6. [Error Codes Reference](#error-codes-reference)
7. [Rate Limiting](#rate-limiting)
8. [CORS](#cors)
9. [Future Considerations](#future-considerations)

---

## Conventions

### UUIDs

All primary keys are UUID v4 strings. Example: `a1b2c3d4-e5f6-7890-abcd-ef1234567890`.

### Timestamps

All timestamps are ISO 8601 UTC strings. Example: `2025-01-15T10:30:00Z`.

### Partial Updates (PATCH semantics via PUT)

`UpdateRequest` bodies use nullable/pointer fields. Only provided fields are applied. Omitting a field means "do not change."

### Soft Deletes

Resources use `is_active` flag. `DELETE` sets `is_active = false`. No hard deletes via API.

### Tenancy

Authenticated endpoints extract `tenant_id` from JWT claims. The frontend never sends `tenant_id` — the server enforces scope.

---

## Authentication

### JWT Structure

```
Authorization: Bearer <access_token>
```

**Access token payload:**

```json
{
  "sub": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
  "email": "jane@example.com",
  "roles": ["admin"],
  "exp": 1705312200,
  "iat": 1705308600
}
```

| Claim | Type | Description |
|-------|------|-------------|
| `sub` | string | User ID (UUID) |
| `tenant_id` | string | Tenant ID (UUID) — all requests scoped here |
| `email` | string | User email |
| `roles` | string[] | Role names assigned to user within tenant |
| `exp` | int64 | Expiration (Unix timestamp) |
| `iat` | int64 | Issued at (Unix timestamp) |

**Token lifetimes:**

| Token | Lifetime |
|-------|----------|
| Access | 15 minutes |
| Refresh | 7 days |

---

## Error Handling

### Standard Error Response

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "must be a valid email address"
      },
      {
        "field": "password",
        "message": "must be at least 8 characters"
      }
    ]
  }
}
```

### Standard Success Response

Single resource:

```json
{
  "data": { ... }
}
```

List resource:

```json
{
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 47,
    "total_pages": 3
  }
}
```

---

## Pagination

**Query parameters:**

| Param | Type | Default | Max | Description |
|-------|------|---------|-----|-------------|
| `page` | int | 1 | — | Page number (1-indexed) |
| `per_page` | int | 20 | 100 | Items per page |

**Response meta (list endpoints only):**

```json
{
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 47,
    "total_pages": 3
  }
}
```

`total_pages` = `ceil(total / per_page)`.

---

## Endpoints

---

### Auth

#### `POST /api/v1/auth/login`

Authenticate with email + password. Returns tokens + user.

**Auth:** Public

**Request:**

```json
{
  "email": "jane@example.com",
  "password": "s3cureP@ss!"
}
```

**Response `200 OK`:**

```json
{
  "data": {
    "user": {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "email": "jane@example.com",
      "first_name": "Jane",
      "last_name": "Smith",
      "is_active": true,
      "email_verified": true,
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:30:00Z"
    },
    "token": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g...",
      "expires_in": 900,
      "token_type": "Bearer"
    }
  }
}
```

**Error responses:**

`401 Unauthorized` — invalid credentials:

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid email or password"
  }
}
```

`403 Forbidden` — user deactivated:

```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Account is deactivated"
  }
}
```

---

#### `POST /api/v1/auth/refresh`

Exchange a valid refresh token for a new token pair.

**Auth:** Public (refresh token acts as credential)

**Request:**

```json
{
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g..."
}
```

**Response `200 OK`:**

```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "bmV3IHJlZnJlc2ggdG9rZW4...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

**Error responses:**

`401 Unauthorized` — invalid/expired refresh token:

```json
{
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Refresh token is invalid or expired"
  }
}
```

`401 Unauthorized` — refresh token revoked:

```json
{
  "error": {
    "code": "TOKEN_EXPIRED",
    "message": "Refresh token has been revoked"
  }
}
```

---

#### `POST /api/v1/auth/logout`

Invalidate the refresh token. Future endpoint — currently a no-op returning 204.

**Auth:** Bearer

**Request:** No body required. Token extracted from `Authorization` header.

**Response `204 No Content`:**

_(empty body)_

---

#### `GET /api/v1/auth/me`

Get the authenticated user's profile.

**Auth:** Bearer

**Response `200 OK`:**

```json
{
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "email": "jane@example.com",
    "first_name": "Jane",
    "last_name": "Smith",
    "is_active": true,
    "email_verified": true,
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

**Error responses:**

`401 Unauthorized`:

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing or invalid access token"
  }
}
```

---

### Users

All user endpoints require Bearer auth. Resources are scoped to the tenant from JWT claims.

#### `GET /api/v1/users`

List users for the authenticated tenant.

**Auth:** Bearer

**Query parameters:**

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `per_page` | int | 20 | Items per page (max 100) |
| `search` | string | — | Fuzzy search on email, first_name, last_name |
| `is_active` | bool | — | Filter by active status |

**Example:** `GET /api/v1/users?page=2&per_page=10&search=jane&is_active=true`

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "email": "jane@example.com",
      "first_name": "Jane",
      "last_name": "Smith",
      "is_active": true,
      "email_verified": true,
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:30:00Z"
    },
    {
      "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "email": "jane.doe@example.com",
      "first_name": "Jane",
      "last_name": "Doe",
      "is_active": true,
      "email_verified": false,
      "created_at": "2025-02-01T08:00:00Z",
      "updated_at": "2025-02-01T08:00:00Z"
    }
  ],
  "meta": {
    "page": 2,
    "per_page": 10,
    "total": 47,
    "total_pages": 5
  }
}
```

---

#### `POST /api/v1/users`

Create a new user within the authenticated tenant.

**Auth:** Bearer

**Request:**

```json
{
  "email": "newuser@example.com",
  "password": "s3cureP@ss!",
  "first_name": "Alex",
  "last_name": "Johnson"
}
```

**Response `201 Created`:**

```json
{
  "data": {
    "id": "d4e5f6a7-b8c9-0123-defa-234567890123",
    "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "email": "newuser@example.com",
    "first_name": "Alex",
    "last_name": "Johnson",
    "is_active": true,
    "email_verified": false,
    "created_at": "2025-06-01T14:00:00Z",
    "updated_at": "2025-06-01T14:00:00Z"
  }
}
```

**Error responses:**

`400 Bad Request` — validation failure:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      { "field": "email", "message": "must be a valid email address" },
      { "field": "first_name", "message": "is required" }
    ]
  }
}
```

`409 Conflict` — email already exists in tenant:

```json
{
  "error": {
    "code": "USER_ALREADY_EXISTS",
    "message": "A user with this email already exists in this tenant"
  }
}
```

---

#### `GET /api/v1/users/:id`

Get a single user by ID.

**Auth:** Bearer

**URL params:** `id` — User UUID

**Response `200 OK`:**

```json
{
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "email": "jane@example.com",
    "first_name": "Jane",
    "last_name": "Smith",
    "is_active": true,
    "email_verified": true,
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

**Error responses:**

`404 Not Found`:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found"
  }
}
```

---

#### `PUT /api/v1/users/:id`

Update a user. Partial — only provided fields are modified.

**Auth:** Bearer

**URL params:** `id` — User UUID

**Request** (all fields optional):

```json
{
  "first_name": "Janet",
  "last_name": "Smith-Johnson",
  "is_active": true
}
```

**Response `200 OK`:**

```json
{
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "email": "jane@example.com",
    "first_name": "Janet",
    "last_name": "Smith-Johnson",
    "is_active": true,
    "email_verified": true,
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-06-01T15:30:00Z"
  }
}
```

**Error responses:**

`400 Bad Request` — validation:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      { "field": "first_name", "message": "must be at most 100 characters" }
    ]
  }
}
```

`404 Not Found`:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found"
  }
}
```

---

#### `DELETE /api/v1/users/:id`

Soft-delete: sets `is_active = false`. User cannot log in after deactivation.

**Auth:** Bearer

**URL params:** `id` — User UUID

**Request:** No body.

**Response `200 OK`:**

```json
{
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "email": "jane@example.com",
    "first_name": "Jane",
    "last_name": "Smith",
    "is_active": false,
    "email_verified": true,
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-06-01T16:00:00Z"
  }
}
```

**Error responses:**

`404 Not Found`:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found"
  }
}
```

---

### Tenants

All tenant endpoints require Bearer auth + admin role.

#### `GET /api/v1/tenants`

List tenants. Admin-only.

**Auth:** Bearer + Admin

**Query parameters:**

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `per_page` | int | 20 | Items per page (max 100) |

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "name": "Acme Corp",
      "slug": "acme-corp",
      "domain": "acme.app.mikrosaas.io",
      "is_active": true,
      "created_at": "2025-01-10T09:00:00Z",
      "updated_at": "2025-03-20T11:00:00Z"
    },
    {
      "id": "e5f6a7b8-c9d0-1234-efab-567890123456",
      "name": "Globex Inc",
      "slug": "globex-inc",
      "domain": null,
      "is_active": true,
      "created_at": "2025-02-05T14:00:00Z",
      "updated_at": "2025-02-05T14:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 12,
    "total_pages": 1
  }
}
```

---

#### `POST /api/v1/tenants`

Create a new tenant. Admin-only.

**Auth:** Bearer + Admin

**Request:**

```json
{
  "name": "Initech",
  "slug": "initech",
  "domain": "initech.app.mikrosaas.io"
}
```

**Response `201 Created`:**

```json
{
  "data": {
    "id": "f6a7b8c9-d0e1-2345-fabc-678901234567",
    "name": "Initech",
    "slug": "initech",
    "domain": "initech.app.mikrosaas.io",
    "is_active": true,
    "created_at": "2025-06-01T10:00:00Z",
    "updated_at": "2025-06-01T10:00:00Z"
  }
}
```

**Error responses:**

`400 Bad Request` — validation:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      { "field": "slug", "message": "must match pattern: lowercase alphanumeric and hyphens, 1-100 chars" }
    ]
  }
}
```

`409 Conflict` — slug already exists:

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "Tenant with this slug already exists"
  }
}
```

---

#### `GET /api/v1/tenants/:id`

Get a single tenant by ID.

**Auth:** Bearer + Admin

**URL params:** `id` — Tenant UUID

**Response `200 OK`:**

```json
{
  "data": {
    "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "name": "Acme Corp",
    "slug": "acme-corp",
    "domain": "acme.app.mikrosaas.io",
    "is_active": true,
    "created_at": "2025-01-10T09:00:00Z",
    "updated_at": "2025-03-20T11:00:00Z"
  }
}
```

**Error responses:**

`404 Not Found`:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Tenant not found"
  }
}
```

---

#### `PUT /api/v1/tenants/:id`

Update a tenant. Partial — only provided fields are modified. Admin-only.

**Auth:** Bearer + Admin

**URL params:** `id` — Tenant UUID

**Request** (all fields optional):

```json
{
  "name": "Acme Corporation",
  "domain": "acme.mikrosaas.io"
}
```

**Response `200 OK`:**

```json
{
  "data": {
    "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "name": "Acme Corporation",
    "slug": "acme-corp",
    "domain": "acme.mikrosaas.io",
    "is_active": true,
    "created_at": "2025-01-10T09:00:00Z",
    "updated_at": "2025-06-01T12:00:00Z"
  }
}
```

**Error responses:**

`400 Bad Request` — validation:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      { "field": "name", "message": "must be between 1 and 200 characters" }
    ]
  }
}
```

`404 Not Found`:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Tenant not found"
  }
}
```

`409 Conflict` — slug already taken:

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "Tenant with this slug already exists"
  }
}
```

---

#### `DELETE /api/v1/tenants/:id`

Soft-delete: sets `is_active = false`. All users in tenant lose access. Admin-only.

**Auth:** Bearer + Admin

**URL params:** `id` — Tenant UUID

**Request:** No body.

**Response `200 OK`:**

```json
{
  "data": {
    "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "name": "Acme Corp",
    "slug": "acme-corp",
    "domain": "acme.app.mikrosaas.io",
    "is_active": false,
    "created_at": "2025-01-10T09:00:00Z",
    "updated_at": "2025-06-01T13:00:00Z"
  }
}
```

**Error responses:**

`404 Not Found`:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Tenant not found"
  }
}
```

---

### Roles

All role endpoints require Bearer auth. Resources scoped to tenant.

#### `GET /api/v1/roles`

List roles for the authenticated tenant.

**Auth:** Bearer

**Query parameters:**

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `per_page` | int | 20 | Items per page (max 100) |

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "name": "admin",
      "description": "Full administrative access",
      "is_system": true,
      "created_at": "2025-01-10T09:00:00Z",
      "updated_at": "2025-01-10T09:00:00Z"
    },
    {
      "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "name": "editor",
      "description": "Can read and write resources",
      "is_system": false,
      "created_at": "2025-02-01T08:00:00Z",
      "updated_at": "2025-02-01T08:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 3,
    "total_pages": 1
  }
}
```

---

#### `POST /api/v1/roles`

Create a new role within the authenticated tenant.

**Auth:** Bearer

**Request:**

```json
{
  "name": "viewer",
  "description": "Read-only access to all resources"
}
```

**Response `201 Created`:**

```json
{
  "data": {
    "id": "d4e5f6a7-b8c9-0123-defa-234567890123",
    "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "name": "viewer",
    "description": "Read-only access to all resources",
    "is_system": false,
    "created_at": "2025-06-01T14:00:00Z",
    "updated_at": "2025-06-01T14:00:00Z"
  }
}
```

**Error responses:**

`400 Bad Request` — validation:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      { "field": "name", "message": "is required" }
    ]
  }
}
```

`409 Conflict` — role name exists in tenant:

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "A role with this name already exists in this tenant"
  }
}
```

---

#### `GET /api/v1/roles/:id`

Get a single role by ID.

**Auth:** Bearer

**URL params:** `id` — Role UUID

**Response `200 OK`:**

```json
{
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "name": "admin",
    "description": "Full administrative access",
    "is_system": true,
    "created_at": "2025-01-10T09:00:00Z",
    "updated_at": "2025-01-10T09:00:00Z"
  }
}
```

**Error responses:**

`404 Not Found`:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Role not found"
  }
}
```

---

#### `PUT /api/v1/roles/:id`

Update a role. System roles can have their description updated but not their name. Admin-only for system roles.

**Auth:** Bearer

**URL params:** `id` — Role UUID

**Request** (all fields optional):

```json
{
  "name": "senior-editor",
  "description": "Can read, write, and publish resources"
}
```

**Response `200 OK`:**

```json
{
  "data": {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "name": "senior-editor",
    "description": "Can read, write, and publish resources",
    "is_system": false,
    "created_at": "2025-02-01T08:00:00Z",
    "updated_at": "2025-06-01T15:00:00Z"
  }
}
```

**Error responses:**

`400 Bad Request` — cannot modify system role name:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Cannot modify system role name"
  }
}
```

`404 Not Found`:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Role not found"
  }
}
```

`409 Conflict` — role name exists in tenant:

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "A role with this name already exists in this tenant"
  }
}
```

---

#### `DELETE /api/v1/roles/:id`

Delete a role. System roles (`is_system = true`) cannot be deleted. Deleting a role removes all user-role and role-permission mappings.

**Auth:** Bearer

**URL params:** `id` — Role UUID

**Request:** No body.

**Response `200 OK`:**

```json
{
  "data": {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "name": "editor",
    "description": "Can read and write resources",
    "is_system": false,
    "created_at": "2025-02-01T08:00:00Z",
    "updated_at": "2025-02-01T08:00:00Z"
  }
}
```

**Error responses:**

`403 Forbidden` — system role:

```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "System roles cannot be deleted"
  }
}
```

`404 Not Found`:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Role not found"
  }
}
```

---

### RBAC / Permissions

#### `GET /api/v1/roles/:id/permissions`

List all permissions assigned to a role.

**Auth:** Bearer

**URL params:** `id` — Role UUID

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": "p1a2b3c4-d5e6-7890-abcd-ef1234567890",
      "name": "users:read",
      "resource": "users",
      "action": "read",
      "description": "View user profiles",
      "created_at": "2025-01-10T09:00:00Z"
    },
    {
      "id": "p2b3c4d5-e6f7-8901-bcde-f123456789012",
      "name": "users:write",
      "resource": "users",
      "action": "write",
      "description": "Create and update users",
      "created_at": "2025-01-10T09:00:00Z"
    },
    {
      "id": "p3c4d5e6-f7a8-9012-cdef-123456789012",
      "name": "roles:read",
      "resource": "roles",
      "action": "read",
      "description": "View roles",
      "created_at": "2025-01-10T09:00:00Z"
    }
  ]
}
```

**Error responses:**

`404 Not Found` — role not found:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Role not found"
  }
}
```

---

#### `PUT /api/v1/roles/:id/permissions`

Replace all permissions on a role. Sends the complete set — previous assignments are removed. Admin-only.

**Auth:** Bearer + Admin

**URL params:** `id` — Role UUID

**Request:**

```json
{
  "permission_ids": [
    "p1a2b3c4-d5e6-7890-abcd-ef1234567890",
    "p2b3c4d5-e6f7-8901-bcde-f123456789012",
    "p3c4d5e6-f7a8-9012-cdef-123456789012"
  ]
}
```

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": "p1a2b3c4-d5e6-7890-abcd-ef1234567890",
      "name": "users:read",
      "resource": "users",
      "action": "read",
      "description": "View user profiles",
      "created_at": "2025-01-10T09:00:00Z"
    },
    {
      "id": "p2b3c4d5-e6f7-8901-bcde-f123456789012",
      "name": "users:write",
      "resource": "users",
      "action": "write",
      "description": "Create and update users",
      "created_at": "2025-01-10T09:00:00Z"
    },
    {
      "id": "p3c4d5e6-f7a8-9012-cdef-123456789012",
      "name": "roles:read",
      "resource": "roles",
      "action": "read",
      "description": "View roles",
      "created_at": "2025-01-10T09:00:00Z"
    }
  ]
}
```

**Error responses:**

`400 Bad Request` — empty permission list:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "At least one permission is required"
  }
}
```

`404 Not Found` — role not found:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Role not found"
  }
}
```

`422 Unprocessable Entity` — invalid permission ID:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "One or more permission IDs are invalid",
    "details": [
      { "field": "permission_ids[2]", "message": "permission not found" }
    ]
  }
}
```

---

#### `PUT /api/v1/users/:id/roles`

Replace all roles assigned to a user. Sends the complete set — previous assignments are removed. Admin-only.

**Auth:** Bearer + Admin

**URL params:** `id` — User UUID

**Request:**

```json
{
  "role_ids": [
    "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "c3d4e5f6-a7b8-9012-cdef-123456789012"
  ]
}
```

**Response `200 OK`:**

```json
{
  "data": {
    "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "roles": [
      {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
        "name": "admin",
        "description": "Full administrative access",
        "is_system": true,
        "created_at": "2025-01-10T09:00:00Z",
        "updated_at": "2025-01-10T09:00:00Z"
      },
      {
        "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
        "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
        "name": "editor",
        "description": "Can read and write resources",
        "is_system": false,
        "created_at": "2025-02-01T08:00:00Z",
        "updated_at": "2025-02-01T08:00:00Z"
      }
    ]
  }
}
```

**Error responses:**

`400 Bad Request` — empty role list:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "At least one role is required"
  }
}
```

`404 Not Found` — user not found:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found"
  }
}
```

`422 Unprocessable Entity` — invalid role ID or cross-tenant role:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "One or more role IDs are invalid",
    "details": [
      { "field": "role_ids[1]", "message": "role not found or belongs to different tenant" }
    ]
  }
}
```

---

#### `GET /api/v1/permissions`

List all available permissions (global, not tenant-scoped).

**Auth:** Bearer

**Query parameters:**

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `resource` | string | — | Filter by resource (e.g. `users`, `roles`, `tenants`) |
| `page` | int | 1 | Page number |
| `per_page` | int | 100 | Items per page (default 100 for permissions) |

**Example:** `GET /api/v1/permissions?resource=users`

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": "p1a2b3c4-d5e6-7890-abcd-ef1234567890",
      "name": "users:read",
      "resource": "users",
      "action": "read",
      "description": "View user profiles",
      "created_at": "2025-01-10T09:00:00Z"
    },
    {
      "id": "p2b3c4d5-e6f7-8901-bcde-f123456789012",
      "name": "users:write",
      "resource": "users",
      "action": "write",
      "description": "Create and update users",
      "created_at": "2025-01-10T09:00:00Z"
    },
    {
      "id": "p4d5e6f7-a8b9-0123-defa-345678901234",
      "name": "users:delete",
      "resource": "users",
      "action": "delete",
      "description": "Deactivate users",
      "created_at": "2025-01-10T09:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 100,
    "total": 3,
    "total_pages": 1
  }
}
```

---

## Error Codes Reference

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 / 422 | Request body or query params failed validation. `details` array lists field-level errors. |
| `UNAUTHORIZED` | 401 | Missing or invalid access token. |
| `INVALID_TOKEN` | 401 | Token is malformed or signature verification failed. |
| `TOKEN_EXPIRED` | 401 | Token has expired. Client should attempt refresh. |
| `FORBIDDEN` | 403 | Authenticated but lacks permission for this action (wrong role, resource owner mismatch, system role protection). |
| `NOT_FOUND` | 404 | Requested resource does not exist in the current tenant scope. |
| `CONFLICT` | 409 | Resource already exists (duplicate email in tenant, duplicate tenant slug, duplicate role name). |
| `USER_ALREADY_EXISTS` | 409 | User with this email already exists in the target tenant. |
| `TENANT_NOT_FOUND` | 404 | Referenced tenant does not exist (e.g. from JWT claims after tenant deletion). |
| `INTERNAL_ERROR` | 500 | Unexpected server error. Client should retry with backoff. |

### HTTP Status Code Summary

| Status | Meaning |
|--------|---------|
| `200 OK` | Successful read or update |
| `201 Created` | Successful resource creation |
| `204 No Content` | Successful deletion (logout) |
| `400 Bad Request` | Malformed request syntax |
| `401 Unauthorized` | Missing, invalid, or expired token |
| `403 Forbidden` | Valid token but insufficient permissions |
| `404 Not Found` | Resource does not exist |
| `409 Conflict` | Duplicate resource |
| `422 Unprocessable Entity` | Valid JSON but semantic validation failure |
| `500 Internal Server Error` | Server-side failure |

---

## Rate Limiting

**Strategy:**

| Scope | Endpoint Group | Limit | Window |
|-------|---------------|-------|--------|
| Per IP | `/api/v1/auth/*` | 10 requests | 1 minute |
| Per IP | `/api/v1/auth/login` | 5 attempts | 15 minutes (after 3 failures, exponential backoff) |
| Per Tenant | All other endpoints | 100 requests | 1 minute |
| Per Tenant | All other endpoints | 1,000 requests | 1 hour |

**Response headers:**

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 87
X-RateLimit-Reset: 1705312200
```

**When rate limit is exceeded:**

```json
{
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many requests. Retry after 45 seconds.",
    "details": [
      { "field": "retry_after", "message": "45" }
    ]
  }
}
```

HTTP `429 Too Many Requests`.

**Implementation plan:** Redis-backed sliding window counters. Per-IP keys: `rl:ip:{ip}:{endpoint}`. Per-tenant keys: `rl:tenant:{tenant_id}`. Auth brute-force protection: lockout after 5 failures, 15-min cooldown with `rl:auth:fail:{email}`.

---

## CORS

**Configuration approach:** Environment-driven via config struct.

```go
type CORSConfig struct {
    AllowedOrigins   []string `env:"CORS_ALLOWED_ORIGINS"`   // e.g. ["https://app.example.com"]
    AllowedMethods   []string `env:"CORS_ALLOWED_METHODS"`   // default: ["GET","POST","PUT","DELETE","OPTIONS"]
    AllowedHeaders   []string `env:"CORS_ALLOWED_HEADERS"`   // default: ["Authorization","Content-Type","X-Request-ID"]
    ExposedHeaders   []string `env:"CORS_EXPOSED_HEADERS"`   // default: ["X-RateLimit-Limit","X-RateLimit-Remaining","X-RateLimit-Reset"]
    AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS"` // true for cookie-based; false for bearer-only
    MaxAge           int      `env:"CORS_MAX_AGE"`           // default: 86400 (24h preflight cache)
}
```

**Default behavior in development:** `AllowedOrigins: ["*"]`, `AllowCredentials: false`.

**Production:** Explicit origin whitelist required. No wildcards with credentials.

**Preflight:** All `OPTIONS` requests return `204 No Content` with CORS headers. No auth required for preflight.

---

## Future Considerations

### Webhooks

Tenant-configurable webhook endpoints for event-driven integrations.

**Planned events:**

| Event | Trigger |
|-------|---------|
| `user.created` | New user created in tenant |
| `user.updated` | User profile modified |
| `user.deactivated` | User soft-deleted |
| `user.role_changed` | User role assignments modified |
| `tenant.created` | New tenant provisioned |
| `tenant.deactivated` | Tenant soft-deleted |
| `role.created` | New role created |
| `role.permissions_changed` | Role permissions modified |

**Webhook payload format:**

```json
{
  "event": "user.created",
  "tenant_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
  "timestamp": "2025-06-01T14:00:00Z",
  "data": { ... },
  "id": "wh-a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**Management endpoints (future):**

- `POST /api/v1/webhooks` — create webhook subscription
- `GET /api/v1/webhooks` — list webhooks
- `DELETE /api/v1/webhooks/:id` — remove webhook
- `POST /api/v1/webhooks/:id/secret` — rotate signing secret

Delivery: HMAC-SHA256 signed payloads, exponential retry (5 max), dead-letter queue after exhaustion.

---

### Bulk Operations

- `POST /api/v1/users/bulk` — create up to 100 users in single request
- `PUT /api/v1/users/bulk/roles` — assign roles to multiple users
- `DELETE /api/v1/users/bulk` — deactivate multiple users by ID list

**Bulk request pattern:**

```json
{
  "items": [
    { "email": "user1@example.com", "password": "pass1!", "first_name": "A", "last_name": "B" },
    { "email": "user2@example.com", "password": "pass2!", "first_name": "C", "last_name": "D" }
  ]
}
```

**Bulk response pattern:**

```json
{
  "data": {
    "succeeded": 2,
    "failed": 0,
    "results": [
      { "index": 0, "status": "created", "id": "..." },
      { "index": 1, "status": "created", "id": "..." }
    ]
  }
}
```

---

### Import / Export

- `GET /api/v1/tenants/:id/export` — export tenant data (users, roles, permissions) as JSON
- `POST /api/v1/tenants/:id/import` — import tenant data from JSON (merge or replace mode)

**Export response:** `Content-Disposition: attachment; filename="tenant-{slug}-{date}.json"`

---

### Additional Planned Endpoints

| Endpoint | Description | Priority |
|----------|-------------|----------|
| `POST /api/v1/auth/forgot-password` | Initiate password reset flow | P1 |
| `POST /api/v1/auth/reset-password` | Complete password reset with token | P1 |
| `POST /api/v1/auth/verify-email` | Verify email with token | P1 |
| `PUT /api/v1/auth/change-password` | Change password (authenticated) | P1 |
| `GET /api/v1/users/:id/roles` | Get user's role assignments | P2 |
| `GET /api/v1/users/:id/permissions` | Get user's effective permissions | P2 |
| `GET /api/v1/audit-log` | Tenant audit trail | P2 |
| `GET /api/v1/health` | Health check (public) | P0 |
| `GET /api/v1/health/ready` | Readiness probe (public) | P0 |
