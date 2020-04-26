// Refactor code to make it easier to test: https://endler.dev/2018/go-io-testing/

package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type credentials struct {
	ExportURI   string `json:"exportUri"`
	ExportToken string `json:"exportToken"`
	ImportURI   string `json:"importUri"`
	ImportToken string `json:"importToken"`
}

var hd, _ = os.UserHomeDir()
var dir = hd + "/.gitlab-migrate"
var fdir = dir + "/config.json"

const unauthorizedError = "401 Unauthorized"

func setup() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("export uri: ")
	exportURI, _ := reader.ReadString('\n')
	exportURI = strings.TrimSuffix(exportURI, "\n")

	fmt.Print("export token: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	exportToken := string(bytePassword)

	// print empty line
	fmt.Println("")

	fmt.Print("import uri: ")
	importURI, _ := reader.ReadString('\n')
	importURI = strings.TrimSuffix(importURI, "\n")

	fmt.Print("import token: ")
	bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	// print empty line
	fmt.Println("")

	importToken := string(bytePassword)

	c := credentials{
		ExportURI:   exportURI,
		ExportToken: exportToken,
		ImportURI:   importURI,
		ImportToken: importToken,
	}

	storeCredentials(c)
	auth(c.ExportURI, c.ExportToken)
	auth(c.ImportURI, c.ImportToken)
}

func storeCredentials(c credentials) error {
	os.Mkdir(dir, 0744)

	c.ExportToken = base64.StdEncoding.EncodeToString([]byte(c.ExportToken))
	c.ImportToken = base64.StdEncoding.EncodeToString([]byte(c.ImportToken))

	bs, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fdir, bs, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("WARNING! Your tokens will be stored unencrypted in  %v\n", fdir)
	return nil
}

func fetchCredentials() credentials {
	bs, err := ioutil.ReadFile(fdir)
	if err != nil {
		log.Fatalf("No credentials found in %v, please run setup", fdir)
	}

	var creds credentials

	json.Unmarshal(bs, &creds)

	de, err := base64.StdEncoding.DecodeString(creds.ExportToken)
	if err != nil {
		log.Fatal(err)
	}
	creds.ExportToken = string(de)

	di, err := base64.StdEncoding.DecodeString(creds.ImportToken)
	if err != nil {
		log.Fatal(err)
	}
	creds.ImportToken = string(di)

	return creds
}

func auth(uri, token string) {
	client := http.Client{}
	url := uri + "/api/v4/projects"

	req, err := http.NewRequest("GET", url, nil)
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

	if m.Message == unauthorizedError {
		log.Println("Login to", uri, "failed:", m.Message)
	} else {
		fmt.Println("Login to", uri, "successful.")
	}
}

func login() {
	if _, err := os.Stat(fdir); os.IsNotExist(err) {
		// path does not exist
		setup()
	} else {
		c := fetchCredentials()
		auth(c.ExportURI, c.ExportToken)
		auth(c.ImportURI, c.ImportToken)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

/* func login(uri, token string) {


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
} */
