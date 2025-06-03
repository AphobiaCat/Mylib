
all:
	cargo build -j 8 -p main

exsamples:
	cargo build -j 8 -p exsamples

release:
	cargo build -j 8 -p main --release

check:
	cargo check -p main

run:
	cargo run -p main

wasm:
	cargo build -j 8 -p main_wasm --target wasm32-unknown-unknown --release
	wasm-bindgen target/wasm32-unknown-unknown/release/main_wasm.wasm --out-dir pkg --target web
clean:
	cargo clean
	rm -rf ./pkg
