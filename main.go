package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"soci-cdn/handlers"
	"soci-cdn/util"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		util.SendResponse(w, "", 200)
		return
	}
	// Parse our multipart form, set a 1GB max upload size
	r.ParseMultipartForm(1 << 30)

	// Get the user's email if we're authorized
	bearerToken := r.Header.Get("Authorization")
	fmt.Println(bearerToken)
	user, err := util.GetUserEmail(bearerToken)
	fmt.Println(user)
	if err != nil {
		util.SendError(w, fmt.Sprintf("User is not authorized. Token: %v", bearerToken), 400)
		return
	}

	// Parse our url, and check if the url is available
	url := r.FormValue("url")
	urlIsAvailable, err := util.CheckIfURLIsAvailable(url)
	if err != nil {
		util.SendError(w, "Error checking requested url.", 500)
		return
	}
	if urlIsAvailable == false {
		util.SendError(w, "Url is taken.", 400)
		return
	}

	// Parse our file and assign it to the proper handlers depending on the type
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	re, _ := regexp.Compile("([a-zA-Z]+)/")
	var mimeType = handler.Header["Content-Type"][0]

	// If all is good, let's log what the hell is going on
	fmt.Printf("%v is uploading a %v of size %v to %v", user, re.FindStringSubmatch(mimeType)[1], handler.Size, url)

	switch re.FindStringSubmatch(mimeType)[1] {
	case "image":
		handlers.HandleImage(w, r, file, url)
	case "video":
		handlers.HandleVideo(w, r, file, url)
	}

}

// Move takes the temp file and renames it to match the url
func moveFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		util.SendResponse(w, "", 200)
		return
	}

	r.ParseMultipartForm(1 << 30)

	// Get the user's email if we're authorized
	bearerToken := r.Header.Get("Authorization")
	fmt.Println(bearerToken)
	user, err := util.GetUserEmail(bearerToken)
	fmt.Println(user)
	if err != nil {
		util.SendError(w, "User is not authorized.", 400)
		return
	}

	// Parse our url, and check if the url is available
	url := r.FormValue("url")
	urlIsAvailable, err := util.CheckIfURLIsAvailable(url)
	if err != nil {
		util.SendError(w, "Error checking requested url.", 500)
		fmt.Println(err)
		return
	}
	if urlIsAvailable == false {
		util.SendError(w, "Url is taken.", 400)
		return
	}

	// Check if the file we're moving exists
	tempFile := r.FormValue("tempName")
	if _, err := os.Stat(fmt.Sprintf("files/temp-images/%v.webp", tempFile)); os.IsNotExist(err) {
		util.SendError(w, "No temp image exists with that name.", 400)
		fmt.Println(err)
		return
	}

	// If everything else looks good, lets move the file.
	err = os.Rename(fmt.Sprintf("files/temp-images/%v.webp", tempFile), fmt.Sprintf("files/images/%v.webp", url))
	if err != nil {
		util.SendError(w, "Error renaming file.", 500)
		return
	}

	// Send back a response.
	util.SendResponse(w, url, 200)
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/move", moveFile)
	http.ListenAndServe(":8081", nil)
}

func main() {
	fmt.Println("Starting media encoding server")
	setupRoutes()
}
