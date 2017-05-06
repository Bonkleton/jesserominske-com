package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

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
//TODO: make this agree with model
func blogHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("blog")
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
	//body := req.FormValue("body") //this will be paragraphized by the front end
	//create new file for bodyHTML
	//bodyPath := "../blog/test.html"
	//err := ioutil.WriteFile(bodyPath, body, 0644)
	//insert into db new entry with path to the file made from bodyHTML
	//_, err = db.Exec("INSERT INTO blog(title, date, path) VALUES(?, ?, ?)", title, date, bodyPath)
  //if err != nil {
  //  http.Error(res, "Server error - unable to create your blog entry.", 500)
  //  return
  //}
}
