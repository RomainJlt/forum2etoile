package forums_Les_Etoiles

import (
	"database/sql"
	"time"
)

type Login struct {
	Name  string
	Email string
	Image string
	Post  int
	Id    int
	Sub   int
}

var User Login
var AllUser []Login

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

var AllResult []Post
var AllData []PostData

type Register struct {
	Id       int
	Pseudo   string
	Email    string
	Password string
	Log      int
}

func InsertIntoPost(db *sql.DB, title string, content string, author string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO post (author, date, title, content, like, dislike, filter) values (?, ?, ?, ?, 0, 0, 0)`, author, time.Now(), title, content)
	return result.LastInsertId()
}

func InsertIntoComment(db *sql.DB, postid int, author string, content string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO comment (postid, date, author, content) values (?, ?, ?, ?)`, postid, time.Now(), author, content)
	return result.LastInsertId()
}

func InsertIntoRegister(db *sql.DB, pseudo string, email string, password string, image string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO register (pseudo, email, password, image, post, subscribers) values (?, ?, ?, ?, 0, 0)`, pseudo, email, password, image)
	return result.LastInsertId()
}
