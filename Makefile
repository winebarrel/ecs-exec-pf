SHELL   := /bin/bash
VERSION := v0.1.0
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: all
all: vet build

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/ecs-exec-pf

.PHONY: vet
vet:
	go vet

.PHONY: package
package: clean vet build
	gzip ecs-exec-pf -c > ecs-exec-pf_$(VERSION)_$(GOOS)_$(GOARCH).gz
	sha1sum ecs-exec-pf_$(VERSION)_$(GOOS)_$(GOARCH).gz > ecs-exec-pf_$(VERSION)_$(GOOS)_$(GOARCH).gz.sha1sum

.PHONY: clean
clean:
	rm -f ecs-exec-pf
