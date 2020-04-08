package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/roloum/gophercises/urlshort"
)

func main() {

	var yamlFile, jsonFile string
	flag.StringVar(&yamlFile, "yaml-file", "", "YAML File")
	flag.StringVar(&jsonFile, "json-file", "", "JSON File")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var yaml []byte
	var err error
	if yamlFile == "" {

		// Build the YAMLHandler using the mapHandler as the
		// fallback
		data := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
		//convert data to byte array
		yaml = []byte(data)
	} else {
		fmt.Printf("Reading YAML from file: %v\n", yamlFile)

		//read file
		yaml, err = ioutil.ReadFile(yamlFile)
		if err != nil {
			panic(err)
		}
	}

	//changed to avoid converting data to string and then back to byte array
	//when yaml file is provided
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	//JSON handler. Since YAML is superset of JSON, we're going to read
	//json from a file, use the yamlHandler as fallback and create the
	//jsonHandler using the same YAMLHandler function
	var handler http.HandlerFunc = yamlHandler

	if jsonFile != "" {
		fmt.Printf("Reading JSON from file: %v\n", jsonFile)

		json, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			panic(err)
		}
		jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
		if err != nil {
			panic(err)
		}
		handler = jsonHandler

	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}
