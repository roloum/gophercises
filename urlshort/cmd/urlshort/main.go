package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/roloum/gophercises/urlshort"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	var yamlFile, jsonFile string
	var mongodb bool

	flag.StringVar(&yamlFile, "yaml-file", "", "YAML File")
	flag.StringVar(&jsonFile, "json-file", "", "JSON File")
	flag.BoolVar(&mongodb, "mongodb", false, "MongoDB collection")
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

	if mongodb {
		fmt.Println("Loading routes from MongoDB collection")

		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			panic(err)
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		collection := client.Database("urlshort").Collection("route")

		mongoHandler, err := urlshort.MongoDBHandler(collection, handler)
		if err != nil {
			panic(err)
		}

		handler = mongoHandler

		err = client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

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
