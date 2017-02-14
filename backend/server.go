package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

// HTTP request handler
func handler(w http.ResponseWriter, r *http.Request) {
	// if the request is for root, load index page
	if (r.URL.Path == "/") || (r.URL.Path == "/index.html") {
		// run grunt task to grab all partials
		g, err := getResources()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// get template function based on index and execute to load page
		t, _ := template.ParseFiles("../index.html")
		t.Execute(w, g)

	} else if r.URL.Path == "/favicon.ico" { // load favicon.ico because who cares
		p, err := ioutil.ReadFile("../favicon.ico")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", p)
	} else { // if anyone tried to naviage to anything else
		fmt.Println("Attempted navigate to: " + r.URL.Path)
		return
	}
}

// start the servers
func main() {
	http.HandleFunc("/", handler)

	// start HTTPS server in goroutine
	go http.ListenAndServeTLS(":8081", "ssl/cert.pem", "ssl/key.pem", nil)

	// start HTTP server old-school and redirect to HTTPS
	http.ListenAndServe(":8080", nil)
	//http.HandlerFunc(httpsRedirect)
}
