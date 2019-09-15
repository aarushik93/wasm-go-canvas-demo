# Compiling
 `GOOS=js GOARCH=wasm go build -o main.wasm main.go`
 
# Running
goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))'

# Explaination
This is the first demo with wasbassembly and GoLang for GoLang Conf. Main take away from this is:
a) Everything can be done in Go, no JS, but it ends up looking like ugly JS. 
b) Super slow because we have to copy every JS value in array to GO value, GO 1.13 provides a better way. 
 
