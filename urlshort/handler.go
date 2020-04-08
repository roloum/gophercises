package urlshort

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v2"
)

const _yaml string = "yaml"
const _json string = "json"

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//if path is in the map, redirect to its URL
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		//otherwise serve fallback
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	routes, err := parse(yml, _yaml)
	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(routes), fallback), nil
}

//JSONHandler will parse the provided JSON and then return an http.HandlerFunc
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	routes, err := parse(json, _json)
	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(routes), fallback), nil
}

//parse parses structure into an array of route struct
//It can parse YAML or JSON
func parse(content []byte, format string) ([]route, error) {

	routes := []route{}
	var err error

	if format == _yaml {
		err = yaml.Unmarshal(content, &routes)
	} else if format == _json {
		err = json.Unmarshal(content, &routes)
	} else {
		err = errors.New("Unknown format")
	}

	return routes, err
}

//MongoDBHandler will load the routes from a MongoDB Collection
//Client must provide the MongoDB collection and the fallback handler
func MongoDBHandler(collection *mongo.Collection, fallback http.Handler) (
	http.HandlerFunc, error) {

	routes := []route{}
	var err error

	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var r route
		err := cur.Decode(&r)
		if err != nil {
			return nil, err
		}
		routes = append(routes, r)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return MapHandler(buildMap(routes), fallback), nil
}

//buildMap receives an array of route struct and creates a map
//that later on we pass to the MapHandler function
func buildMap(routes []route) map[string]string {

	m := make(map[string]string)

	for _, r := range routes {
		m[r.Path] = r.URL
	}

	return m
}

//route holds the URL a given path needs to be redirected to
type route struct {
	//Path
	Path string `yaml:"path" json:"path"`
	//URL ...
	URL string `yaml:"url" json:"url"`
}
