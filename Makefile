.EXPORT_ALL_VARIABLES:
SHELL := /bin/bash -Eeuo pipefail

.PHONY: build
build:
	CGO_ENABLED=0 go build -o bin/getlush .