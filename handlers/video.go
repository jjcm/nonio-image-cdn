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

type videoUploadResponse struct {
	Status string
	Name   string
}

// HandleVideo encodes the video into a webm and returns the path to it
func HandleVideo(w http.ResponseWriter, r *http.Request, file multipart.File, name string) {
	// Create a temp file
	tempFile, err := ioutil.TempFile("files/temp-videos", "video-*.mp4")
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

	// since this is a video we'll use ffmpeg to encode it
	cmd := exec.Command("ffmpeg", "-y", "-i", tempFile.Name(), "-c:v", "libvpx-vp9", "-b:v", "2M", fmt.Sprintf("files/videos/%v.webm", name))
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
	res := videoUploadResponse{"success", fmt.Sprintf("videos/%v.webm", name)}
	SendResponse(w, res, 200)
}
