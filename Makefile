# This file is the development Makefile for the bully-algorithm project.
# All variables listed here are used as substitution in these Makefile targets.

SERVICE-NAME = bully-algorithm

define ENV-CONFIGURATION
ENV='dev'
endef

################################################################################


# Install all dependencies required.
#
# NOTE: Docker & Docker Compose should already be installed.
.PHONY: install
install:
		curl https://glide.sh/get | sh
		go get -u github.com/alecthomas/gometalinter
		gometalinter --install
		glide update
		glide install

# Build project binaries.
.PHONY: build
build: lint
		cd cmd/bully && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bully
		cd cmd/data-viz && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o data-viz

# Build project Docker images for dev environment.
.PHONY: docker-build
docker-build: build
		cp cmd/bully/conf/bully.conf.yaml build/docker/bully/ && \
		cp cmd/bully/bully build/docker/bully/ && \
		docker build -t timtosi/bully:latest build/docker/bully/ && \
		cp cmd/data-viz/data-viz build/docker/data-viz/ && \
		cp -r cmd/data-viz/assets build/docker/data-viz/ && \
		docker build -t timtosi/data-viz:latest build/docker/data-viz/

# Runs linter against the service codebase.
#
# NOTE: This rule require gcc to be found in the `$PATH`.
.PHONY: lint
lint:
		@gometalinter --config=conf/gometalinter_conf.json ./... && \
		echo "linter pass ok !"

# Runs test suite.
.PHONY: test
test: lint
		go test -v -coverprofile=coverage.txt -tags integration -race -cover -timeout=120s $$(glide novendor)

# Run project locally.
.PHONY: run
run:
		docker-compose -f deployments/docker-compose.yaml up
