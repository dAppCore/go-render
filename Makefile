.PHONY: wasm test clean

WASM_OUT := dist/go-html.wasm

test:
	go test ./...

wasm: $(WASM_OUT)

$(WASM_OUT):
	@mkdir -p dist
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o $(WASM_OUT) ./cmd/wasm/
	@ls -lh $(WASM_OUT)
	@SIZE=$$(stat -c%s "$(WASM_OUT)" 2>/dev/null || stat -f%z "$(WASM_OUT)"); \
	if [ "$$SIZE" -gt 2097152 ]; then \
		echo "WARNING: WASM binary exceeds 2MB target ($${SIZE} bytes)"; \
	else \
		echo "OK: WASM binary within 2MB target ($${SIZE} bytes)"; \
	fi

clean:
	rm -rf dist/
