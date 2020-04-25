// Refactor code to make it easier to test: https://endler.dev/2018/go-io-testing/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const u = "http://192.168.0.16:4080"
const t = "5eoo8jg3aiwyDJN3Z9g8"

func main() {
	//setup()
	//fetchCredentials()
	//login()
	client := http.Client{}
	params := map[string]string{
		"name":      "hello",
		"path":      "hello",
		"parent_id": "2",
	}

	req := newGroup(u, params)
	req.Header.Add("PRIVATE-TOKEN", t)
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}

func migrate(pid, dst string) {

	//Create subgroups
	//split string

	/* 	params := map[string]string{
		"namespace": "sample",
		"path":      "upload-test",
		"overwrite": "false",
	} */
}
