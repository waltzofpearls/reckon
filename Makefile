APP = reckon

.PHONY: build
build:
	cargo build --release

.PHONY: run
run:
	cargo run -- --config $(APP).toml --log-level info

.PHONY: lint
lint:
	cargo clippy --workspace --tests --all-features -- -D warnings

.PHONY: test
test:
	cargo test

.PHONY: cover
cover:
	docker run \
		--security-opt seccomp=unconfined \
		-v ${PWD}:/volume \
		xd009642/tarpaulin \
		cargo tarpaulin --out Html --output-dir ./target
	open target/tarpaulin-report.html
