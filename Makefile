PACKAGE_NAME         := github.com/migratemgr8/mgr8
GOLANG_CROSS_VERSION ?= v1.18.2

.PHONY: build
build: main.go
	go build -o bin/mgr8 main.go

install-tools:
	go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: test
test:
	go test ./... -coverprofile=coverage.txt

.PHONY: coverage-report
coverage-report:
	go tool cover -html=coverage.txt

UNIT_COVERAGE:= $(shell go tool cover -func=coverage.txt | tail -n 1 | cut -d ' ' -f 3 | rev | cut -c 1-5 | rev)

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

.PHONY: mock
mock:
	@mockgen -source=domain/driver.go -destination=mock/domain/driver_mock.go -package=domain_mock
	@mockgen -source=domain/diff_deque.go -destination=mock/domain/diff_deque_mock.go -package=domain_mock
	@mockgen -source=infrastructure/clock.go -destination=mock/infrastructure/clock_mock.go -package=infrastructure_mock
	@mockgen -source=infrastructure/file.go -destination=mock/infrastructure/file_mock.go -package=infrastructure_mock
	@mockgen -source=applications/migrationscripts.go -destination=mock/applications/migrationscripts_mock.go -package=applications_mock
	@mockgen -source=applications/log.go -destination=mock/applications/log_mock.go -package=applications_mock
