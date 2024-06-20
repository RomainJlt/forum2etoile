package forum2etoile

import (
		
	_ "github.com/mattn/go-sqlite3"
)

type Login struct {
	Name  string
	Email string
	Image string
	Post  int
	Id    int
	Sub   int
}

type Post struct {
	Title          string
	Content        string
	Date           string
	Id             int
	Like           int
	Dislike        int
	Image          string
	Author         string
	Filter         int
	category    string
	//Username string 
	AuthorComment  string
	ContentComment string
	DateComment    string
	// CountCom       int
}

type PostData struct {
	Id       int
    Author   string
    Date     string
    Title    string
    Content  string
    Like     int
    Dislike  int
    Filter   int
    Category string
}

var user Login
var allUser []Login
var allResult []Post
var allData []PostData

type Register struct {
	Id       int
	Pseudo   string
	Email    string
	Password string
	Log      int
}

type Category struct {
	Id   int
	Name string
}
