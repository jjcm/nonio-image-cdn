package route

import (
	"fmt"
	"net/http"
	"os"
	"soci-cdn/util"
)

// MoveFile takes the temp file and renames it to match the url
func MoveFile(w http.ResponseWriter, r *http.Request) {
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
	if _, err := os.Stat(fmt.Sprintf("files/images/%v.webp", tempFile)); os.IsNotExist(err) {
		util.SendError(w, "No temp image exists with that name.", 400)
		fmt.Println(err)
		return
	}

	// If everything else looks good, lets move the file.
	err = os.Rename(fmt.Sprintf("files/images/%v.webp", tempFile), fmt.Sprintf("files/images/%v.webp", url))
	if err != nil {
		util.SendError(w, "Error renaming file.", 500)
		return
	}

	// Send back a response.
	util.SendResponse(w, url, 200)
}
