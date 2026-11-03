# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    for name in api worker migrate; do \
      CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/$name ./cmd/$name || exit 1; \
    done

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H -u 10001 app

WORKDIR /app
COPY --from=build /out/ ./bin/
COPY migrations ./migrations

USER app

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO /dev/null "http://127.0.0.1:${PORT}/healthz" || exit 1

CMD ["/app/bin/api"]
