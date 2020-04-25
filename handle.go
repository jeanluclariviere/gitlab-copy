package main

import (
	"log"
	"net/http"
)

func handle() {
	client := http.Client{}
	uri, token := fetchCredentials()

	//req := exportStatus(uri, "1")
	path := "/home/jeanluc/go/src/github.com/jeanluclariviere/gitlab-migrate/project.tar.gz"
	extraParams := map[string]string{
		"path":      "upload-test",
		"overwrite": "false",
	}

	req := importFile(uri, extraParams, "file", path)
	req.Header.Add("PRIVATE-TOKEN", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	printResponse(resp)
}
