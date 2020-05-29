package main

import (
	"encoding/json"
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

type authorizationResponse struct {
	Error string
	Email string
	ID    int
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		handlers.SendResponse(w, "", 200)
		return
	}
	// Parse our multipart form, set a 1GB max upload size
	r.ParseMultipartForm(1 << 30)

	// Get the user's email if we're authorized
	bearerToken := r.Header.Get("Authorization")
	fmt.Println(bearerToken)
	user, err := getUserEmail(bearerToken)
	fmt.Println(user)
	if err != nil {
		fmt.Println("User is not authorized.")
		fmt.Println(err)
		res := uploadResponse{"User is not authorized", bearerToken}
		handlers.SendResponse(w, res, 400)
		return
	}

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

func getUserEmail(bearerToken string) (string, error) {
	// Send a req to the api to get the email from our token, if it's valid
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.non.io/protected", nil)
	req.Header.Add("Authorization", bearerToken)
	userAuthRes, err := client.Do(req)
	if err != nil {
		fmt.Println("Error checking if the user is authorized")
		fmt.Println(err)
		return "", err
	}
	defer userAuthRes.Body.Close()

	// Parse the body of the request once it comes back
	body, err := ioutil.ReadAll(userAuthRes.Body)
	if err != nil {
		fmt.Println("Error parsing the body of the user auth check")
		fmt.Println(err)
		return "", err
	}

	// Create a authResponse struct, fill it with the parsed json values
	authResponse := authorizationResponse{}
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		fmt.Println("Error parsing the json of the user auth check")
		fmt.Println(err)
		return "", err
	}

	// Populate our error with the json response's error
	if authResponse.Error != "" {
		err = fmt.Errorf(authResponse.Error)
	}

	return authResponse.Email, err
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8081", nil)
}

func main() {
	fmt.Println("Starting media encoding server")
	setupRoutes()
}
