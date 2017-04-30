package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

//list of keys and paths for resource loading
var paths = map[string]string{
	"Style":         "../css/css.css",
	"Blog":          "../view/blog.html",
	"About":         "../view/about.html",
	"Connect":       "../view/connect.html",
	"Wiki":          "../view/wiki.html",
	"ContentScript": "../js/content.js"}

//loads all relevant partials
func loadPartials() (map[string]string, error) {
	g := make(map[string]string)

	//load resources from paths
	for key, path := range paths {
		body, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		g[key] = string(body)
	}
	return g, nil
}

//HTTP request handler
func handler(wrt http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.Path)
	if req.URL.Path == "/" { //load the index
		//grab all partials
		partials, err := loadPartials()
		if err != nil {
			http.Error(wrt, err.Error(), http.StatusInternalServerError)
			return
		}
		//get template function based on index and execute to load page
		t, _ := template.ParseFiles("../index.html")
		t.Execute(wrt, partials)
	} else if req.URL.Path == "/favicon.ico" { //load favicon.ico because who cares
		f, err := ioutil.ReadFile("../favicon.ico")
		if err != nil {
			http.Error(wrt, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(wrt, "%s", f)
	} else if req.URL.Path[:11] == "/css/images" { //if we want to load an image
		f, err := ioutil.ReadFile("../" + req.URL.Path)
		if err != nil {
			http.Error(wrt, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(wrt, "%s", f)
	} else { //in all other cases load the index
		//grab all partials
		partials, err := loadPartials()
		if err != nil {
			http.Error(wrt, err.Error(), http.StatusInternalServerError)
			return
		}
		//get template function based on index and execute to load page
		t, _ := template.ParseFiles("../index.html")
		t.Execute(wrt, partials)
	}
}

//start the servers
func main() {
	http.HandleFunc("/", handler)

	//start HTTPS server in goroutine
	go http.ListenAndServeTLS(":8081", "ssl/cert.pem", "ssl/key.pem", nil)

	//start HTTP server old-school and redirect to HTTPS
	http.ListenAndServe(":8080", nil)
	//http.HandlerFunc(httpsRedirect)
}
