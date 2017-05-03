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
