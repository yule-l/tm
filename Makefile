all: build

.PHONY: cover ## Run coverage
cover:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

.PHONY: update-go-ref ## Update go ref
update-go-ref:
	GOPROXY=https://proxy.golang.org GO111MODULE=on go install github.com/yule-l/tm

.PHONY: build
build: ## Build application
	go build -o bin/tm cmd/tm/main.go
