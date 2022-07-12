
BUILD_PATH="./build"

GOBIN ?= $(shell go env GOPATH)/bin

CMD_NAME="omecmd"
CMD_PATH="$(GOBIN)/$(CMD_NAME)"

prepare:
	mkdir -p $(BUILD_PATH)

echo-env: ## Echo environment variables
	@echo go version $(shell go version)
	@echo GOPATH: $(GOPATH)
	@echo GOROOT: $(GOROOT)
	@echo GOBIN:  $(GOBIN)
	@echo GOOS: $(GOOS)
	@echo GO111MODULE: $(GO111MODULE)
	@echo PKG_CONFIG_PATH: $(PKG_CONFIG_PATH)
	@echo AWS_PROFILE: $(AWS_PROFILE)
	@echo CLIENT_PATH: $(CMD_PATH)
	@echo BUILD_PATH: $(BUILD_PATH)

.PHONY: clean
clean: echo-env ## Clean
	@echo "==> Clean"
	rm -rf $(BUILD_PATH)
	rm -fv cover.out
	rm -fv cprofile.out
	rm -fv cover.out.original
	rm -fv $(CMD_PATH)
	env GO111MODULE=off go clean -cache -testcache

.PHONY: nuke
nuke: clean ## Clean with -modcache
	@echo "==> Nuke-em!"
	env GO111MODULE=off go clean -modcache -x

.PHONY: build
build: echo-env prepare ## Build binary
	@echo "==> Build"
	GOOS=linux
	go build -a -trimpath -ldflags="-s -w" -o $(BUILD_PATH)/app ./

.PHONY: install
install: echo-env ## Install command line tool
	@echo "==> Install"
	go install -trimpath -ldflags="-s -w" ./cmd/...

PACKAGES=$(shell go list ./...)
PACKAGES_WITH_TESTS = $(shell go list -f '{{if len .XTestGoFiles}}{{.ImportPath}}{{end}}' ./... \
							&& go list -f '{{if len .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

.PHONY: test
test: echo-env ## Run tests
	@for package in $(PACKAGES_WITH_TESTS); do \
		echo "==> Testing ==> $$package" ; \
		go test $$package -test.v || exit 1; \
	done
	@echo "==> SUCCESS"

test-cover: echo-env ## Run tests with -covermode
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
	@for package in $(PACKAGES_WITH_TESTS); do \
		echo "==> Testing race ==> $$package" ; \
		go test -race -run=. -test.timeout=4000s $$package || exit 1; \
	done

BENCH_TIME=1s
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
	golangci-lint run

buf-gen: ## Run buf lint format generate
	buf export buf.build/envoyproxy/protoc-gen-validate -o .
	# Not tracking breaking changes just yet, uncomment after first major release.
	# buf breaking --against '.git#branch=main'
	buf lint --exclude-path validate/validate.proto
	buf format -w
	buf mod update api
	buf generate

buf-ls: ## Run buf list
	buf build -o -#format=json | jq '.file[] | .package' | sort | uniq | head

help: echo-env
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.DEFAULT_GOAL := help
