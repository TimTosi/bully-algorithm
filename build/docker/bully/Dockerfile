# Dockerfile for the bully binary.

# ------------------------------------------------------------------------------

# Builder Image using golang:1.11.5-alpine3.9.
FROM golang@sha256:a6435c88400d0d25955ccdea235d2b2bd72650bbab45e102e4ae31a0359dbfb2 AS builder

# Create non-root user.
RUN adduser -D -g '' appuser

# ------------------------------------------------------------------------------

# Final Image.
FROM scratch
LABEL maintainer="timothee.tosi@gmail.com"

# Copy configuration file.
COPY ./bully.conf.yaml /bully.conf.yaml

# Copy service binary.
COPY ./bully /bully

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd

# Change to non-root user.
USER appuser

# Change working directory.
WORKDIR /

# Define entrypoint
ENTRYPOINT [ "/bully" ]

# Entry point arguments.
# NOTE: Override CMD instruction for spawning multiples bullies.
CMD [ "0" ]
