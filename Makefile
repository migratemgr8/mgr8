.PHONY: build
build: main.go
	go build -o bin/mgr8 main.go

.PHONY: test
test:
	go test ./...
