// Refactor code to make it easier to test: https://endler.dev/2018/go-io-testing/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	//login("", "")
	handle()
}

func printResponse(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
