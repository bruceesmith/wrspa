.PHONY: all audit build run test fmt fmt-fix clippy

all: build

## audit: run quality control checks
.PHONY: audit
audit: fmt clippy test

## build: build the application
.PHONY: build
build:
	cargo build

clippy:
	cargo clippy -- -D warnings

## coverage: run tests and show code coverage
coverage:
	cargo llvm-cov test

## fmt: check the source formatting
fmt:
	cargo fmt -- --check

## fmt-fix: format the si=ource
fmt-fix:
	cargo fmt

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

run:
	cargo run

## test: run the test suite
test:
	cargo test
