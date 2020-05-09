package main

import (
	"fmt"
	"net/http"
	"regexp"
	"soci-cdn/handlers"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		handlers.SendResponse(w, "", 200)
		return
	}
	// Parse our multipart form, set a 1GB max upload size
	r.ParseMultipartForm(1 << 30)

	name := r.FormValue("name")
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)

	re, _ := regexp.Compile("([a-zA-Z]+)/")
	var mimeType = handler.Header["Content-Type"][0]
	fmt.Println(re.FindStringSubmatch(mimeType)[1])

	switch re.FindStringSubmatch(mimeType)[1] {
	case "image":
		handlers.HandleImage(w, r, file, name)
	case "video":
		handlers.HandleVideo(w, r, file, name)
	}

}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8081", nil)
}

func main() {
	fmt.Println("Starting media encoding server")
	setupRoutes()
}
