package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// initMongoDir initializes mongodirectory with default structure
func initMongoDir(path, name string) error {
	st, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(path, 0777)

		if err != nil {
			return err
		}
	}

	if err == nil && st.IsDir() {
		return nil
	}

	initPath := filepath.Join(path, name)

	err = os.Mkdir(initPath, 0777)

	if err != nil {
		return err
	}

	return nil
}

func downloadMongo(path string) error {
	url := "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-3.0.1.tgz?_ga=1.52666216.179553068.1401809707"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	tname := filepath.Join(path, "mongodb.tgz")

	tfile, err := os.Create(tname)
	if err != nil {
		return err
	}

	io.Copy(tfile, resp.Body)

	err = tfile.Close()
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil

}

//download and extract mongodb images from internet (admin could serve mongodb images by providing metadata as json)
//{version: "3.0.1", "os": "linux", "arch":"x86_64", path: "/dw/mongodb-3.0.1.tgz"}
//use mongod for running of version
//read configuration of existing db from existing configuration files and merge results with the configuration that need to be available on that machine
func main() {
	workDir := "/Users/mgenov/temp/mdd/"

	err := initMongoDir(workDir, "mongo1")
	if err != nil {
		log.Fatal(err)
	}

	err = downloadMongo(workDir)

	if err != nil {
		log.Fatal(err)
	}
}
