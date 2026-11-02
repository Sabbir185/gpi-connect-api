#### GPI Connect — Connect Your People. Simplify Your Operations.

Your daily activity solution
- Smart Geo location based attendance
- Selfie based attendance
- QR code based attendance
- Daily notes and share with group or team members
- Store files and documents in cloud storage
- Offline mode support and sync when online
- Meetings and events management
- AI based insights and analytics
- Group or one-to-one chat and messaging


#### Migration
```
make migrate-create name="users table"
make migrate-up
make migrate-down
```

#### Run
```
make run       # API
make worker    # background jobs: welcome email (Brevo), welcome push (FCM)
make keygen    # random value for JWT_SECRET / OTP_SECRET
```
The API and the worker share Redis (Asynq), so the worker can run as its own process and scale separately.


#### Auth
Mobile apps sign in with a phone OTP (SMS) or Google through Firebase. Either way the API returns its own session:

| Endpoint | Purpose |
|---|---|
| `POST /v1/auth/otp` | Send a 6-digit code. 60s resend cooldown, 5 codes per number per hour |
| `POST /v1/auth/otp/verify` | Verify the code; signs in or creates the account |
| `POST /v1/auth/firebase` | Verify a Firebase ID token (Google); signs in, links by verified email, or creates the account |
| `POST /v1/auth/refresh` | Rotate the refresh token and issue a new access token |
| `POST /v1/auth/sign-out` | Revoke the current session and detach the device's push token |
| `GET /v1/me`, `PATCH /v1/me` | Read the profile; save name, country, avatar URL and designation (completes onboarding) |
| `GET /v1/designations?q=&limit=` | Search the designation catalogue. A picked title is approved, a typed-in one waits for review |
| `GET /v1/me/devices`, `PUT /v1/me/devices/{installationId}`, `DELETE /v1/me/devices/{deviceId}` | Signed-in devices, push token registration, remote sign-out |

- **Access token:** HS256 JWT (15 min) carrying the user and session id. Every request also checks that the session is live and the user active, so sign-out and blocking apply at once.
- **Refresh token:** 32 random bytes, stored only as a SHA-256 hash, one session per app install. Each refresh rotates it. Reusing a rotated token after a 60s grace window revokes the whole session.
- **OTP:** stored as an HMAC digest bound to the phone number, single use, 5 attempts per code.
- New accounts get a welcome push and, when an email is known, a welcome email, both queued so sign-in never waits on them.
