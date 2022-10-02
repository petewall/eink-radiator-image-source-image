HAS_GINKGO := $(shell command -v ginkgo;)
HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)
HAS_COUNTERFEITER := $(shell command -v counterfeiter;)
PLATFORM := $(shell uname -s)

# #### DEPS ####
.PHONY: deps-counterfeiter deps-ginkgo deps-modules

deps-counterfeiter: deps-go-binary
ifndef HAS_COUNTERFEITER
	go install github.com/maxbrunsfeld/counterfeiter/v6@latest
endif

deps-ginkgo:
ifndef HAS_GINKGO
	go install github.com/onsi/ginkgo/v2/ginkgo
endif

deps-modules:
	go mod download

# #### TEST ####
.PHONY: lint test

lint:
ifndef HAS_GOLANGCI_LINT
ifeq ($(PLATFORM), Darwin)
	brew install golangci-lint
endif
ifeq ($(PLATFORM), Linux)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
endif
endif
	golangci-lint run

test: deps-modules deps-ginkgo
	ginkgo -r .

# #### BUILD ####
.PHONY: build
SOURCES = $(shell find . -name "*.go" | grep -v "_test\." )
VERSION := $(or $(VERSION), dev)
LDFLAGS="-X github.com/petewall/eink-radiator-image-source-image/v2/cmd.Version=$(VERSION)"

build/image: $(SOURCES) deps-modules
	go build -o build/image -ldflags ${LDFLAGS} github.com/petewall/eink-radiator-image-source-image/v2

build: build/image
