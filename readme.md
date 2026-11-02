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


#### Deployment (CI/CD)
`.github/workflows/ci-cd.yml` checks formatting, vets, builds and tests every push and pull request to `main` and `dev`. A push to `main` that passes is deployed: the workflow connects over SSH, checks out that exact commit on the server and runs `deploy/deploy.sh`, which downloads dependencies, builds `api`, `worker` and `migrate`, runs migrations and restarts both services. If `/healthz` doesn't answer within 15s, the previous binaries come back. Migrations are not rolled back, so keep each one compatible with the previous release.

One-time server setup (Ubuntu, as root; adjust the user and path to taste):
```bash
# Go (same version as go.mod), git and curl
apt-get install -y git curl
curl -fsSL https://go.dev/dl/go1.27.0.linux-amd64.tar.gz | tar -C /usr/local -xz

# App user with a read-only GitHub deploy key (repo → Settings → Deploy keys)
useradd --system --create-home --shell /bin/bash gpi
sudo -u gpi ssh-keygen -t ed25519 -N "" -f /home/gpi/.ssh/id_ed25519   # add the .pub as a deploy key
sudo -u gpi git clone git@github.com:Sabbir185/gpi-connect-api.git /opt/gpi-connect-api  # create /opt/gpi-connect-api owned by gpi first

# Production config, readable only by the app user
sudo -u gpi cp /opt/gpi-connect-api/.env.example /opt/gpi-connect-api/.env   # fill in, ENV="production"
chmod 600 /opt/gpi-connect-api/.env

# Services, and permission for the deploy to restart them
cp /opt/gpi-connect-api/deploy/gpi-*.service /etc/systemd/system/
systemctl daemon-reload && systemctl enable gpi-api gpi-worker
echo 'gpi ALL=(root) NOPASSWD: /usr/bin/systemctl restart gpi-api, /usr/bin/systemctl restart gpi-worker' > /etc/sudoers.d/gpi-deploy

# Key GitHub Actions logs in with: put the private key in DEPLOY_SSH_KEY
ssh-keygen -t ed25519 -N "" -C github-actions -f gh_deploy && cat gh_deploy.pub >> /home/gpi/.ssh/authorized_keys
```
Then run the first deploy by hand: `sudo -u gpi bash /opt/gpi-connect-api/deploy/deploy.sh`.

GitHub → repo Settings → Environments → `production`:

| Name | Kind | Value |
|---|---|---|
| `DEPLOY_HOST` | secret | Server IP or hostname |
| `DEPLOY_USER` | secret | `gpi` |
| `DEPLOY_PORT` | secret | SSH port, optional (default 22) |
| `DEPLOY_SSH_KEY` | secret | Private key whose public half is in `authorized_keys` |
| `DEPLOY_KNOWN_HOSTS` | secret | Output of `ssh-keyscan -p <port> <host>`, so the server's identity is verified |
| `DEPLOY_PATH` | variable | `/opt/gpi-connect-api` |

App secrets (database, JWT, Firebase, SMS, Brevo) stay in the server's `.env` and never pass through GitHub.
