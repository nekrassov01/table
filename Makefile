GOBIN ?= $(shell go env GOPATH)/bin

VERSION_DIR := ./internal/version/
VERSION     := $$(make version)

HAS_LINT := $(shell command -v $(GOBIN)/golangci-lint 2> /dev/null)
HAS_VULN := $(shell command -v $(GOBIN)/govulncheck 2> /dev/null)
HAS_BUMP := $(shell command -v $(GOBIN)/gobump 2> /dev/null)

BIN_LINT := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
BIN_VULN := golang.org/x/vuln/cmd/govulncheck@latest
BIN_BUMP := github.com/x-motemen/gobump/cmd/gobump@latest

TEST_CMD     = go test -race -cover -v -coverprofile coverage.out -covermode atomic
TEST_TARGETS = all table internal text html markdown backlog csv contract golden
TEST_TARGET  = $(if $(target),$(target),all)
TEST_PATTERN = $(if $(run),$(value run),$(TEST_PATTERN_$(TEST_TARGET)))

TEST_PACKAGE_all      = ./...
TEST_PACKAGE_table    = .
TEST_PACKAGE_internal = ./internal/...
TEST_PACKAGE_text     = ./text
TEST_PACKAGE_html     = ./html
TEST_PACKAGE_markdown = ./markdown
TEST_PACKAGE_backlog  = ./backlog
TEST_PACKAGE_csv      = ./csv
TEST_PACKAGE_contract = ./...
TEST_PACKAGE_golden   = ./...

TEST_PATTERN_contract = ^TestContract_
TEST_PATTERN_golden   = ^TestGolden_

.PHONY: deps deps-lint deps-vuln deps-bump clean example check test cover bench lint vuln version check-git check-branch bump

# -------
#  deps
# -------

deps: deps-lint deps-vuln deps-bump

deps-lint:
ifndef HAS_LINT
	go install $(BIN_LINT)
endif

deps-vuln:
ifndef HAS_VULN
	go install $(BIN_VULN)
endif

deps-bump:
ifndef HAS_BUMP
	go install $(BIN_BUMP)
endif

# --------
#  utils
# --------

clean:
	go clean
	rm -f $(NAME) coverage.out coverage.html cpu.prof mem.prof $(NAME).test
	@cd benchmarks && $(MAKE) clean

example:
	@cd examples && $(MAKE) example

# --------
#  check
# --------

check: test cover bench lint vuln

test:
ifeq ($(filter $(TEST_TARGET),$(TEST_TARGETS)),)
	$(error unknown test target: $(TEST_TARGET))
endif
	$(TEST_CMD)$(if $(TEST_PATTERN), -run '$(TEST_PATTERN)') $(TEST_PACKAGE_$(TEST_TARGET))

cover:
	go tool cover -html coverage.out -o coverage.html

bench:
	@cd benchmarks && $(MAKE) bench

lint: deps-lint
	golangci-lint run --verbose ./...

vuln: deps-vuln
	govulncheck -test -show verbose ./...

# ----------
#  version
# ----------

version: deps-bump
	@echo $(shell gobump show -r $(VERSION_DIR))

check-git:
ifneq ($(shell git status --porcelain),)
	$(error git workspace is dirty)
endif

check-branch:
ifndef branch
	$(error branch is undefined)
endif

bump: check-branch deps-bump
	gobump up -w $(VERSION_DIR)
	git commit -am "chore: bump up version to $(VERSION)"
	git push origin $(branch)
