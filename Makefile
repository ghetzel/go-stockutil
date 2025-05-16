
LOCALS      := $(shell find -type f -name '*.go')
PKGS        := log $(wildcard *util)
COUNT       ?= 1
GO111MODULE  = on

all: fmt deps test docs

deps:
	which peg || go install github.com/pointlander/peg@latest
	go generate -x ./...
	go get ./...
	-go mod tidy

fmt: gofmt
	@go vet ./...

gofmt: $(LOCALS)
$(LOCALS):
	@gofmt -w $(@)

docs:
	@owndoc render --property rootpath=/go-stockutil/

$(PKGS):
	$(info Testing $(@))
	@go test -count=$(COUNT) ./$(@)/...

test: $(PKGS)


.PHONY: test deps docs $(PKGS) $(LOCALS)