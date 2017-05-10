package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

//datatype interface for JSON decoder to represent blog posts
type BlogPost struct {
	Title string `json:"title"`
	Date string `json:"date"`
	Body string `json:"body"`
	UName string `json:"uname"`
	PWord string `json:"pword"`
}
//datatype for db results of extant blogs
type BlogRecord struct {
	Id string
	Title string
	Date string
	Body string
}
//datatype for login failure
type LoginResponse struct {
	Failure bool
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
//returns whatever file is currently available in blog folder
func availBlogFile() string {
	for fileCtr := 1; fileCtr < 10000; fileCtr++ {
		fileBase := "000" + strconv.Itoa(fileCtr)
		fileBase = fileBase[len(fileBase)-4:]
		fmt.Println(fileBase)
		_, err := os.Stat("../blog/" + fileBase + ".html")
		if os.IsNotExist(err) {
			return "../blog/" + fileBase + ".html"
		}
	}
	return ""
}
//function for getting all blogs in db
func processBlogs(blogs *[]BlogRecord, rows *sql.Rows, res http.ResponseWriter) {
	for rows.Next() {
		var (
			recordId string
			recordTitle string
			recordDate string
			recordBodyPath string
		)
		err := rows.Scan(&recordId, &recordTitle, &recordDate, &recordBodyPath)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		recordBody, err := ioutil.ReadFile(recordBodyPath)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		blogRecord := BlogRecord{Id: recordId,
														 Title: recordTitle,
													   Date: recordDate,
													   Body: string(recordBody)}
		*blogs = append(*blogs, blogRecord)
	}
}

//request handler for blog
func blogHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("blog")
	//if not a POST request just get blogs
	if req.Method != "POST" {
		//get current list of blogs from db
		rows, err := db.Query("SELECT * FROM blog")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
	    return
		}
		defer rows.Close()
		//process list of blogs and send response
		var blogs []BlogRecord
		processBlogs(&blogs, rows, res)
		json.NewEncoder(res).Encode(blogs)
  } else {
		//otherwise, decode the JSON
		decoder := json.NewDecoder(req.Body)
	  var b BlogPost
	  err := decoder.Decode(&b)
	  if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
	    return
	  }
	  defer req.Body.Close()
		//validate user signature
	  uName := b.UName
	  pWord := b.PWord
		validSig, err := signature(uName, pWord)
		if !validSig {
			json.NewEncoder(res).Encode(LoginResponse{Failure: true})
			return
		} else {
			fmt.Println("Signature accepted!")
		}
		//if signature succeeded, get the other stuff
		title := b.Title
	  date := b.Date
		body := []byte(b.Body)
		//create new file for bodyHTML
		bodyPath := availBlogFile()
		if bodyPath == "" {
			fmt.Println("No files available in blog folder!")
			return
		}
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
		//get current list of blogs from db
		rows, err := db.Query("SELECT * FROM blog")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
	    return
		}
		defer rows.Close()
		//process list of blogs and send response
		var blogs []BlogRecord
		processBlogs(&blogs, rows, res)
		json.NewEncoder(res).Encode(blogs)
	}
}
