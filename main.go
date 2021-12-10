package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

var FileDirectory = "./files/"

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func copyToFile(dst *os.File, src multipart.File) {
	io.Copy(dst, src)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 50)

	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		log.Printf("Error retrieving file: %s", err)
		return
	}
	defer file.Close()

	path := fmt.Sprintf("%s%s", FileDirectory, handler.Filename)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return
	}
	defer f.Close()

	go copyToFile(f, file)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadFile)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
