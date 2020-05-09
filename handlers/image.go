package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
)

type imageUploadResponse struct {
	Status string
	Name   string
}

// HandleImage encodes the image into a webp and returns the path to it
func HandleImage(w http.ResponseWriter, r *http.Request, file multipart.File, name string) {
	// Create a temp file
	tempFile, err := ioutil.TempFile("temp-images", "image-*")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read the uploaded file into a buffer and write it to our temp file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	// since this is an image we'll use magick to encode it
	cmd := exec.Command("magick", tempFile.Name(), fmt.Sprintf("files/images/%v.webp", name))
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Dir = workingDir
	var output bytes.Buffer
	cmd.Stderr = &output
	err = cmd.Run()
	if err != nil {
		panic(output.String())
	}

	// if everything looks good, send back a response
	res := imageUploadResponse{"success", fmt.Sprintf("files/images/%v.webp", name)}
	SendResponse(w, res, 200)
}
