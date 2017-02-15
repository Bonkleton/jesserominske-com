package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

// HTTP request handler
func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" { // load favicon.ico because who cares
		p, err := ioutil.ReadFile("../favicon.ico")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", p)
	} else { // in all other cases load the index

		// run grunt task to grab all partials
		g, err := runGrunt()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get template function based on index and execute to load page
		t, _ := template.ParseFiles("../index.html")
		t.Execute(w, g)
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
