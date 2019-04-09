# Executables
GO := /usr/local/go/bin/go
SHELL := /bin/bash

# Naming and directories
NAME := $(shell echo $${PWD\#\#*/})
PKG := github.com/W1lkins/$(NAME)

VERSION := $(shell cat VERSION.txt)
GITCOMMIT := $(shell git rev-parse --short HEAD)

PREFIX?=$(shell pwd)
BUILDDIR := ${PREFIX}/cross

# Flags and other vars
CGO_ENABLED := 0
BUILDTAGS :=

CTIMEVAR=-X main.VERSION=$(VERSION) -X main.GITCOMMIT=$(GITCOMMIT)
LDFLAGS=-ldflags "-w $(CTIMEVAR)"
LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

# Files
GO_SRC := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOOSARCHES = $(shell cat .goosarch)

.PHONY: build
build: $(GO_SRC) VERSION.txt ## Build dynamic executable
	@echo "+ $@"
	$(GO) build -tags "$(BUILDTAGS)" ${LDFLAGS} -o $(NAME) .

.PHONY: static
static: ## Builds static executable
	@echo "+ $@"
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
				-tags "$(BUILDTAGS) static_build" \
				${GO_LDFLAGS_STATIC} -o $(NAME) .

all: clean build fmt lint test staticcheck vet install ## Runs a clean, build, fmt, lint, test, staticcheck, vet and install

.PHONY: fmt
fmt: ## Runs gofmt
	@echo "+ $@"
	@gofmt -s -l -w $(GO_SRC)

.PHONY: lint
lint: ## Runs golint
	@echo "+ $@"
	@golint $(GO_SRC)

.PHONY: test
test: ## Runs tests
	@echo "+ $@"
	@$(GO) test -v -tags "$(BUILDTAGS) cgo" $(GO_SRC)

.PHONY: vet
vet: ## Runs go vet
	@echo "+ $@"
	@$(GO) vet $(GO_SRC)

.PHONY: staticcheck
staticcheck: ## Runs staticcheck
	@echo "+ $@"
	@staticcheck $(GO_SRC)

.PHONY: install
install: ## Installs executable/package
	@echo "+ $@"
	$(GO) install -a -tags "$(BUILDTAGS)" ${GO_LDFLAGS} .

define buildpretty
mkdir -p $(BUILDDIR)/$(1)/$(2);
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
	 -o $(BUILDDIR)/$(1)/$(2)/$(NAME) \
	 -a -tags "$(BUILDTAGS) static_build netgo" \
	 -installsuffix netgo ${LDFLAGS_STATIC} .;
md5sum $(BUILDDIR)/$(1)/$(2)/$(NAME) > $(BUILDDIR)/$(1)/$(2)/$(NAME).md5;
sha256sum $(BUILDDIR)/$(1)/$(2)/$(NAME) > $(BUILDDIR)/$(1)/$(2)/$(NAME).sha256;
endef

.PHONY: cross
cross: *.go VERSION.txt ## Builds cross-compiled binaries in the form GOOS/GOARCH/x
	@echo "+ $@"
	$(foreach GOOSARCH,$(GOOSARCHES), $(call buildpretty,$(subst /,,$(dir $(GOOSARCH))),$(notdir $(GOOSARCH))))

define buildrelease
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
	 -o $(BUILDDIR)/$(NAME)-$(1)-$(2) \
	 -a -tags "$(BUILDTAGS) static_build netgo" \
	 -installsuffix netgo ${LDFLAGS_STATIC} .;
md5sum $(BUILDDIR)/$(NAME)-$(1)-$(2) > $(BUILDDIR)/$(NAME)-$(1)-$(2).md5;
sha256sum $(BUILDDIR)/$(NAME)-$(1)-$(2) > $(BUILDDIR)/$(NAME)-$(1)-$(2).sha256;
endef

.PHONY: release
release: *.go VERSION.txt ## Builds cross-compiled binaries in the form x-GOOS-GOARCH
	@echo "+ $@"
	$(foreach GOOSARCH,$(GOOSARCHES), $(call buildrelease,$(subst /,,$(dir $(GOOSARCH))),$(notdir $(GOOSARCH))))

.PHONY: tag
tag: ## Create new git tag
	git tag -sa $(VERSION) -m "$(VERSION)"
	@echo "Run git push origin $(VERSION) to push your new tag to the remote."

.PHONY: vendor
vendor: ## Vendors dependencies
	@$(RM) go.sum
	@$(RM) -r vendor
	GO111MODULE=on $(GO) mod init || true
	GO111MODULE=on $(GO) mod tidy
	GO111MODULE=on $(GO) mod download
	GO111MODULE=on $(GO) mod vendor

.PHONY: clean
clean: ## Cleanup any binaries
	@echo "+ $@"
	$(RM) $(NAME)
	$(RM) -r $(BUILDDIR)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/^[^:]*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

