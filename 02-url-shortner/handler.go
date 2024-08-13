package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return a http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths (keys in the map) to their corresponding URL
// (values that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if we can match a path redirect to it
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		// else do the fallback
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return a http.HandlerFunc
// (which also implements http.Handler) that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAMl, then the fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having invalid YAML data.
//
// See MapHandler to create similar http.HandlerFunc via a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// 1. Parse the YAML
	// 2. Convert YAML array into map
	// 3. return a map handler using the map
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pathUrl := range pathUrls {
		pathToUrls[pathUrl.Path] = pathUrl.URL
	}
	return pathToUrls
}

func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
