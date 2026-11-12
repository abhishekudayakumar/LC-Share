
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed web/*
var webFiles embed.FS

func main() {
	web, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(web)))

	address := "127.0.0.1:8765"

	fmt.Println("================================")
	fmt.Println("          LC SHARE")
	fmt.Println("================================")
	fmt.Println("Dashboard: http://" + address)
	fmt.Println("Status: Running")

	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
