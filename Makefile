build_go:
	docker run --rm \
	-v `pwd`/src:/game \
	--env GOOS=js --env GOARCH=wasm \
	golang:1.12-rc \
	/bin/bash -c "go build -o /wasm-example/main.wasm /wasm-example/main.go; cp /usr/local/go/misc/wasm/wasm_exec.js /wasm_exec.js"

serve:
	(sleep 2 && open http://localhost:8080) &
	docker run --rm -p 8080:8043 -v `pwd`/src:/srv/http pierrezemb/gostatic:latest
