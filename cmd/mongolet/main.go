package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"os"
	"path/filepath"

	"github.com/mgenov/mongodb-provisioning/tarutil"
)

var (
	port    = flag.Int("port", 8080, "the port on which app will listen")
	workDir = flag.String("workDir", "workdir", "the work folder which will be used")
)

func main() {
	flag.Parse()

	initDir(*workDir)

	http.HandleFunc("/initMongo", initMongoDb)

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	var err error
	for try := 0; try < 10; try++ {
		_, err := http.Get("localhost:8080")

		if err == nil {
			break
		}
	}

	if err == nil {
		log.Printf("mongolet was started on port: %d\n", *port)
	} else {
		log.Printf("got error '%v' while starting server, so give up", err)
	}

	select {}
}

func initMongoDb(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("expected method POST, but got %s", r.Method), http.StatusBadRequest)
		return
	}

	err := initDir(filepath.Join(*workDir, "mongo1"))

	if err != nil {
		http.Error(w, "got error while initialzing work dir.", http.StatusInternalServerError)
		return
	}

	rawUrl := r.FormValue("url")

	_, err = url.Parse(rawUrl)
	if err != nil {
		log.Printf("got error %v while parsing url.\n", err)
		http.Error(w, "expected valid url param", http.StatusBadRequest)
		return
	}

	err = downloadMongo(*workDir, rawUrl)

	if err != nil {
		log.Printf("got error %v while downloading mongo.\n", err)

		http.Error(w, "got error while downloading mongo", http.StatusInternalServerError)
		return
	}

	f, err := os.Open(filepath.Join(*workDir, "mongodb.tgz"))

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	err = tarutil.Untar(f, *workDir)

	if err != nil {
		http.Error(w, "got error while extracting mongo version", http.StatusInternalServerError)
	}
}

func downloadMongo(path, url string) error {
	mongoPath := filepath.Join(path, "mongodb-current.tgz")

	_, err := os.Open(mongoPath)

	// file is already available so we can skip the same download
	if err == nil {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	tfile, err := os.Create(mongoPath)
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

// initDir initializes directory with default structure
func initDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0777)

		if err != nil {
			return err
		}
	} else if err != nil { // maybe file is same as name of a file?
		return err
	}
	return nil
}
