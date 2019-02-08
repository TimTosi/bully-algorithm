# Dockerfile for the bully binary.

# ------------------------------------------------------------------------------

# Builder Image.
FROM golang:1.11.5-alpine3.9 AS builder

# Install dependencies.
RUN apk update && apk add --no-cache git

# Retrieve source files.
RUN go get -d -v github.com/timtosi/bully-algorithm

# Create non-root user.
RUN adduser -D -g '' appuser

# Define working directory.
WORKDIR $GOPATH/src/github.com/timtosi/bully-algorithm/

# Build static binary.
RUN cd $GOPATH/src/github.com/timtosi/bully-algorithm/cmd/bully/ && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/bully

# ------------------------------------------------------------------------------

# Final Image.
FROM scratch
LABEL maintainer="timothee.tosi@gmail.com"

# Copy service binary.
COPY --from=builder /go/bin/bully /bully

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd

# Expose HTTP ports.
EXPOSE 9990-9999

# Change to non-root user.
USER appuser

# Define entrypoint
ENTRYPOINT [ "/bully" ]

# Entry point arguments.
# NOTE: Override CMD instruction for spawning multiples bullies.
CMD [ "0" ]
