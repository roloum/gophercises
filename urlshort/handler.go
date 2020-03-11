package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

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

	routes, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(routes), fallback), nil
}

//parseYAML converts the YAML structure into an array of route struct
func parseYAML(yml []byte) ([]route, error) {
	routes := []route{}

	err := yaml.Unmarshal(yml, &routes)

	return routes, err
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
	Path string `yaml:"path"`
	//URL ...
	URL string `yaml:"url"`
}
