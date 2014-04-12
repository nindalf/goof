package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	ip       = flag.String("i", "", "The IP Address the server should run on")
	port     = flag.Int("p", 8086, "The port on which the server listens")
	root     = flag.String("f", "", "The name of the file/folder to be shared")
	count    = flag.Int("c", 1, "The number of times the file/folder should be shared")
	duration = flag.Int("t", 0, "Server timeout")
	archive  = flag.Bool("a", false, "Whether the folder should be compressed before serving")
	upload   = flag.Bool("u", false, "Serve a form that enables users to upload files.")
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

func serveFile(handler *fileHandler, endpoint string) {
	http.Handle("/", handler)
	log.Println("Serving", handler.root, "at", externalEndpoints(endpoint))
	log.Fatal(http.ListenAndServe(endpoint, nil))
}

func serveFolderArchive(root string, count, duration int, endpoint string) {
	tarfile, err := archiveDir(root)
	if err != nil {
		log.Fatal(err)
	}
	go exitAfter(duration)
	serveFile(&fileHandler{tarfile, count}, endpoint)
}

func serveFolderInteractive(root string, duration int, endpoint string) {
	log.Println("Serving", root, "at", externalEndpoints(endpoint))
	exitAfter(duration)
	log.Fatal(http.ListenAndServe(endpoint, http.FileServer(http.Dir(root))))
}

func exitAfter(minutes int) {
	if minutes == 0 {
		return
	}
	delay := fmt.Sprintf("%dm", minutes)
	duration, _ := time.ParseDuration(delay)
	log.Println("Will exit automatically after", duration)
	<-time.After(duration)
	log.Fatal("Server timed out.")
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

func externalEndpoints(endpoint string) []string {
	var ips []string
	if strings.Index(endpoint, ":") != 0 {
		return append(ips, "http://"+endpoint)
	}
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if strings.Index(addr.String(), ":") != -1 {
			//ipv6 address
			continue
		}
		ips = append(ips, "http://"+strings.Split(addr.String(), "/")[0]+endpoint)
	}
	return ips
}

func main() {
	flag.Parse()
	endpoint := fmt.Sprintf("%s:%d", *ip, *port)
	fi, err := os.Stat(*root)
	if err != nil {
		log.Fatal("Path is invalid")
	}
	if fi.IsDir() == true {
		if *archive == false {
			serveFolderInteractive(*root, *duration, endpoint)
		} else {
			serveFolderArchive(*root, *count, *duration, endpoint)
		}
	} else {
		// is a file
		go exitAfter(*duration)
		serveFile(&fileHandler{*root, *count}, endpoint)
	}
}
