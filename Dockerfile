# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.26 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Embed version/commit info (works for `docker build`, not `go install`)
ARG VERSION=dev
ARG COMMIT=unknown
RUN go build \
    -ldflags="-X 'github.com/aallbrig/allbctl/cmd.Version=${VERSION}' \
              -X 'github.com/aallbrig/allbctl/cmd.Commit=${COMMIT}'" \
    -o /allbctl .

# ── Runtime stage ──────────────────────────────────────────────────────────────
# Use debian:bookworm-slim (not distroless) because allbctl calls external
# binaries: git, systemctl, lsblk, ip, ss, docker, kubectl, aws, gcloud, etc.
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    curl \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /allbctl /usr/local/bin/allbctl

# NOTE: For `allbctl status` to see real host info, run with host namespaces:
#   docker run --pid=host --net=host \
#     -v /proc:/proc:ro -v /sys:/sys:ro -v /etc:/etc:ro \
#     -v $HOME:$HOME:ro \
#     allbctl status
#
# See `make docker-run` for a ready-made wrapper.

ENTRYPOINT ["/usr/local/bin/allbctl"]
CMD ["help"]