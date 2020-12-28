.PHONY: build

build:
	go build -v -tags=jsoniter .

build-static:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter -a -o go-there

tests:
	go test -v -tags=jsoniter ./...

integration-tests:
	bash .test/run.sh
