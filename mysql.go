package main

import (
	"net/url"

	"codeberg.org/momar/ternary"
)

func mysql(uri *url.URL) (string, []byte, error) {
	password, _ := uri.User.Password()

	options := []string{
		"--host", ternary.Default(resolveIP(uri.Hostname()), "127.0.0.1").(string),
		"--port", ternary.Default(uri.Port(), "3306").(string),
		"--user", ternary.Default(uri.User.Username(), "root").(string),
		"--password", password,
	}

	databaseOptions := []string{"--events", "--routines", "--triggers"}
	// TODO: populate databaseOptions

	output, err := runInDocker([]string{}, []string{}, "mysqldump", append(options, databaseOptions...)...).Output()
	return ".sql", output, err
}
