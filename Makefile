OUTPUT?=gitstrap
BUILD_VERSION?=0.0
BUILD_HASH?=stub
BUILD_DATE?=2019.01.01

all: clean build test lint

build: $(OUTPUT)


$(OUTPUT):
	go build \
		-ldflags "-X main.buildVersion=$(BUILD_VERSION) -X main.buildCommit=$(BUILD_HASH) -X main.buildDate=$(BUILD_DATE)" \
		-o $(OUTPUT)

clean:
	rm -f $(OUTPUT)

test: $(OUTPUT)
	go test .

lint: $(OUTPUT)
	gometalinter .

.PHONY: build clean test lint

