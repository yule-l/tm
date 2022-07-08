all: cover

.PHONY: cover
cover:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

.PHONY: update-go-ref
update-go-ref:
	GOPROXY=https://proxy.golang.org GO111MODULE=on go install github.com/yule-l/tm
