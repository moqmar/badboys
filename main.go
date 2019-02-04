package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/moqmar/gonfig"
)

var cfg = gonfig.Open("/etc/badboys.yaml", `---
# List of folders with "databases.yaml" files to back up; can also contain wildcards
repositories:
- /var/backup/badboys

# See https://golang.org/pkg/time/#pkg-constants for an explanation of the format
filename: 2006-01-02.15-04

# How many backups to keep, if one condition matches (true: keep all, false: ignore)
retention:
  latest: 3      # always keep the latest backup
  hourly: false  # don't keep hourly backups
  daily: 7       # keep daily backups for 7 days
  weekly: 4      # keep weekly backups for 4 weeks
  monthly: 12    # keep monthly backups for 12 weeks
  yearly: true   # keep one yearly backup forever

# Shell command to run on completion; gets the exit code of badboys as $1.
oncomplete: "echo '\\o/'"
`)

var now = time.Now()
var dry = false

func main() {
	result := 0

	// TODO: command line arguments - "--dry", "--restore" and a manual list of repositories (+ WARNING if the repository is not matched by the config file)
	for _, glob := range cfg.Get("repositories").StringList() {
		globResult, err := filepath.Glob(glob)
		if err != nil {
			fmt.Printf("Glob error: %s\n", err)
			result = 1
		}

		for _, dir := range globResult {
			fmt.Printf("Backing up databases from: %s\n", dir)
			repo := gonfig.Open(path.Join(dir, "databases.yaml"))

			for name, uri := range repo.StringMap() {
				err = backup(name, uri, dir)
				if err != nil {
					fmt.Printf("%s\n", err)
					result = 1
				}

				err = prune(name, dir)
				if err != nil {
					fmt.Printf("%s\n", err)
					result = 1
				}
			}
		}
	}

	err := exec.Command("/bin/sh", "-c", cfg.Get("oncomplete").String(), "--", string(result)).Run()
	if err != nil {
		fmt.Printf("Oncomplete error: %s\n", err)
		result = 1
	}

	os.Exit(result)
}
