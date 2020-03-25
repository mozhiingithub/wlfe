package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	http.Handle("/", http.FileServer(http.Dir(p)))
	http.ListenAndServe(":8000", nil)
}
