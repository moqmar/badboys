package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/moqmar/gonfig"
)

func backupDir(dir string) {
	running++

	fmt.Printf("Backing up databases from: %s\n", dir)
	repo := gonfig.Open(path.Join(dir, "databases.yaml"))

	dirRunning := 0
	dirUpdate := make(chan bool)
	for name, uri := range repo.StringMap() {
		dirRunning++
		go backupAndPrune(dirUpdate, name, uri, dir)
	}

	for <-dirUpdate {
		dirRunning--
		if dirRunning <= 0 {
			break
		}
	}

	running--
	update <- running
}

func backupAndPrune(dirUpdate chan bool, name string, uri string, dir string) {
	err := backup(name, uri, dir)
	if err != nil {
		fmt.Printf("%s\n", err)
		result = 1

		dirUpdate <- true
		return
	}

	err = prune(name, dir)
	if err != nil {
		fmt.Printf("%s\n", err)
		result = 1
	}

	dirUpdate <- true
}

func backup(name, rawuri, dir string) error {
	var extension string
	var content []byte
	uri, err := url.Parse(rawuri)
	if err != nil {
		return fmt.Errorf("[%s] - url error: %s", name, err)
	}

	driver, ok := drivers[uri.Scheme]
	if !ok {
		return fmt.Errorf("[%s] - no handler for %s databases", name, uri.Scheme)
	}

	extension, content, err = driver.Backup(uri)
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
