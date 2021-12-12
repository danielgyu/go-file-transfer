package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"text/template"
)

var FileDirectory = "./files/"

func index(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("upload.html")
	if err != nil {
		log.Println("Error loading static html")
	}

	template.Execute(w, nil)
}

func copyToFile(dst *os.File, src multipart.File, wg *sync.WaitGroup) {
	io.Copy(dst, src)
	wg.Done()
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

	var wg = new(sync.WaitGroup)
	wg.Add(1)
	go copyToFile(f, file, wg)
	wg.Wait()
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadFile)

	log.Println("Serving file server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
