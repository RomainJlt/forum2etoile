package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
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
	AuthorComment  string
	ContentComment string
	DateComment    string
}

type PostData struct {
	Title   string
	Content string
	Date    string
	Id      int
	Like    int
	Dislike int
	Image   string
	Author  string
	Filter  int
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

// Define new CookiesSessions
var store = sessions.NewCookieStore([]byte("mysession"))

func initDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "script/database/db.db")
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
			filter INT NOT NULL
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
		`
	_, err = db.Exec(sqlStmt)
	if err != nil {
	log.Fatal(err)
	}
	return db

}

func getBookLastID() int {
	db := initDatabase()
	defer db.Close()
	var id int

	err := db.QueryRow("select ifnull(max(id), 0) as id from post").Scan(&id)
	if err != nil {
		panic(err)
	}
	return id + 1
}

func insertIntoRegister(db *sql.DB, pseudo string, email string, password string, image string, creatDate string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO register (pseudo, email, password, image, post, subscribers , creatDate) values (?, ?, ?, ?, 0, 0, ?)`, pseudo, email, password, image, time.Now())
	return result.LastInsertId()
}

func insertIntoPost(db *sql.DB, title string, content string, author string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO post (author, date, title, content, like, dislike, filter) values (?, ?, ?, ?, 0, 0, 0)`, author, time.Now(), title, content)
	return result.LastInsertId()
}

func insertIntoComment(db *sql.DB, postid int, author string, content string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO comment (postid, date, author, content) values (?, ?, ?, ?)`, postid, time.Now(), author, content)
	return result.LastInsertId()
}

func insertIntoLike(db *sql.DB, postid string, author string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO like (postid, author, like, dislike) values (?, ?, 1, 1)`, postid, author)
	return result.LastInsertId()
}

func getPostData() {
	db := initDatabase()
	defer db.Close()
	var temp Post

	rows, _ := db.Query(`SELECT * FROM post`)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allResult = append([]Post{temp}, allResult...)
	}
}

func getCommentData(idInfo int) {
	db := initDatabase()
	defer db.Close()
	var temp Post

	rows, _ := db.Query("SELECT author, content, date FROM comment WHERE postid = ?", idInfo)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.AuthorComment, &temp.ContentComment, &temp.DateComment)
		allResult = append(allResult, temp)
	}
}