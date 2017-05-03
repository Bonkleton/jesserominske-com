package main

import (
	"database/sql"
	"fmt"
	"github.com/kabukky/httpscerts"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
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

//signature validation for blog form
func signature(u string, p string) (bool, error) {
  var dbUName, dbPWord string
	//validate username
	result := db.QueryRow("SELECT username, password FROM users WHERE username=?", u)
  err := result.Scan(&dbUName, &dbPWord)
  if err != nil {
    fmt.Println("Invalid username or password!")
    return false, err
  }
  //validate password
  err = bcrypt.CompareHashAndPassword([]byte(dbPWord), []byte(p))
  if err != nil {
		fmt.Println("Hash error!")
    return false, err
  }
  // If the login succeeded
	return true, nil
}

//finds rune in string
func runeInString(s string, r rune) (int, bool) {
	for i,c := range s {
  	if c == r {
			return i, true
		}
	}
	fmt.Println("Rune " + string(r) + " not found in string " + s)
	return 0, false
}

//blog body formatter (tail-recursive)
func paragraphize(body string, result string) (string, string) {
	index, isInBody := runeInString(body, '\n')
	//if there are no newlines, just wrap the body
	if !isInBody {
		return "", "<p class='lead'>\n" + body + "\n</p>\n"
	} else {
		wrapped := "<p class='lead'>\n" + body[:index] + "\n</p>\n"
		//if the newline is the end of the body, we're done
		if len(body[index:]) <= 1 {
			return "", result + wrapped
		}
		//chop off current paragraph, add wrapped paragraph to result, and call again
		return paragraphize(body[index + 1:], result + wrapped)
	}
}

//request handler for blog
func blogHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
    http.Redirect(res, req, "localhost:8080", 301) //change this once appropriate
    return
  }
	//get the username/password from the submitted post form
  uName := req.FormValue("uname")
  pWord := req.FormValue("pword")
	//validate user signature
	validSig, err := signature(uName, pWord)
	if !validSig {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Signature accepted!")
	}

	//if signature succeeded, get the other stuff
	//title := req.FormValue("title")
  //date := req.FormValue("date")
	//body := req.FormValue("body")
	//bodyHTML := paragraphize(body)
	//create new file for bodyHTML
	//insert into db new entry with path to the file made from bodyHTML
}

//redirect HTTP to HTTPS (works, but unused)
func redirectToHttps(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "localhost:8081"+req.RequestURI, http.StatusMovedPermanently)
}

//start the server
func main() {
	//TODO: actually get TLS working, dork

	//check if the cert files are available
	err := httpscerts.Check("../ssl/cert.pem", "../ssl/key.pem")
	//if they are not available, generate new ones
	if err != nil {
		err = httpscerts.Generate("../ssl/cert.pem", "../ssl/key.pem", "localhost:8081")
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}

	//create sql.DB and check for errors
  db, dbErr = sql.Open("mysql", dbOptions["dbAdmin"] + ":" + dbOptions["dbPWord"] + "@/" + dbOptions["dbName"])
  if dbErr != nil {
    panic(dbErr.Error())
  }
  defer db.Close()
  //test the connection to the database
  dbErr = db.Ping()
  if dbErr != nil {
    panic(dbErr.Error())
  }

	//create a new ServeMux for HTTP connections (delete once TLS works)
	httpMux := http.NewServeMux()
	httpMux.Handle("/", http.HandlerFunc(indexHandler))
	httpMux.Handle("/favicon.ico", http.HandlerFunc(faviconHandler))
	httpMux.Handle("/css/images/", http.HandlerFunc(imageHandler))
	httpMux.Handle("/blog/post/", http.HandlerFunc(blogHandler))

	//create a new ServeMux for HTTPS connections
	httpsMux := http.NewServeMux()
	httpsMux.Handle("/", http.HandlerFunc(indexHandler))
	httpsMux.Handle("/favicon.ico", http.HandlerFunc(faviconHandler))
	httpsMux.Handle("/css/images/", http.HandlerFunc(imageHandler))
	httpsMux.Handle("/blog/post/", http.HandlerFunc(blogHandler))

	//start HTTPS server in goroutine
	go http.ListenAndServeTLS(":8081", "../ssl/cert.pem", "../ssl/key.pem", httpsMux)

	//start HTTP server old-school and redirect to HTTPS
	http.ListenAndServe(":8080", httpMux) //replace with line below once TLS works
	//http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))
}
