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
var dir = hd + "/.gitlab-copy"
var fdir = dir + "/config.json"

const unauthorizedError = "401 Unauthorized"

// setup is takes user input and populates the config.json
// used to authenticate against the source and destination
// gitlab servers.
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

// storeCredentails stores the source and destination URIs
// and tokens (unencrypted!) in ~/.gitlab-copy/config.json
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

// fetchCredentials reads the ~/.gitlab-copy/config.json file
// and returns the credentials to be used.
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

// auth validates a given URI and token combination
// by attempting to list all projects
// returns success if error 401 is not encountered.
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

// login is a high level wrapper around auth which
// simply validates the existing credentials
// useful to check if the tokens have not expired.
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

// exists is a small helper function
// which validates whether or not
// a path/file exists.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}
