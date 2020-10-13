.PHONY: build

build:
	go build -v -tags=jsoniter .

test:
	go test -v -tags=jsoniter ./...