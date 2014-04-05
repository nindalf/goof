package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	ip       = flag.String("i", "127.0.0.1", "The IP Address the server should run on")
	port     = flag.Int("p", 8086, "The port on which the server listens")
	root     = flag.String("f", "", "The name of the file/folder to be shared")
	count    = flag.Int("c", 1, "The number of times the file/folder should be shared")
	duration = flag.Int("t", 0, "Server timeout")
)

type fileHandler struct {
	root  string
	count int
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving " + path.Base(f.root) + " to " + strings.Split(r.RemoteAddr, ":")[0])
	w.Header().Set("Content-Disposition", "attachment;filename=\""+path.Base(f.root)+"\"")
	http.ServeFile(w, r, f.root)
	f.count = f.count - 1
	if f.count == 0 {
		log.Fatal("Finished serving. Server exiting.")
	}
}

func exitafter(minutes int) {
	if minutes == 0 {
		return
	}
	delay := fmt.Sprintf("%dm", minutes)
	duration, _ := time.ParseDuration(delay)
	log.Println("Will exit automatically after", duration)
	<-time.After(duration)
	log.Fatal("Server timed out.")
}

func serveFile(handler http.Handler, endpoint string) {
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(endpoint, nil))
}

func serveFolder(root string, count, duration int, endpoint string) {
	tarfile, err := archiveDir(root)
	if err != nil {
		log.Fatal(err)
	}
	go exitafter(duration)
	log.Println("Serving", tarfile, "at", endpoint)
	serveFile(&fileHandler{tarfile, count}, endpoint)
}

func newArchWriter(dirname string) (*tar.Writer, error) {
	w, err := os.Create(dirname + ".tar")
	if err != nil {
		return new(tar.Writer), err
	}
	cw := gzip.NewWriter(w)
	return tar.NewWriter(cw), nil
}

func archiveDir(root string) (string, error) {
	log.Println("Creating archive of", root)
	dir := filepath.Dir(root)
	tw, err := newArchWriter(root)
	if err != nil {
		return "", err
	}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		header, _ := tar.FileInfoHeader(info, "")
		header.Name = path[len(dir):]
		tw.WriteHeader(header)
		if info.IsDir() == false {
			data, _ := ioutil.ReadFile(path)
			tw.Write(data)
			tw.Flush()
		}
		return nil
	})
	tw.Close()
	log.Println("Created", root+".tar")
	return root + ".tar", nil
}

func main() {
	flag.Parse()
	endpoint := fmt.Sprintf("%s:%d", *ip, *port)
	fi, err := os.Stat(*root)
	if err != nil {
		log.Fatal("Path is invalid")
	}
	if fi.IsDir() == true {
		serveFolder(*root, *count, *duration, endpoint)
	} else {
		// is a file
		go exitafter(*duration)
		log.Println("Serving", *root, "at", endpoint)
		serveFile(&fileHandler{*root, *count}, endpoint)
	}
}
