package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"soci-cdn/handlers"
)

type uploadResponse struct {
	Status string
	Name   string
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		handlers.SendResponse(w, "", 200)
		return
	}
	// Parse our multipart form, set a 1GB max upload size
	r.ParseMultipartForm(1 << 30)

	// Parse our url, and check if the url is available
	url := r.FormValue("url")
	urlIsAvailable, err := checkIfURLIsAvailable(url)
	if err != nil {
		res := uploadResponse{"Error checking requested url", url}
		handlers.SendResponse(w, res, 400)
		return
	}
	if urlIsAvailable == false {
		res := uploadResponse{"Url is taken", url}
		handlers.SendResponse(w, res, 400)
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
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)

	re, _ := regexp.Compile("([a-zA-Z]+)/")
	var mimeType = handler.Header["Content-Type"][0]
	fmt.Println(re.FindStringSubmatch(mimeType)[1])

	switch re.FindStringSubmatch(mimeType)[1] {
	case "image":
		handlers.HandleImage(w, r, file, url)
	case "video":
		handlers.HandleVideo(w, r, file, url)
	}

}

// hits our api server and checks to see if the url is available to upload to
func checkIfURLIsAvailable(url string) (bool, error) {
	urlCheckRes, err := http.Get(fmt.Sprintf("https://api.non.io/posts/url-is-available/%v", url))
	if err != nil {
		fmt.Println("Error checking if url is available")
		fmt.Println(err)
		return false, err
	}
	defer urlCheckRes.Body.Close()
	body, err := ioutil.ReadAll(urlCheckRes.Body)
	if err != nil {
		fmt.Println("Error parsing the body of the url check")
		fmt.Println(err)
		return false, err
	}

	return string(body) == "true", err
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8081", nil)
}

func main() {
	fmt.Println("Starting media encoding server")
	setupRoutes()
}
