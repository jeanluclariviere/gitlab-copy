// Refactor code to make it easier to test: https://endler.dev/2018/go-io-testing/

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type credentials struct {
	URI   string `json:"uri"`
	Token string `json:"token"`
}

func storeCredentials(uri, token string) error {
	hd, _ := os.UserHomeDir()
	dir := hd + "/.gitlab-migrate"
	fdir := dir + "/config.json"

	os.Mkdir(dir, 0744)

	token = base64.StdEncoding.EncodeToString([]byte(token))

	creds := credentials{
		URI:   uri,
		Token: token,
	}
	bs, err := json.MarshalIndent(creds, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fdir, bs, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("WARNING! Your token will be stored unencrypted in  %v\n", fdir)
	return nil
}

func fetchCredentials() (string, string) {
	hd, _ := os.UserHomeDir()
	dir := hd + "/.gitlab-migrate"
	fdir := dir + "/config.json"

	bs, err := ioutil.ReadFile(fdir)
	if err != nil {
		log.Fatalf("No credentials found in %v, please login", fdir)
	}

	var creds credentials

	json.Unmarshal(bs, &creds)

	dt, err := base64.StdEncoding.DecodeString(creds.Token)
	if err != nil {
		log.Fatal(err)
	}

	return creds.URI, string(dt)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func login(uri, token string) {
	hd, _ := os.UserHomeDir()
	dir := hd + "/.gitlab-migrate"
	fdir := dir + "/config.json"

	if _, err := os.Stat(fdir); err == nil {
		// path exists
		if uri != "" && token != "" {
			storeCredentials(uri, token)
			uri, token = fetchCredentials()
		} else {
			uri, token = fetchCredentials()
		}
	} else if os.IsNotExist(err) {
		// path does not exist
		storeCredentials(uri, token)
		uri, token = fetchCredentials()
	}

	auth(uri, token)
}

func auth(uri, token string) {
	client := http.Client{}
	u := uri + "/api/v4/projects"
	req, err := http.NewRequest("GET", u, nil)

	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	m := struct {
		Message string `json:"message"`
	}{}

	json.Unmarshal(body, &m)

	if m.Message == "401 Unauthorized" {
		log.Fatalf("Login failed, %v", m.Message)
	} else {
		fmt.Println("Login successful")
	}
}
