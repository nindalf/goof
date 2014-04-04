package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	filename = flag.String("f", "", "The name of the file to be shared")
	n        = flag.Int("n", 1, "The number of times the file should be shared")
	t        = flag.Int("t", 0, "Server timeout")
)

type fileHandler struct {
	filename string
	n        int
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.n = f.n - 1
	if f.n == -1 {
		log.Fatal("Finished serving. Server exiting.")
	}
	log.Println(f.filename)
	http.ServeFile(w, r, f.filename)
}

func exitafter(minutes int) {
	delay := fmt.Sprintf("%dm", minutes)
	duration, _ := time.ParseDuration(delay)
	<-time.After(duration)
	log.Fatal("Server timed out.")
}

func main() {
	flag.Parse()
	go exitafter(*t)
	handler := fileHandler{"/home/nindalf/Pictures/wallpapers/octocats/chellocat.jpg", *n}
	http.ListenAndServe(":8086", &handler)
}
