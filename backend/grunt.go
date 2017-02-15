package main

import (
	"io/ioutil"
)

// loads a partial
func loadResource(file string) ([]byte, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// runs the grunt-task
func runGrunt() (map[string]string, error) {

	// list of keys and paths
	paths := map[string]string{
		"Style":         "../css/css.css",
		"Home":          "../view/home.html",
		"Blog":          "../view/blog.html",
		"About":         "../view/about.html",
		"Connect":       "../view/connect.html",
		"ContentScript": "../js/content.js"}

	// declare new map to be returned
	g := make(map[string]string)

	// replace the paths with the loaded resources
	for k, v := range paths {
		// load the resource
		r, e := loadResource(v)
		if e != nil {
			return nil, e
		}

		// put resource into returned map
		g[k] = string(r)
	}

	return g, nil
}
