package main

import (
	"io/ioutil"
)

// container for the resources we will need to load into index.html
type Resources struct {
	Home          []byte
	Blog          []byte
	About         []byte
	Connect       []byte
	ContentScript []byte
}

// loads a partial
func loadResource(file string) ([]byte, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// runs the grunt-task
func getResources() (*Resources, error) {
	// load partials
	home, e1 := loadResource("view/home.html")
	blog, e2 := loadResource("view/blog.html")
	about, e3 := loadResource("view/about.html")
	connect, e4 := loadResource("view/connect.html")

	// load scripts
	contentScript, e5 := loadResource("js/content.js")

	// error checking
	if e1 != nil {
		return nil, e1
	} else if e2 != nil {
		return nil, e2
	} else if e3 != nil {
		return nil, e3
	} else if e4 != nil {
		return nil, e4
	} else if e5 != nil {
		return nil, e5
	}

	// return Grunt object of resources to insert
	return &Resources{
		Home:          home,
		Blog:          blog,
		About:         about,
		Connect:       connect,
		ContentScript: contentScript}, nil
}
