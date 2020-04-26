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

type groupResp struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func newGroup(uri, token, path string) (groupID string) {
	client := http.Client{}
	URL := uri + groups

	// Create the path slice
	ps := strings.Split(path, "/")
	data := url.Values{}
	var pid string

	// Range over the path slice and create each group.
	// If the group already exists, retrieve it's ID
	// and continue until the full path has been constructed.
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

		// ID is returned if the object exists
		// Message is returned in the event of an error,
		// including if the group already exists.
		var r groupResp

		json.Unmarshal(bs, &r)

		if r.Message == "" {
			// If the group gets created
			// set the current groups ID as the parent
			pid = strconv.Itoa(r.ID)
		} else if r.Message == alreadyTakenError {
			// If the current group already exists, rebuild the
			// full_path up to the current value
			// and retrieve it's ID
			fp := strings.Join(ps[:i+1], "/")
			pid = getParentID(uri, token, v, fp)
		} else {
			log.Fatal(r.Message)
		}
	}

	return pid
}

type parentResp struct {
	ID       int    `json:"id"`
	FullPath string `json:"full_path"`
}

// getParentID return the parent ID the last group in the path
// the full path must be provided
func getParentID(uri, token, name, path string) (parentID string) {
	client := http.Client{}
	URL := uri + groups + "?search=" + name

	// Perform the search
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

	// API response struct, condensed.
	var r []parentResp

	// Unmarshal search response
	if err := json.Unmarshal(bs, &r); err != nil {
		log.Fatal(err)
	}

	// Return the ID from list of responses
	// where the requested path is equal to the groups full_path
	for _, v := range r {
		if v.FullPath == path {
			return strconv.Itoa(v.ID)
		}
	}

	return ""
}

func scheduleExport(uri, token, pid string) (*http.Response, error) {
	client := http.Client{}

	URL := uri + projects + pid + "/export/"
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)

	return resp, err
}

type statusResp struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	ExportStatus string `json:"export_status"`
}

func exportStatus(uri, token, pid string) (*statusResp, *http.Response, error) {
	client := http.Client{}

	URL := uri + projects + pid + "/export/"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp, err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var r statusResp

	json.Unmarshal(bs, &r)

	return &r, resp, nil
}

func exportDownload(uri, token, pid, filename string) (*http.Response, error) {
	client := http.Client{}
	URL := uri + projects + pid + "/export/download"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	return resp, nil
}

func importFile(uri, token, namespace, path, filename string) *http.Response {
	client := http.Client{}

	URL := uri + projects + "/import"

	// Open a file
	file, err := os.Open(filename)
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
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		log.Fatal(err)
	}

	// Copy the contents of the file directly in to part
	// the io.writer
	_, err = io.Copy(part, file)

	// Add additional fields to our writer

	if err = writer.WriteField("path", path); err != nil {
		log.Fatal(err)
	}

	if err = writer.WriteField("namespace", namespace); err != nil {
		log.Fatal(err)
	}

	if err = writer.WriteField("file", filename); err != nil {
		log.Fatal(err)
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
