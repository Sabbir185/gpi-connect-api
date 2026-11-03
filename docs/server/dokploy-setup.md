# Dokploy + GitHub Actions CI/CD Deployment Guide

This guide explains how to deploy `gpi-connect-api` to a **Dokploy** instance using GitHub Actions and Dokploy's **Deploy Webhook**.

---

## 1. Overview of the Flow

```mermaid
flowchart LR
    A["git push main"] --> B["GitHub Actions CI"]
    B --> C["Run Tests & Linter"]
    C --> D["Build & Push to GHCR"]
    D --> E["Trigger Dokploy Webhook (curl)"]
    E --> F["Dokploy Server"]
    F --> G["Pull new image & Restart Containers"]
```

---

## 2. Step-by-Step Dokploy Setup

### Step 1: Create a Service in Dokploy
You can run this project in Dokploy using either **Compose** or **Application**:

#### Option A: Dokploy Compose (Recommended - runs API, Worker & Redis together)
1. Go to your **Dokploy Dashboard** ➔ Select/Create a Project ➔ **Create Service** ➔ Choose **Compose**.
2. Service Type: Choose **Raw Docker Compose** (or Git).
3. Paste the contents of your [compose.yaml](file:///home/sabbir185/Documents/sabbir/projects/GPI%20Connect/gpi-connect-api/compose.yaml) into the Compose editor.
4. Set the image tag to `latest`:
   ```yaml
   image: ghcr.io/sabbir185/gpi-connect-api:latest
   ```

#### Option B: Standalone Application
1. **Create Service** ➔ Choose **Application**.
2. Provider / Build Type: Select **Docker** (or Docker Image).
3. Image URL: `ghcr.io/sabbir185/gpi-connect-api:latest`
4. Container Port: `8080` (or your configured `PORT`).

---

### Step 2: Configure Environment Variables in Dokploy
1. In Dokploy, click on your service ➔ open the **Environment** tab.
2. Paste all required production environment variables:
   ```env
   PORT=8080
   DATABASE_URL=postgres://user:password@host:5432/dbname?sslmode=disable
   REDIS_URL=redis:6379
   JWT_SECRET=your-production-jwt-secret
   # Brevo, Firebase, and other API credentials...
   ```
3. Click **Save**.

---

### Step 3: Copy the Deploy Webhook URL from Dokploy
1. In your Dokploy service page, navigate to the **Deployments** or **General** tab.
2. Look for the **"Deploy Webhook"** field.
3. Click the copy icon to copy the full URL (it will look like `https://dokploy.yourdomain.com/api/deploy/xxxx-xxxx-xxxx` or similar).

---

### Step 4: Add the Webhook to GitHub Secrets
1. Go to your GitHub Repository:
   **Settings** ➔ **Secrets and variables** ➔ **Actions**.
2. Under **Repository secrets**, click **New repository secret**.
3. Name: **`DOKPLOY_WEBHOOK_URL`**
4. Secret: Paste the webhook URL copied from Dokploy.
5. Click **Add secret**.

---

## 3. GHCR Access Permission (Important!)

Dokploy needs permission to pull images from `ghcr.io/sabbir185/gpi-connect-api`.

### Make the Image Package Public (Easiest):
1. On GitHub, go to your profile ➔ **Packages** ➔ click `gpi-connect-api`.
2. Click **Package settings** (bottom of right sidebar).
3. Scroll to **Danger Zone** ➔ **Change package visibility** ➔ select **Public**.

*OR*

### If Keeping Image Private:
1. Generate a GitHub Personal Access Token (PAT Classic) with `read:packages` scope.
2. In Dokploy: **Settings** ➔ **Docker Registries** ➔ **Add Registry**:
   - Provider: `GitHub Packages`
   - Server: `ghcr.io`
   - Username: `sabbir185`
   - Password/Token: Your GitHub PAT.

---

## 4. Triggering Deployments

Whenever you push to the `main` branch:
1. GitHub Actions runs code checks (`gofmt`, `go vet`, `go test`).
2. Buildx compiles and pushes `ghcr.io/sabbir185/gpi-connect-api:latest`.
3. GitHub Actions calls your Dokploy webhook.
4. Dokploy automatically pulls the new image and restarts the container gracefully with zero manual SSH commands required.

