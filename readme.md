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


#### Docker
One image holds the `api`, `worker` and `migrate` binaries. `compose.yaml` runs Redis, applies migrations, then starts the API and the worker, restarts them if they crash or the server reboots, and rotates their logs. Configuration comes from `.env`.
```bash
docker build -t ghcr.io/sabbir185/gpi-connect-api .   # local image
docker compose up -d                                  # start everything
docker compose logs -f api worker
docker compose down
```
The API listens on `127.0.0.1:${PORT}`; put a reverse proxy (e.g. Caddy) in front of it for HTTPS.


#### Deployment
On every push and pull request to `main` and `dev`, GitHub Actions checks formatting, vets and tests. On a push to `main` it then builds the image, pushes it to GitHub Container Registry (tagged with the commit and `latest`), copies `compose.yaml` to the server and runs `docker compose pull && docker compose up -d`. To roll back, re-run an earlier deploy from the Actions tab.

Server setup, once (Ubuntu):
```bash
curl -fsSL https://get.docker.com | sh            # Docker, starts on boot
useradd -m -s /bin/bash -G docker gpi             # deploy user
install -d -o gpi -g gpi /opt/gpi-connect-api     # put the production .env here (chmod 600)
```
Create an SSH key pair for GitHub Actions and add the public key to `/home/gpi/.ssh/authorized_keys`.

GitHub → Settings → Environments → `production`:

| Name | Kind | Value |
|---|---|---|
| `DEPLOY_HOST` | secret | Server IP or hostname |
| `DEPLOY_USER` | secret | `gpi` |
| `DEPLOY_SSH_KEY` | secret | The private key |
| `DEPLOY_KNOWN_HOSTS` | secret | Output of `ssh-keyscan <host>` |
| `DEPLOY_PATH` | variable | `/opt/gpi-connect-api` |

App secrets stay in the server's `.env`.
