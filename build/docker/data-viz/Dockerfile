# Dockerfile for the data-viz binary.

# ------------------------------------------------------------------------------

# Builder Image using golang:1.11.5-alpine3.9.
FROM golang@sha256:a6435c88400d0d25955ccdea235d2b2bd72650bbab45e102e4ae31a0359dbfb2 AS builder

# Install dependencies.
RUN apk update && apk add --no-cache ca-certificates && \
    update-ca-certificates

# Create non-root user.
RUN adduser -D -g '' appuser

# ------------------------------------------------------------------------------

# Final Image.
FROM scratch
LABEL maintainer="timothee.tosi@gmail.com"

# Copy service binary.
COPY ./data-viz /data-viz

# Copy web assets.
COPY ./assets /assets

# Import certificates from the builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd

# Expose HTTP ports.
EXPOSE 8080-8081

# Change to non-root user.
USER appuser

# Change working directory.
WORKDIR /

# Define entrypoint
ENTRYPOINT [ "/data-viz" ]
