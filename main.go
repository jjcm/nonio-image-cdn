package main

import (
	"fmt"
	"net/http"
	"os"
	"soci-image-cdn/route"
)

func setupRoutes() {
<<<<<<< HEAD
	http.Handle("/", http.FileServer(http.Dir("./files/images")))
	http.Handle("/thumbnail/", http.StripPrefix("/thumbnail/", http.FileServer(http.Dir("./files/thumbnails"))))
=======
	http.Handle("/", http.FileServer(http.Dir("./files")))
>>>>>>> 3fd91a1b4c0aeef3453b2289d1f56510b9fd557f
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
