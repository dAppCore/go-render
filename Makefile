.PHONY: wasm test clean

WASM_OUT := dist/go-html.wasm
# Raw size limit: 3MB (Go WASM has ~2MB runtime floor)
WASM_RAW_LIMIT := 3145728
# Gzip transfer size limit: 1MB (what users actually download)
WASM_GZ_LIMIT := 1048576

test:
	go test ./...

wasm: $(WASM_OUT)

$(WASM_OUT):
	@mkdir -p dist
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o $(WASM_OUT) ./cmd/wasm/
	@RAW=$$(stat -c%s "$(WASM_OUT)" 2>/dev/null || stat -f%z "$(WASM_OUT)"); \
	GZ=$$(gzip -c "$(WASM_OUT)" | wc -c); \
	echo "WASM size: $${RAW} bytes raw, $${GZ} bytes gzip"; \
	if [ "$$GZ" -gt $(WASM_GZ_LIMIT) ]; then \
		echo "FAIL: gzip transfer size exceeds 1MB limit ($${GZ} bytes)"; \
		exit 1; \
	elif [ "$$RAW" -gt $(WASM_RAW_LIMIT) ]; then \
		echo "WARNING: raw binary exceeds 3MB ($${RAW} bytes) — check imports"; \
	else \
		echo "OK: gzip $${GZ} bytes (limit 1MB), raw $${RAW} bytes (limit 3MB)"; \
	fi

clean:
	rm -rf dist/
