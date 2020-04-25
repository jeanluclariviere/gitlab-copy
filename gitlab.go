package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const projects = "/api/v4/projects/"

func scheduleExport(uri, pid string) *http.Request {

	url := uri + projects + pid + "/export/"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return req
}

func exportStatus(uri, pid string) *http.Request {
	url := uri + projects + pid + "/export/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return req
}

func exportDownload(uri, pid string) *http.Request {

	url := uri + projects + pid + "/export/download"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return req
}

func importFile(uri string, params map[string]string, paramName, path string) *http.Request {

	url := uri + projects + "/import"

	// Open a file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	// Close the file on exit
	defer file.Close()

	//Create a variable sized buffer
	//of type io writer
	body := &bytes.Buffer{}

	// Create a multipart writer
	// that writes to body of type io.writer
	writer := multipart.NewWriter(body)

	// Create a new form-data header with the provided field name and file name
	// And returns an io.writer
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		log.Fatal(err)
	}

	// Copy the contents of the file directly in to part
	// the io.writer
	_, err = io.Copy(part, file)

	// Adds additional fields to our writer
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	// Close the writer now that we are finished
	// adding contents
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new request and return it
	req, err := http.NewRequest("POST", url, body)
	// You must set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
