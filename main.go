// Refactor code to make it easier to test: https://endler.dev/2018/go-io-testing/

package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	var args []string

	if len(os.Args) > 1 {
		args = os.Args[1:]
	} else {
		fmt.Println(`Usage of ./gitlab-migrate:
setup                  Initial setup
login                  Validate credentials
<pid> <destination>    Migrate <pid> to <destination>`)
		os.Exit(1)
	}

	switch args[0] {
	case "login":
		login()
		os.Exit(0)
	case "setup":
		setup()
		os.Exit(0)
	}

	if len(args) >= 3 {
		fmt.Println("Too many arguments")
	} else {
		migrate(args[0], args[1])
		os.Exit(0)
	}
}

func migrate(pid, dst string) {

	c := fetchCredentials()
	log.Println("Scheduling export...")
	_, err := scheduleExport(c.ExportURI, c.ExportToken, pid)
	if err != nil {
		log.Fatal(err)
	}

	var r *statusResp
	var filename string
	for {
		r, _, err = exportStatus(c.ExportURI, c.ExportToken, pid)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Export status:", r.ExportStatus)
		if r.ExportStatus == "finished" {

			t := time.Now()
			filename = "./" + t.Format("01-02-2006") + "-" + r.Path + ".tar.gz"
			log.Println("Downloading", filename)
			_, err := exportDownload(c.ExportURI, c.ExportToken, pid, filename)
			if err != nil {
				log.Fatal(err)
			}
			break
		}
		time.Sleep(10 * time.Second)
	}

	// create the groups
	log.Println("Creating groups")
	gid := newGroup(c.ImportURI, c.ImportToken, dst)
	// import the project
	log.Println("Importing project")
	importFile(c.ImportURI, c.ImportToken, gid, r.Path, filename)
	log.Println("Import complete")
}
