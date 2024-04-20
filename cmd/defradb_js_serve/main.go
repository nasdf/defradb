package main

// GOOS=js GOARCH=wasm go build -o main.wasm ../defradb_js/main.go

import "net/http"

func main() {
	http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))
}
