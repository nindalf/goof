package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var (
	ip       = flag.String("i", "127.0.0.1", "The IP Address the server should run on")
	port     = flag.Int("p", 8086, "The port on which the server listens")
	filepath = flag.String("f", "", "The name of the file to be shared")
	count    = flag.Int("c", 1, "The number of times the file should be shared")
	duration = flag.Int("t", 0, "Server timeout")
)

type fileHandler struct {
	filepath string
	count    int
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.count = f.count - 1
	if f.count == -1 {
		log.Fatal("Finished serving. Server exiting.")
	}
	log.Println(f.filepath)
	w.Header().Set("Content-Disposition", "attachment;filename=\""+path.Base(f.filepath)+"\"")
	http.ServeFile(w, r, f.filepath)
}

func exitafter(minutes int) {
	if minutes == 0 {
		return
	}
	delay := fmt.Sprintf("%dm", minutes)
	duration, _ := time.ParseDuration(delay)
	<-time.After(duration)
	log.Fatal("Server timed out.")
}

func checkFile(filepath string) {
	if fi, err := os.Stat(filepath); err != nil || fi.IsDir() == true {
		log.Fatal("File does not exist")
	}
}

func main() {
	flag.Parse()
	go exitafter(*duration)
	checkFile(*filepath)
	handler := fileHandler{*filepath, *count}
	endpoint := fmt.Sprintf("%s:%d", *ip, *port)
	http.Handle("/", &handler)
	log.Fatal(http.ListenAndServe(endpoint, nil))
}
