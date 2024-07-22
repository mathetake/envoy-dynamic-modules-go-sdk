goimports := golang.org/x/tools/cmd/goimports@v0.21.0
golangci_lint := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0

.PHONY: build
build:
	@go build ./...
	@cd example && go build -buildmode=c-shared -o main .

.PHONY: test
test:
	@go test $(shell go list ./... | grep -v e2e)
	@cd example && CGO_ENABLED=0 go test ./... -count=1

.PHONY: lint
lint:
	@echo "lint => ./..."
	@go run $(golangci_lint) run ./...
	@echo "lint => example/"
	@cd example && go run $(golangci_lint) run ./...

.PHONY: format
format:
	@find . -type f -name '*.go' | xargs gofmt -s -w
	@for f in `find . -name '*.go'`; do \
	    awk '/^import \($$/,/^\)$$/{if($$0=="")next}{print}' $$f > /tmp/fmt; \
	    mv /tmp/fmt $$f; \
	done
	@go run $(goimports) -w -local github.com/envoyproxyx/go-sdk `find . -name '*.go'`

.PHONY: tidy
tidy: ## Runs go mod tidy on every module
	@find . -name "go.mod" \
	| grep go.mod \
	| xargs -I {} bash -c 'dirname {}' \
	| xargs -I {} bash -c 'echo "tidy => {}"; cd {}; go mod tidy -v; '

.PHONY: precommit
precommit: format lint tidy

.PHONY: check
check:
	@$(MAKE) precommit
	@if [ ! -z "`git status -s`" ]; then \
		echo "The following differences will fail CI until committed:"; \
		git diff --exit-code; \
	fi
