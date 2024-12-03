
LOCALS := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PKGS   := log $(wildcard *util)
COUNT  ?= 1

# TEST_GOSTOCKUTIL_RETRIEVE_VIA_SFTP ?= sftp://cortex/motd
# TEST_GOSTOCKUTIL_RETRIEVE_VIA_SSH  ?= ssh://cortex/hostname

.PHONY: test deps docs $(PKGS)

.EXPORT_ALL_VARIABLES:
GO111MODULE = on

all: fmt deps test docs

deps:
	which peg || go install github.com/pointlander/peg@latest
	go generate -x ./...
	go get ./...
	-go mod tidy

fmt:
	$(info Formatting)
	@gofmt -w $(LOCALS)
	$(info Vetting code)
	@go vet ./...

docs:
	@owndoc render --property rootpath=/go-stockutil/

$(PKGS):
	$(info Testing $(@))
	@go test -count=$(COUNT) ./$(@)/...

test: $(PKGS)

