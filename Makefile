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

.PHONY: test
test: $(OUTPUT)
	go test ./internal/...

.PHONY: lint
lint: $(OUTPUT)
	golangci-lint run

.PHONY: install
install: $(OUTPUT)
	install $(OUTPUT) /usr/bin

