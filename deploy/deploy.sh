#!/usr/bin/env bash
# Runs on the server from the app directory, after CI checked out the commit to deploy.
# Builds next to the running binaries, migrates, swaps them in and restarts. If the API
# fails its health check, the previous binaries are restored. Migrations are not rolled
# back, so keep them backward compatible with the previous release.
set -euo pipefail

SERVICES=(gpi-api gpi-worker)
BINARIES=(api worker migrate)
HEALTH_RETRIES=15

export PATH="$PATH:/usr/local/go/bin"
cd "$(dirname "$0")/.."

port=$(grep -E '^PORT=' .env | head -1 | cut -d= -f2 | tr -d '"' || true)
health_url="http://127.0.0.1:${port:-8080}/healthz"

echo "==> Deploying $(git rev-parse --short HEAD)"

echo "==> Installing dependencies"
go mod download

echo "==> Building"
mkdir -p bin
for name in "${BINARIES[@]}"; do
  go build -trimpath -ldflags="-s -w" -o "bin/$name.new" "./cmd/$name"
done

echo "==> Running migrations"
./bin/migrate.new up

swap_in() {
  for name in "${BINARIES[@]}"; do
    [[ -f "bin/$name" ]] && mv "bin/$name" "bin/$name.prev"
    mv "bin/$name.new" "bin/$name"
  done
}

restore_previous() {
  for name in "${BINARIES[@]}"; do
    [[ -f "bin/$name.prev" ]] && mv "bin/$name.prev" "bin/$name"
  done
}

restart() {
  for service in "${SERVICES[@]}"; do
    sudo systemctl restart "$service"
  done
}

healthy() {
  for _ in $(seq "$HEALTH_RETRIES"); do
    curl -fsS --max-time 2 "$health_url" > /dev/null 2>&1 && return 0
    sleep 1
  done
  return 1
}

echo "==> Restarting services"
swap_in
restart

if ! healthy; then
  echo "!! Health check failed at $health_url, rolling back" >&2
  restore_previous
  restart
  exit 1
fi

for service in "${SERVICES[@]}"; do
  systemctl is-active --quiet "$service" || { echo "!! $service is not running" >&2; exit 1; }
done

rm -f bin/*.prev
echo "==> Deployed $(git rev-parse --short HEAD)"
