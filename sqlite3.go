package main

import (
	"net/url"
	"strings"
)

func sqlite3(uri *url.URL) (string, []byte, error) {
	path := strings.SplitN(uri.String(), "://", 2)[1]
	output, err := runInDocker([]string{path + ":/database.sqlite3"}, []string{},
		"sqlite3", "/database.sqlite3", ".dump").Output()
	return ".sql", output, err
}
