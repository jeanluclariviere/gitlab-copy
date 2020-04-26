// Refactor code to make it easier to test: https://endler.dev/2018/go-io-testing/

package main

import "fmt"

const u = "http://192.168.0.16:4080"
const t = "5eoo8jg3aiwyDJN3Z9g8"

func main() {
	//setup()
	c := fetchCredentials()
	//login()

	fmt.Println("GroupID:", newGroup(c.ExportURI, c.ExportToken, "bird/bird"))

	//fmt.Println(getParentID(c.ExportURI, c.ExportToken, "test", "hello/test"))

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
