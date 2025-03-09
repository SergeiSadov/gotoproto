build-srv:
	go build -o app cmd/srv/main.go

build-wasm:
	GOOS=js GOARCH=wasm go build -o docs/main.wasm cmd/wasm/main.go

serve: build-wasm build-srv
	./app