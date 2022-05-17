.PHONY: build
build: main.go
	go build -o bin/mgr8 main.go

.PHONY: test
test:
	go test ./...

release:
	docker run --rm --privileged -e CGO_ENABLED=1 -v /var/run/docker.sock:/var/run/docker.sock -v `pwd`:/go/src/$(PACKAGE_NAME) -v `pwd`/sysroot:/sysroot -w /go/src/$(PACKAGE_NAME) goreleaser/goreleaser-cross:v1.17.6 release