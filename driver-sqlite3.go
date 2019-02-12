package main

import (
	"net/url"
	"strings"
)

func init() {
	drivers["sqlite3"] = Driver{
		Backup: func(uri *url.URL) (string, []byte, error) {
			path := strings.SplitN(uri.String(), "://", 2)[1]
			output, err := runInDocker([]string{path + ":/database.sqlite3"}, []string{},
				"sqlite3", "/database.sqlite3", ".dump").Output()
			return ".sql", output, err
		},
		Restore: func(*url.URL, []byte) error {
			// TODO:
			return nil
		},
	}
}
