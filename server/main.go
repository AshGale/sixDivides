package main

import (
	"embed"
	"fmt"
	"net/http"
)

// https://golangbot.com/webassembly-using-go/
// https://go-app.dev/getting-started#build-the-server
// go install golang.org/x/mobile/cmd/gomobile@latest
// https://pkg.go.dev/golang.org/x/mobile/app
/* powershell
$Env:GOOS = "windows"; $Env:GOARCH = "amd64"; go run main.go

$Env:GOOS = "darwin"; $Env:GOARCH = "amd64"; go build -o mac.dmg main.go 		// mac
$Env:GOOS = "linux"; $Env:GOARCH = "arm64"; go build -o android.apk main.go 	// android
$Env:GOOS = "windows"; $Env:GOARCH = "amd64"; go build -o windows.exe main.go 	// windows
$Env:GOOS = "js"; $Env:GOARCH = "wasm"; go build -o browser.wasm main.go 		// browser


*/

//go:embed static/*
var staticFiles embed.FS

func main() {
	port := "9090"

	done := make(chan bool)
	//go http.ListenAndServe(fmt.Sprintf(":%v", port), http.FileServer(http.Dir("../")))
	http.Handle("/", http.FileServer(http.FS(staticFiles)))
	go http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	fmt.Printf("Server started at port %v\n", port)
	<-done
}
