SHELL := /usr/bin/env bash
.SHELLFLAGS := -euo pipefail -c
.DEFAULT_GOAL := help

BINARY   := terraform-provider-truenas
GOOS     := $(shell go env GOOS)
GOARCH   := $(shell go env GOARCH)
VERSION  ?= $(or $(patsubst v%,%,$(shell git describe --tags --abbrev=0 2>/dev/null)),0.1.0)
LDFLAGS  := -s -w -X main.version=$(VERSION)
TOOL     := go tool -modfile=tools/go.mod

PLUGIN_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/twilightcoders/truenas/$(VERSION)/$(GOOS)_$(GOARCH)

TIMEOUT ?= 30m
RUN     ?= .
DOC     ?= index

##@ Build

.PHONY: build
build: ## Build the provider binary into bin/
	go build -trimpath -ldflags '$(LDFLAGS)' -o bin/$(BINARY) .

.PHONY: install
install: build ## Install into the local Terraform plugin mirror
	mkdir -p '$(PLUGIN_DIR)'
	cp bin/$(BINARY) '$(PLUGIN_DIR)/$(BINARY)_v$(VERSION)'
	@echo "Installed $(VERSION) to $(PLUGIN_DIR)"

.PHONY: dev-override
dev-override: build ## Print a ~/.terraformrc dev_overrides block for this checkout
	@printf 'provider_installation {\n  dev_overrides {\n    "%s" = "%s"\n  }\n  direct {}\n}\n' \
		'registry.terraform.io/twilightcoders/truenas' '$(CURDIR)/bin'
	@echo "# Add the block above to ~/.terraformrc yourself; this target never writes it." >&2

.PHONY: clean
clean: ## Remove build and coverage output
	rm -rf bin dist coverage.out coverage.html

##@ Quality

.PHONY: test
test: ## Run unit tests
	# The provider package drives a real plugin server per test; -race on a shared CI runner
	# takes it past Go's default 10 minute limit.
	go test -race -cover -timeout 30m ./...

.PHONY: testacc
testacc: ## Run acceptance tests against TRUENAS_HOST (RUN=regex TIMEOUT=30m)
	TF_ACC=1 go test ./... -run '$(RUN)' -v -timeout '$(TIMEOUT)'

.PHONY: coverage
coverage: ## Write coverage.out and coverage.html, counting coverage wherever it happens
	# -coverpkg matters here: the generated resources are exercised by the provider tests, so a
	# per-package figure reports them as untested and invites the wrong fix.
	go test -race -count=1 -coverpkg=./... -coverprofile=coverage.out -timeout 30m ./...
	go tool cover -func=coverage.out | tail -1
	go tool cover -html=coverage.out -o coverage.html

.PHONY: schema-golden
schema-golden: ## Accept resource schema changes into internal/provider/testdata/schema.golden.json
	go test ./internal/provider -run TestSchemaGolden -update

.PHONY: lint
lint: ## Run golangci-lint
	$(TOOL) golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	$(TOOL) golangci-lint fmt ./...

##@ Documentation

.PHONY: docs
docs: ## Generate registry documentation into docs/
	$(TOOL) tfplugindocs generate --provider-name truenas

.PHONY: api-docs
api-docs: ## Print TrueNAS API docs from a box (DOC=api_methods RAW=1)
	@test -n "$${TRUENAS_HOST:-}" || { echo "TRUENAS_HOST is required" >&2; exit 1; }
	@url="https://$${TRUENAS_HOST}/api/docs/current/_sources/$(DOC).rst.txt"; \
	if [[ -n "$${RAW:-}" ]] || ! command -v lynx >/dev/null; then curl -fsSk "$$url"; \
	else curl -fsSk "$$url" | lynx -stdin -dump; fi

.PHONY: api-method
api-method: ## Print the JSON schema for one API method (METHOD=pool.snapshottask.create)
	@test -n "$${METHOD:-}" || { echo "METHOD is required" >&2; exit 1; }
	@scripts/api-method.sh '$(METHOD)'

##@ Release

.PHONY: release
release: ## Update CHANGELOG, commit and tag (VERSION=x.y.z); does not push
	@[[ "$(VERSION)" =~ ^[0-9]+\.[0-9]+\.[0-9]+$$ ]] || { echo "VERSION must be x.y.z" >&2; exit 1; }
	@command -v git-cliff >/dev/null || { echo "git-cliff is required: brew install git-cliff" >&2; exit 1; }
	@! git rev-parse 'v$(VERSION)' >/dev/null 2>&1 || { echo "tag v$(VERSION) already exists" >&2; exit 1; }
	@git diff --quiet && git diff --cached --quiet || { echo "working tree is not clean" >&2; exit 1; }
	git-cliff --tag 'v$(VERSION)' -o CHANGELOG.md
	git add CHANGELOG.md
	git commit -m 'chore(release): v$(VERSION)'
	git tag -a 'v$(VERSION)' -m 'v$(VERSION)'
	@echo "Tagged v$(VERSION). Push with: git push origin main --tags"

##@ Help

.PHONY: help
help: ## List targets
	@awk 'BEGIN {FS = ":.*##"} /^##@/ {printf "\n%s\n", substr($$0, 5)} /^[a-z-]+:.*##/ {printf "  %-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
