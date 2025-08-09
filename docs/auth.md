## Backend Authentication System

This document explains how authentication works in the backend and how to integrate it in the frontend.

### Overview

- **Session-backed JWT**: Short-lived access JWT (15m) and a refresh token (7d) are issued per session.
- **Cookies**: Tokens are set as secure, HttpOnly cookies by the backend; the frontend does not store tokens.
- **CSRF protection**: Double-submit cookie; non-GET requests must send `X-CSRF-Token` that matches the `csrf_token` cookie.
- **Providers**: Email/password and Google OAuth.
- **Email flows**: Verification (on signup), resend verification, and password reset.

### Key Components

- **JWT Service** (`backend/core/auth/jwt.go`):

  - Issues short-lived access tokens embedding `user_id`, `email`, and a `session_token` reference.
  - Validates tokens on requests.

- **Session Service** (`backend/core/auth/session.go`):

  - Persists sessions in `user_sessions` with rotation on refresh.
  - Issues `access_token` (JWT) and `refresh_token` per session.
  - Supports: create, validate, refresh, revoke one, revoke all.

- **CSRF Middleware** (`backend/api/middlewares/auth_middleware.go`):

  - `GET /csrf-token` sets a `csrf_token` cookie and returns the same token in JSON.
  - For non-GET/HEAD/OPTIONS, requires header `X-CSRF-Token` equal to the `csrf_token` cookie.

- **CORS and Security Headers** (`backend/api/middlewares/auth_middleware.go`):
  - Allowed origins: `http://localhost:3000`, `https://localhost:3000`, `https://syncer.net`.
  - Sets standard security headers and a restrictive Content Security Policy.

### Cookies

- `access_token`: JWT access token
  - Flags: HttpOnly, Secure, SameSite=Strict, Path=/
  - Lifetime: 15 minutes
- `refresh_token`: Refresh token
  - Flags: HttpOnly, Secure, SameSite=Strict, Path=/
  - Lifetime: 7 days
- `csrf_token`: CSRF token
  - Flags: Secure, SameSite=Strict, Path=/, not HttpOnly (so only the server reads cookie; frontend receives the token in JSON response)
  - Lifetime: 1 hour

Important: With `SameSite=Strict` and `Secure=true`, cookies are sent only over HTTPS and are not sent in cross-site requests. For browser clients, deploy frontend and backend on the same site and over HTTPS.

### Data Model (relevant tables)

- `users`: primary profile (email, name, avatar, timestamps)
- `user_auth_methods`: provider records (email, google), verification/reset tokens, password hash
- `user_sessions`: session and refresh tokens with metadata (IP, UA, expiries)

### Endpoints

- `GET /health` → { status }
- `GET /csrf-token` → Sets `csrf_token` cookie and returns `{ csrf_token }`

Auth (some routes require CSRF header):

- `POST /auth/register` body `{ email, password, full_name }` → sets cookies, returns `{ user }`
- `POST /auth/login` body `{ email, password }` → sets cookies, returns `{ user }`
- `POST /auth/verify-email` body `{ token }` → `{ message }`
- `POST /auth/resend-verification` body `{ email }` → `{ message }`
- `POST /auth/forgot-password` body `{ email }` → `{ message }`
- `POST /auth/reset-password` body `{ token, password }` → `{ message }`
- `GET /auth/validate-reset-token?token=...` → `{ message, email }`
- `GET /auth/google` → `{ auth_url }` (frontend should redirect user to this URL)
- `GET /auth/google/callback` → sets cookies, returns `{ user }` (see integration note below)
- `POST /auth/refresh` → rotates session, sets new cookies
- `POST /auth/logout` → revokes current session
- `POST /auth/logout-all` → revokes all sessions for the user

Protected API:

- `GET /auth/profile` → requires valid `access_token` cookie (or `Authorization: Bearer <jwt>`), returns `{ user_id, email }`

## Frontend Integration Guide

The repository includes a ready-to-use client: `frontend/src/api.ts` (`ApiService`). It handles CSRF fetching, credentials, 401 refresh flow, and retries.

### 1) Configure environment

- **VITE_API_URL**: Base URL of the backend (e.g., `https://yourdomain.com` or `http://localhost:8080`).
- Ensure the backend CORS `allowedOrigins` contains your frontend origin if they differ in dev.
- Use HTTPS in the browser so `Secure` cookies are set and sent.

### 2) Bootstrap CSRF and credentials

- All state-changing requests (POST/PUT/DELETE) require `X-CSRF-Token` and cookies.
- `ApiService` automatically fetches `GET /csrf-token` before non-GET requests and includes `credentials: 'include'`.

Optional prefetch on app startup to avoid the very first POST incurring a CSRF roundtrip:

```ts
// Example: call once on app mount
import { apiService } from '~/api';

await apiService.healthCheck(); // any GET with credentials primes cookies
// or explicitly fetch: await fetch(`${API}/csrf-token`, { credentials: 'include' })
```

### 3) Register and Login

```ts
import { apiService } from '~/api';

// Register
const { user } = await apiService.register({
    email: 'user@example.com',
    password: 'hunter2hunter2',
    full_name: 'Jane Doe',
});

// Login
const { user } = await apiService.login({
    email: 'user@example.com',
    password: 'hunter2hunter2',
});

// The backend sets HttpOnly cookies; no token storage on the client.
```

### 4) Using protected endpoints

```ts
const profile = await apiService.getProfile(); // cookies are sent automatically
```

`ApiService` will:

- Retry once on 401 by POSTing `/auth/refresh` with CSRF header and then re-issuing the original request.
- Retry once on 403 (CSRF mismatch) by re-fetching `/csrf-token` and re-issuing the request.

### 5) Email verification and password reset

```ts
// Verify email from a link containing a token
await apiService.verifyEmail({ token });

// Resend verification
await apiService.resendVerification({ email });

// Forgot password
await apiService.forgotPassword({ email });

// Reset password
await apiService.resetPassword({ token, password: 'new-strong-pass' });
```

### 6) Google OAuth login

```ts
// 1) Ask backend for the provider URL
const { auth_url } = await apiService.getGoogleAuthUrl();

// 2) Redirect the user
window.location.href = auth_url;

// 3) After Google redirects back to /auth/google/callback, the backend
//    sets cookies and returns JSON. Navigate the user back to your app
//    (or open the flow in a popup and then refresh state):
//    e.g., after returning to your app, call:
const profile = await apiService.getProfile(); // should now succeed
```

Notes:

- For a slick UX, consider opening the OAuth flow in a popup and, after callback completion, closing the popup and reloading session state in the main app (custom HTML in callback would be needed if you want postMessage; backend currently returns JSON).

### 7) Logout

```ts
await apiService.logout(); // revoke current session
await apiService.logoutAll(); // revoke all sessions
```

### 8) Fallback: Bearer tokens

The middleware accepts `Authorization: Bearer <jwt>` if no cookie is present. In the browser, you cannot read the HttpOnly cookie to create this header. This path is mainly for non-browser clients or custom tooling.

### Common pitfalls

- **HTTPS required for cookies**: `Secure=true` means cookies are not set over plain HTTP.
- **SameSite=Strict**: Cookies will not be sent in cross-site requests. Prefer same-origin deployment (frontend served from the same site as the API) in production.
- **CSRF**: Always send `X-CSRF-Token` for non-GET requests; the repository client already handles this.

### Required environment variables (backend)

- `DATABASE_URL`
- `JWT_SECRET`
- `RESEND_API_KEY`, `RESEND_FROM_EMAIL`
- `FRONTEND_URL` (used in email links)
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL`
