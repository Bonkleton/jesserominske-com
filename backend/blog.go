package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

//datatype interface for JSON decoder to represent blog posts
type BlogPost struct {
	Title string `json:"title"`
	Date string `json:"date"`
	Body string `json:"body"`
	UName string `json:"uname"`
	PWord string `json:"pword"`
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

//request handler for blog
func blogHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("blog")
	if req.Method != "POST" {
		fmt.Println(req.Method)
    http.Redirect(res, req, "localhost:8080", 301) //change this once appropriate
    return
  }
	//decode JSON data
	decoder := json.NewDecoder(req.Body)
  var b BlogPost
  err := decoder.Decode(&b)
  if err != nil {
    panic(err)
  }
  defer req.Body.Close()
	//validate user signature
  uName := b.UName
  pWord := b.PWord
	validSig, err := signature(uName, pWord)
	if !validSig {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Signature accepted!")
	}
	//if signature succeeded, get the other stuff
	title := b.Title
  date := b.Date
	body := []byte(b.Body)
	//create new file for bodyHTML
	//TODO: create new file based on those that already exist
	bodyPath := "../blog/test.html"
	err = ioutil.WriteFile(bodyPath, body, 0644)
	if err != nil {
    http.Error(res, err.Error(), http.StatusInternalServerError)
    return
  }
	//insert into db new entry with path to the file made from bodyHTML
	_, err = db.Exec("INSERT INTO blog(title, date, path) VALUES(?, ?, ?)", title, date, bodyPath)
  if err != nil {
    http.Error(res, err.Error(), http.StatusInternalServerError)
    return
  }
	fmt.Println("Blog posted successfully!")
}
