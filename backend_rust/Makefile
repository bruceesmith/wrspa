.PHONY: all audit build run test fmt fmt-fix clippy wasm

all: build

## audit: run quality control checks
.PHONY: audit
audit: fmt clippy test

## build: build the application
.PHONY: build
build:
	cargo build

## wasm: build the application for WebAssembly (WASI)
# NOTE: Currently fails due to tokio features incompatibility on wasm32-wasip1
.PHONY: wasm
wasm:
	cargo build --target wasm32-wasip1

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
