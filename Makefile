
BUILD_PATH="./build"

prepare:
	mkdir -p $(BUILD_PATH)

echo-env: ## Echo environment variables
	@echo go version $(shell go version)
	@echo GOPATH: $(GOPATH)
	@echo GOROOT: $(GOROOT)
	@echo GOOS: $(GOOS)
	@echo GO111MODULE: $(GO111MODULE)
	@echo PKG_CONFIG_PATH: $(PKG_CONFIG_PATH)
	@echo AWS_PROFILE: $(AWS_PROFILE)
	@echo BUILD_PATH: $(BUILD_PATH)

.PHONY: clean
clean: echo-env ## Clean
	@echo "==> Clean"
	rm -rf $(BUILD_PATH)
	rm -fv cover.out
	rm -fv cprofile.out
	rm -fv cover.out.original
	env GO111MODULE=off go clean -cache -testcache

.PHONY: nuke
nuke: clean ## Clean with -modcache
	@echo "==> Nuke-em!"
	env GO111MODULE=off go clean -modcache -x

.PHONY: build
build: echo-env prepare ## Build binary
	@echo "==> Build"
	GOOS=linux
	go build -a -o $(BUILD_PATH)/app ./

PACKAGES=$(shell go list ./...)
PACKAGES_COMMA=$(shell echo $(PACKAGES) | tr ' ' ',')
PACKAGES_WITH_TESTS = $(shell go list -f '{{if len .XTestGoFiles}}{{.ImportPath}}{{end}}' ./... \
							&& go list -f '{{if len .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

checkstyle: ## Run quick checkstyle (govet + goimports (fail on errors))
	@echo "==> Running govet"
	go vet $(PACKAGES) || exit 1
	@echo "==> Running goimports"
	goimports -l -w .
	@echo "==> SUCCESS";

.PHONY: test
test: echo-env ## Run tests
	@echo "==> Testing"
	@for package in $(PACKAGES_WITH_TESTS); do \
		echo "==> Testing ==> $$package" ; \
		go test $$package -test.v || exit 1; \
	done
	@echo "==> SUCCESS"

test-cover: echo-env ## Run tests with -covermode
	@echo "==> Cover Testing"
	rm -fv cover.out;
	rm -fv cprofile.out
	rm -fv cover.out.original;
	@for package in $(PACKAGES_WITH_TESTS); do \
  		echo "==> Cover Testing ==> $$package" ; \
  		go test $$package -test.v -test.coverprofile=cprofile.out -covermode=count || exit 1; \
		go test -c -i -o $(BUILD_PATH)/$$package.test -covermode=count $$package || exit 1; \
		if [ -f cprofile.out ]; then \
			tail -n +2 cprofile.out >> cover.out; \
			rm cprofile.out; \
		fi; \
	done
	@sed -i'.original' '1s;^;mode: count\n;' cover.out
	@echo "==> SUCCESS"

test-race: echo-env ## Run tests with -race
	@echo "==> Testing race conditions"
	@for package in $(PACKAGES_WITH_TESTS); do \
		echo "==> Testing ==> $$package" ; \
		go test -race -run=. -test.timeout=4000s $$package || exit 1; \
	done

BENCH_TIME=10s
test-bench: echo-env ## Run tests with -bench
	@echo "==> Benchmark (-benchtime=${BENCH_TIME})"
	@for package in $(PACKAGES_WITH_TESTS); do \
		echo "==> Testing ==> $$package" ; \
		go test -run=NO_MATCH -bench=. -benchtime=${BENCH_TIME} -benchmem -v $$package || exit 1; \
	done

lint_install: ## Install linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.44.0
	golangci-lint --version

lint: ## Run linter
	@echo "==> Running linter"
	golangci-lint run

help: echo-env
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.DEFAULT_GOAL := help
