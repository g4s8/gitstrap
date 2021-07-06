OUTPUT?=gitstrap
BUILD_VERSION?=0.0
BUILD_HASH?=stub
BUILD_DATE?=2019.01.01

all: clean build test lint

.PHONY: build
build: $(OUTPUT)

$(OUTPUT):
	go build \
		-ldflags "-X main.buildVersion=$(BUILD_VERSION) -X main.buildCommit=$(BUILD_HASH) -X main.buildDate=$(BUILD_DATE)" \
		-o $(OUTPUT) ./cmd/gitstrap

.PHONY: clean
clean:
	rm -f $(OUTPUT)

# run_tests_dir - run all tests in provided directory
define _run_tests_dir
  go test -v ${TEST_OPTS} "./$(1)/..."
endef

.PHONY: test
test: $(OUTPUT)
	$(call _run_tests_dir,internal)

.PHONY: test-race
test-race: TEST_OPTS := ${TEST_OPTS} -race
test-race: test

.PHONY: bench
bench: TEST_OPTS := ${TEST_OPTS} -bench=. -run=^$
bench: test

.PHONY: lint
lint: $(OUTPUT)
	golangci-lint run

.PHONY: install
install: $(OUTPUT)
	install $(OUTPUT) /usr/bin

