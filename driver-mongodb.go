package main

import (
	"net/url"
	"strings"

	"codeberg.org/momar/ternary"
)

func init() {
	drivers["mongodb"] = Driver{
		Backup: func(uri *url.URL) (string, []byte, error) {
			ip := resolveIP(uri.Hostname())
			if strings.Contains(ip, ":") {
				uri.Host = "[" + ip + "]" + ternary.If(uri.Port() != "", ":"+uri.Port(), "").(string)
			} else {
				uri.Host = ip + ternary.If(uri.Port() != "", ":"+uri.Port(), "").(string)
			}

			options := []string{
				"--uri", uri.String(),
				"--archive=-",
			}

			output, err := runInDocker([]string{}, []string{}, "mongodump", options...).Output()
			return ".mongoarchive", output, err
		},
		Restore: func(*url.URL, []byte) error {
			// TODO:
			return nil
		},
	}
}
