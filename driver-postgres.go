package main

import (
	"net/url"
	"strings"

	"codeberg.org/momar/ternary"
)

func init() {
	drivers["postgres"] = Driver{
		Backup: func(uri *url.URL) (string, []byte, error) {
			password, _ := uri.User.Password()

			options := append([]string{
				"--format", "custom",
				"--oids",
			},
				// database and tables from path
				ternary.Default(strings.Split(strings.Replace(strings.Trim(uri.Path, " /"), "/", "/-d/", -1), "/"), []string{}).([]string)...,
			)

			env := []string{
				"PGHOST=" + ternary.Default(resolveIP(uri.Hostname()), "127.0.0.1").(string),
				"PGPORT=" + ternary.Default(uri.Port(), "5432").(string),
				"PGUSER=" + ternary.Default(uri.User.Username(), "postgres").(string),
				"PGPASSWORD=" + password,
			}

			output, err := runInDocker([]string{}, env, "pg_dump", options...).Output()
			return ".pgdump", output, err
		},
		Restore: func(*url.URL, []byte) error {
			// TODO:
			return nil
		},
	}
}
