.DEFAULT_GOAL := build

.PHONY: fmt vet test build gen

gen:
	goa gen github.com/oleshko-g/gophermart/api/design -o internal/ 	&& \
	sqlc generate && \
	go generate ./...

test: vet
	go fmt
	go vet
	revive ./...
	go generate ./...
	go test ./...

build:
	go build

fullbuild: test
	go build
