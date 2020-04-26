// Refactor code to make it easier to test: https://endler.dev/2018/go-io-testing/

package main

import (
	"fmt"
	"log"
	"time"
)

const u = "http://192.168.0.16:4080"
const t = "5eoo8jg3aiwyDJN3Z9g8"

func main() {
	//setup()
	c := fetchCredentials()
	//login()
	newGroup(c.ExportURI, c.ExportToken, "hello/test2")
	fmt.Println(getParentID(c.ExportURI, c.ExportToken, "test2", "hello/test2"))

}

func migrate(pid, dst string) {

	c := login()
	_, err := scheduleExport(c.ExportURI, c.ExportToken, pid)
	if err != nil {
		log.Fatal(err)
	}

	for retry := 0; retry < 5; retry++ {
		r, _, err := exportStatus(c.ExportURI, c.ExportToken, pid)
		if err != nil {
			log.Fatal(err)
		}

		if r.ExportStatus == "finished" {

			t := time.Now()
			filename := "./" + t.Format("01-02-2006") + "-" + r.Path + "tar.gz"
			exportDownload(c.ExportURI, c.ExportToken, pid, filename)
			break
		}

		time.Sleep(10 * time.Second)
	}
}
