package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	
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
// Initialise DataBase, and create it with his tables
func initDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", "database/db.db")
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
				CREATE TABLE IF NOT EXISTS register (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					pseudo TEXT NOT NULL, 
					email TEXT NOT NULL, 
					password TEXT NOT NULL,
					image TEXT NOT NULL,
					post INT NOT NULL,
					subscribers INT NOT NULL
				);

				CREATE TABLE IF NOT EXISTS post (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					author TEXT NOT NULL,
					date TEXT NOT NULL,
					title TEXT NOT NULL,
					content TEXT NOT NULL,
					like INT NOT NULL,
					dislike INT NOT NULL,
					filter INT NOT NULL,
					category TEXT NOT NULL
					
				);

				CREATE TABLE IF NOT EXISTS comment (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					postid INT NOT NULL,
					date TEXT NOT NULL,
					author TEXT NOT NULL,
					content TEXT NOT NULL
				);

				CREATE TABLE IF NOT EXISTS like (
					postid INTEGER NOT NULL,
					author TEXT NOT NULL,
					like INT NOT NULL,
					dislike INT NOT NULL,
					PRIMARY KEY (postid, author)
				);

					 CREATE TABLE IF NOT EXISTS category (
   					 id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
   					 name TEXT NOT NULL
				);
				`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	return db
	
}
