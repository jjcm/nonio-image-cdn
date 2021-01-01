package main

import (
	"fmt"
	"net/http"
	"os"
	"soci-image-cdn/route"
)

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./files/images")))
	http.Handle("/thumbnail/", http.StripPrefix("/thumbnail/", http.FileServer(http.Dir("./files/thumbnails"))))
	http.HandleFunc("/upload", route.UploadFile)
	http.HandleFunc("/move", route.MoveFile)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4203"
	}

	fmt.Printf("Listening on %v\n", port)
	http.ListenAndServe(":"+port, nil)
}

func main() {
	fmt.Println("Starting image encoding server...")
	setupRoutes()
}
