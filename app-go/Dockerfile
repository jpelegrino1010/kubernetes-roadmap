# =============================================================================
# Stage 1 — Build
# Use the official Go image; pin to a specific minor for reproducibility.
# =============================================================================
FROM golang:1.22-alpine AS builder

WORKDIR /src

# Copy dependency manifests first so Docker layer-caches the download.
COPY go.mod ./

RUN go mod download

# Copy source after modules are cached.
COPY . .

# Build a fully static binary:
#   CGO_ENABLED=0  — no cgo, no libc dependency
#   -trimpath      — reproducible, strips local paths from stack traces
#   -ldflags       — strip debug info and DWARF to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/api \
    ./cmd/api

# =============================================================================
# Stage 2 — Runtime (scratch)
# scratch has no shell, no libc, no attack surface.
# =============================================================================
FROM scratch

# Copy the CA bundle so the binary can make outbound TLS calls if needed.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy only the compiled binary.
COPY --from=builder /bin/api /api

# Document the port; does not actually publish it.
EXPOSE 8080

# Run as a non-root UID (scratch has no user database; set numeric UID).
USER 65532:65532

ENTRYPOINT ["/api"]
