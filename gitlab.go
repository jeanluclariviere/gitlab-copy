//{"message":"404 Namespace Not Found"}

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

const projects = "/api/v4/projects/"
const groups = "/api/v4/groups/"

func newGroup(uri string, params map[string]string) *http.Request {
	uri = uri + groups

	data := url.Values{}
	for k, v := range params {
		data.Add(k, v)
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	return req
}

func searchGroup(uri string, name string) string {
	uri = uri + groups + "?search=" + name

	resp, err := http.Post("POST", uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := struct {
		ID          int
		name        string
		path        string
		description string
	}{}

	if err := json.Unmarshal(bs, &r); err != nil {
		log.Fatal(err)
	}

	return strconv.Itoa(r.ID)
}

func scheduleExport(uri, pid string) *http.Request {
	uri = uri + projects + pid + "/export/"
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	return req
}

func exportStatus(uri, pid string) *http.Request {
	uri = uri + projects + pid + "/export/"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	return req
}

func exportDownload(uri, pid string) *http.Request {
	uri = uri + projects + pid + "/export/download"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	return req
}

func importFile(uri string, params map[string]string, path string) *http.Request {
	uri = uri + projects + "/import"

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
	part, err := writer.CreateFormFile("file", filepath.Base(path))
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
	req, err := http.NewRequest("POST", uri, body)
	// You must set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
