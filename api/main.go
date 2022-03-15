package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"golang.org/x/text/encoding/charmap"
)

func GetMainPage(rw http.ResponseWriter, r *http.Request) {
	path := filepath.Join("../html", "transcriber.html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	err = tmpl.Execute(rw, nil)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}

func Transcribe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Transcribe audio endpoint hit")
	file, handler, err := r.FormFile("audio")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	dir, _ := os.Getwd()

	os.Chdir("../")

	// read all of the contents of uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(err)
	}

	// write this byte array to our temporary file
	err = ioutil.WriteFile("temp-audio/temp.mp3", fileBytes, 0644)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cmd := exec.Command("py", "transcriber/transcriber.py", "temp-audio/temp.mp3")

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))

	f, e := os.Open("output.txt")

	if e != nil {
		fmt.Println(err.Error())
		return
	}

	defer f.Close()

	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(f)
	result, err := ioutil.ReadAll(reader)

	if e != nil {
		fmt.Println(err.Error())
		return
	}

	// return that we have successfully uploaded our file!
	fmt.Fprint(w, string(result))

	os.Chdir(dir)
}

func setupRoutes() {
	http.HandleFunc("/transcribe", Transcribe)
	http.HandleFunc("/", GetMainPage)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
}
