package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type GoosemonkeyConfig struct {
	Dir  string
	Port int
}

type FilesResponse struct {
	Files []string
}

func getPortAndDirFromConfig() (string, string) {
	config := GoosemonkeyConfig{}
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("failed to open config.json: ", err)
	} else {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&config)
		if err != nil {
			log.Fatal("failed to decode config.json: ", err)
		}
	}
	return strconv.Itoa(config.Port), config.Dir
}

func getFilesResponseFromReadDirResponse(dirFiles []os.FileInfo) *FilesResponse {
	var files []string
	for _, f := range dirFiles {
		files = append(files, f.Name())
	}
	filesResponse := &FilesResponse{Files: files}
	return filesResponse
}

func rootRouteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
}

func getFilesRouteHandler(dirname string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dirFiles, err := ioutil.ReadDir(dirname)
		if err != nil {
			fmt.Fprintf(w, "Failed to read directory")
		} else {
			filesResponse := getFilesResponseFromReadDirResponse(dirFiles)
			response, err := json.Marshal(filesResponse)
			if err != nil {
				fmt.Fprintf(w, "Failed to marshal JSON")
			} else {
				fmt.Fprintf(w, string(response))
			}
		}
	}
}

func main() {
	port, dirname := getPortAndDirFromConfig()

	http.HandleFunc("/", rootRouteHandler)

	http.HandleFunc("/files", getFilesRouteHandler(dirname))

	log.Println("Hello World! We will be serving out of %s today.", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
