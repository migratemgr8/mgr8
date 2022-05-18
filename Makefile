PACKAGE_NAME         := github.com/kenji-yamane/mgr8
GOLANG_CROSS_VERSION ?= v1.17.6

.PHONY: build
build: main.go
	go build -o bin/mgr8 main.go

.PHONY: test
test:
	go test ./... -coverprofile=unit_coverage.out

.PHONY: coverage-report
coverage-report:
	go tool cover -html=coverage.out

UNIT_COVERAGE:= $(shell go tool cover -func=unit_coverage.out | tail -n 1 | cut -d ' ' -f 3 | rev | cut -c 1-5 | rev)

.PHONY: display-coverage
display-coverage:
	@echo "Unit Coverage: $(UNIT_COVERAGE)"

.PHONY: release
release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist

