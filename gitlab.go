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
	"strings"
)

const projects = "/api/v4/projects/"
const groups = "/api/v4/groups/"
const alreadyTakenError = `Failed to save group {:path=>["has already been taken"]}`

func newGroup(uri, token, path string) string {
	URL := uri + groups
	ps := strings.Split(path, "/")
	data := url.Values{}
	client := http.Client{}
	pid := ""
	for i, v := range ps {
		data.Set("name", v)
		data.Set("path", v)
		data.Set("parent_id", pid)

		req, err := http.NewRequest("POST", URL, bytes.NewBufferString(data.Encode()))
		req.Header.Add("PRIVATE-TOKEN", token)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		r := struct {
			ID      int    `json:"id"`
			Message string `json:"message"`
		}{}

		json.Unmarshal(bs, &r)

		if r.Message == "" {
			pid = strconv.Itoa(r.ID)
		} else if r.Message == alreadyTakenError {
			fp := strings.Join(ps[:i+1], "/")
			pid = getParentID(uri, token, v, fp)
		}
	}

	return pid
}

// getParentID return the parent ID the last group in the path
// the full path must be provided
func getParentID(uri, token, name, path string) string {
	client := http.Client{}
	URL := uri + groups + "?search=" + name

	req, err := http.NewRequest("GET", URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := []struct {
		ID       int    `json:"id"`
		FullPath string `json:"full_path"`
	}{}

	if err := json.Unmarshal(bs, &r); err != nil {
		log.Fatal(err)
	}

	for _, v := range r {
		if v.FullPath == path {
			return strconv.Itoa(v.ID)
		}
	}

	return ""
}

func scheduleExport(uri, token, pid string) *http.Response {
	client := http.Client{}

	URL := uri + projects + pid + "/export/"
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func exportStatus(uri, token, pid string) *http.Response {
	client := http.Client{}

	URL := uri + projects + pid + "/export/"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func exportDownload(uri, token, pid string) *http.Response {
	client := http.Client{}

	URL := uri + projects + pid + "/export/download"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func importFile(uri, token string, params map[string]string, path string) *http.Response {
	client := http.Client{}

	URL := uri + projects + "/import"

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
	req, err := http.NewRequest("POST", URL, body)
	// You must set the content type
	req.Header.Add("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}
