package main

import (
	"net/http"
)

type fileHandler struct {
	filename string
	times    int
}

func (f fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, f.filename)
}

func main() {
	handler := fileHandler{"/home/nindalf/Pictures/wallpapers/octocats/baracktocat.jpg", 1}
	http.ListenAndServe(":8086", handler)
}
