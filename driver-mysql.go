package main

import (
	"net/url"
	"strings"

	"codeberg.org/momar/ternary"
)

func init() {
	drivers["mysql"] = Driver{
		Backup: func(uri *url.URL) (string, []byte, error) {
			password, _ := uri.User.Password()

			options := append([]string{
				"--host", ternary.Default(resolveIP(uri.Hostname()), "127.0.0.1").(string),
				"--port", ternary.Default(uri.Port(), "3306").(string),
				"--user", ternary.Default(uri.User.Username(), "root").(string),
				"--password", password,
				"--events", "--routines",
			},
				// database and tables from path, --all-databases otherwise
				ternary.Default(strings.Split(strings.Trim(uri.Path, " /"), "/"), []string{"--all-databases"}).([]string)...,
			)

			output, err := runInDocker([]string{}, []string{}, "mysqldump", options...).Output()
			return ".sql", output, err
		},
		Restore: func(*url.URL, []byte) error {
			// TODO:
			return nil
		},
	}
}
