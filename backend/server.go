package main

import (
	"database/sql"
	"fmt"
	"github.com/kabukky/httpscerts"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"text/template"
)

//global sql.DB to access the database by all handlers
var db *sql.DB
var dbErr error
//list of keys and paths for resource loading
var paths = map[string]string{
	"Style":         "../css/css.css",
	"Blog":          "../view/blog.html",
	"AddBlogModal":  "../view/addBlog.html",
	"About":         "../view/about.html",
	"Connect":       "../view/connect.html",
	"ContentScript": "../js/content.js",
	"ModelScript":   "../js/model.js",
  "BlogScript":    "../js/blog.js"}

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
func indexHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("website index")
	//grab all partials
	partials, err := loadPartials()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	//get template function based on index and execute to load page
	t, _ := template.ParseFiles("../index.html")
	t.Execute(res, partials)
}

//request handler for favicon because who cares
func faviconHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("favicon")
	f, err := ioutil.ReadFile("../favicon.ico")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(res, "%s", f)
}

//request handler for images
func imageHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.Path[12:])
	f, err := ioutil.ReadFile("../" + req.URL.Path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(res, "%s", f)
}

//redirect HTTP to HTTPS (works, but unused)
func redirectToHttps(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "localhost:8081"+req.RequestURI, http.StatusMovedPermanently)
}

//basic panic response to error, saving lines
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

//start the server
func main() {
	//TODO: actually get TLS working, dork

	//check if the cert files are available
	err := httpscerts.Check("../ssl/cert.pem", "../ssl/key.pem")
	//if they are not available, generate new ones
	if err != nil {
		err = httpscerts.Generate("../ssl/cert.pem", "../ssl/key.pem", "localhost:8081")
		checkError(err)
	}

	//create sql.DB and check for errors
  db, dbErr = sql.Open("mysql", dbOptions["dbAdmin"] + ":" + dbOptions["dbPWord"] + "@/" + dbOptions["dbName"])
	checkError(dbErr)
  defer db.Close()
  //test the connection to the database
  dbErr = db.Ping()
  checkError(dbErr)

	//create a new ServeMux for HTTP connections (delete once TLS works)
	httpMux := http.NewServeMux()
	httpMux.Handle("/", http.HandlerFunc(indexHandler))
	httpMux.Handle("/favicon.ico", http.HandlerFunc(faviconHandler))
	httpMux.Handle("/css/images/", http.HandlerFunc(imageHandler))
	httpMux.Handle("/blog/post/", http.HandlerFunc(blogHandler)) //in blog.go

	//create a new ServeMux for HTTPS connections
	httpsMux := http.NewServeMux()
	httpsMux.Handle("/", http.HandlerFunc(indexHandler))
	httpsMux.Handle("/favicon.ico", http.HandlerFunc(faviconHandler))
	httpsMux.Handle("/css/images/", http.HandlerFunc(imageHandler))
	httpsMux.Handle("/blog/post/", http.HandlerFunc(blogHandler)) //in blog.go

	//start HTTPS server in goroutine
	go http.ListenAndServeTLS(":8081", "../ssl/cert.pem", "../ssl/key.pem", httpsMux)

	//start HTTP server old-school and redirect to HTTPS
	http.ListenAndServe(":8080", httpMux) //replace with line below once TLS works
	//http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))
}
