package main

import (
	"fmt"
	//	"github.com/kabukky/httpscerts"
	"io/ioutil"
	//	"log"
	"net/http"
	"text/template"
)

//list of keys and paths for resource loading
var paths = map[string]string{
	"Style":         "../css/css.css",
	"Blog":          "../view/blog.html",
	"About":         "../view/about.html",
	"Connect":       "../view/connect.html",
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

//request handler for index
func indexHandler(wrt http.ResponseWriter, req *http.Request) {
	fmt.Println("website index")
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

//request handler for favicon because who cares
func faviconHandler(wrt http.ResponseWriter, req *http.Request) {
	fmt.Println("favicon")
	f, err := ioutil.ReadFile("../favicon.ico")
	if err != nil {
		http.Error(wrt, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(wrt, "%s", f)
}

//request handler for images
func imageHandler(wrt http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.Path[12:])
	f, err := ioutil.ReadFile("../" + req.URL.Path)
	if err != nil {
		http.Error(wrt, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(wrt, "%s", f)
}

//redirect HTTP to HTTPS (works, but unused)
//func redirectToHttps(wrt http.ResponseWriter, req *http.Request) {
//http.Redirect(wrt, req, "https://jesserominske-com-bonkleton.c9users.io:8081"+req.RequestURI, http.StatusMovedPermanently)
//}

//start the server
func main() {
	//TODO: actually get TLS working, you idiot

	//check if the cert files are available
	//err := httpscerts.Check("../ssl/cert.pem", "../ssl/key.pem")
	//if they are not available, generate new ones
	//if err != nil {
	//	err = httpscerts.Generate("../ssl/cert.pem", "../ssl/key.pem", "https://jesserominske-com-bonkleton.c9users.io:8081")
	//	if err != nil {
	//		log.Fatal("Error: Couldn't create https certs.")
	//	}
	//}

	//create a new ServeMux for HTTP connections
	httpMux := http.NewServeMux()
	httpMux.Handle("/", http.HandlerFunc(indexHandler))
	httpMux.Handle("/favicon.ico", http.HandlerFunc(faviconHandler))
	httpMux.Handle("/css/images/", http.HandlerFunc(imageHandler))

	//create a new ServeMux for HTTPS connections
	httpsMux := http.NewServeMux()
	httpsMux.Handle("/", http.HandlerFunc(indexHandler))
	httpsMux.Handle("/favicon.ico", http.HandlerFunc(faviconHandler))
	httpsMux.Handle("/css/images/", http.HandlerFunc(imageHandler))

	//start HTTPS server in goroutine
	go http.ListenAndServeTLS(":8081", "../ssl/cert.pem", "../ssl/key.pem", httpsMux)

	//start HTTP server old-school and redirect to HTTPS
	http.ListenAndServe(":8080", httpMux)
	//http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))
}
