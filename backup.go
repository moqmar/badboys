package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

func backup(name, rawuri, dir string) error {
	var extension string
	var content []byte
	uri, err := url.Parse(rawuri)
	if err != nil {
		return fmt.Errorf("[%s] - url error: %s", name, err)
	}

	switch uri.Scheme {
	case "sqlite3":
		extension, content, err = sqlite3(uri)
	case "mysql":
		extension, content, err = mysql(uri)
	case "postgres":
		extension, content, err = postgres(uri)
	case "mongodb":
		extension, content, err = mongodb(uri)
	default:
		return fmt.Errorf("[%s] - no handler for %s databases", name, uri.Scheme)
	}
	if err != nil {
		return fmt.Errorf("[%s] - database error: %s", name, err)
	}

	if !dry {
		err = os.MkdirAll(path.Join(dir, name), 0700|os.ModeDir)
		if err != nil {
			return fmt.Errorf("[%s] - mkdir error: %s", name, err)
		}

		err = ioutil.WriteFile(path.Join(dir, name, filename()+extension), content, 0700)
		if err != nil {
			return fmt.Errorf("[%s] - write error: %s", name, err)
		}
	}

	return nil
}

func filename() string {
	return now.Format(cfg.Get("filename").Default("2006-01-02.15-04").String())
}
