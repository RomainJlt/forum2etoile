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
	AuthorComment  string
	ContentComment string
	DateComment    string
	// CountCom       int
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
	db := initDatabase("database/db.db/")
	var id int

	err := db.QueryRow("select ifnull(max(id), 0) as id from post").Scan(&id)
	if err != nil {
		panic(err)
	}
	return id + 1
}

func insertIntoRegister(db *sql.DB, pseudo string, email string, password string, image string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO register (pseudo, email, password, image, post, subscribers) values (?, ?, ?, ?, 0, 0)`, pseudo, email, password, image)
	return result.LastInsertId()
}

func insertIntoPost(db *sql.DB, title string, content string, author string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO post (author, date, title, content, like, dislike, filter) values (?, ?, ?, ?, 0, 0, 0)`, author, time.Now(), title, content)
	return result.LastInsertId()
}

func insertIntoComment(db *sql.DB, postid int, author string, content string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO comment (postid, date, author, content) values (?, ?, ?, ?)`, postid, "0", author, content)
	return result.LastInsertId()
}

func insertIntoLike(db *sql.DB, postid string, author string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO like (postid, author, like, dislike) values (?, ?, 1, 1)`, postid, author)
	return result.LastInsertId()
}


func getPostData() {
	db := initDatabase("database/db.db")
	var temp Post

	rows, _ :=
		db.Query(`SELECT * FROM post`)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allResult = append([]Post{temp}, allResult...)
	}
}

func getCommentData(idInfo int) {
	db := initDatabase("database/db.db")
	var temp Post

	rows, _ :=
		db.Query("SELECT author, content, date FROM comment WHERE postid = ?", idInfo)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.AuthorComment, &temp.ContentComment, &temp.DateComment)
		allResult = append(allResult, temp)
	}
}

func getPostDataById(idInfo int) {
	db := initDatabase("database/db.db")
	var temp PostData

	rows, _ :=
		db.Query("SELECT id, author, date, title, content, like, dislike, filter FROM post WHERE id = ?", idInfo)
	allData = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allData = append(allData, temp)
	}
}

func getPostDataByFilter(filter int) {
	db := initDatabase("database/db.db")
	var temp Post

	rows, _ :=
		db.Query(`SELECT * FROM post WHERE filter = ?`, filter)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allResult = append([]Post{temp}, allResult...)
	}
}

func getUserInfo(userInfo string) {
	db := initDatabase("database/db.db")
	var temp Login

	rows, _ :=
		db.Query(`SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?`, userInfo)
	allUser = nil
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}

func getUserInfoByCookie(username string) {
	db := initDatabase("database/db.db")
	var temp Login

	rows, _ :=
		db.Query(`SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?`, username)
	allUser = nil
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}

func login(LogPseudo string, LogPassword string) bool {
	db := initDatabase("database/db.db")
	var pseudo string
	var password string
	var result = false
	rows, _ :=
		db.Query("SELECT pseudo, password FROM register")
	for rows.Next() {
		rows.Scan(&pseudo, &password)
		if LogPseudo == pseudo && CheckPasswordHash(LogPassword, password) { 
			result = true
		}
	}
	return result
}

func checkLike(username string, likeId string) {
	db := initDatabase("database/db.db")
	var author string
	var postid int
	var like int
	var dislike int
	rows, _ := db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
	for rows.Next() {
		rows.Scan(&postid, &author, &like, &dislike)
	}
	if author != "" && postid != 0 && like != 0 {
		db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
		if like == 1 && dislike == 1 {
			db.Exec("UPDATE post SET like = like + 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET like = 2 WHERE author = ?", username)
		}
		if like == 2 {
			db.Exec("UPDATE post SET like = like - 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET like = 1 WHERE author = ?", username)
		}
	} else {
		insertIntoLike(db, likeId, username)
		db.Exec("UPDATE post SET like = like + 1 WHERE id = ?", likeId)
		db.Exec("UPDATE like SET like = 2 WHERE author = ?", username)
	}
}


func checkDislike(username string, likeId string) {
	db := initDatabase("database/db.db")
	var author string
	var postid int
	var like int
	var dislike int
	rows, _ := db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
	for rows.Next() {
		rows.Scan(&postid, &author, &like, &dislike)
	}
	if author != "" && postid != 0 && dislike != 0 {
		db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
		if like == 1 && dislike == 1 {
			db.Exec("UPDATE post SET dislike = dislike + 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET dislike = 2 WHERE author = ?", username)
		}
		if dislike == 2 {
			db.Exec("UPDATE post SET dislike = dislike - 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET dislike = 1 WHERE author = ?", username)
		}
	} else {
		insertIntoLike(db, likeId, username)
		db.Exec("UPDATE post SET dislike = dislike + 1 WHERE id = ?", likeId)
		db.Exec("UPDATE like SET dislike = 2 WHERE author = ?", username)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func register(RegisterPseudo string, RegisterEmail string) bool {
	db := initDatabase("database/db.db")

	var pseudo string
	var email string
	var result = true
	rows, _ :=
		db.Query("SELECT  pseudo, email FROM register")
	for rows.Next() {
		rows.Scan(&pseudo, &email)
		if RegisterPseudo == pseudo || RegisterEmail == email { 
			result = false
		}
	}
	return result
}



func indexHandler(w http.ResponseWriter, r *http.Request) {
    
    getPostData()
    
    
    t, _ := template.ParseFiles("index.html")
    t.Execute(w, allResult)
}


func registerHandler(w http.ResponseWriter, r *http.Request) {
	pseudoForm := r.FormValue("pseudoCreate")
	emailForm := r.FormValue("emailCreate")
	passwordForm := r.FormValue("passwordCreate")
	imageForm := r.FormValue("imageCreate")
	pseudoLog := r.FormValue("pseudoLog")
	passwordLog := r.FormValue("passwordLog")

	db := initDatabase("database/db.db")

	hash, _ := HashPassword(passwordForm)
	if pseudoForm != "" && emailForm != "" && passwordForm != "" {
		if register(pseudoForm, emailForm) { //If true
			if imageForm != "" {
				insertIntoRegister(db, pseudoForm, emailForm, hash, imageForm)
			} else {
				insertIntoRegister(db, pseudoForm, emailForm, hash, "http://marclimoservices.com/wp-content/uploads/2017/05/facebook-default.png")
			}
		}
	}
	
	if login(pseudoLog, passwordLog) {
		user.Name = pseudoLog
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{Name: "username", Value: pseudoLog, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("register.html")
	t.Execute(w, nil)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	getUserInfoByCookie(username)
	t, _ := template.ParseFiles("profile.html")
	t.Execute(w, allUser)
}


func userHandler(w http.ResponseWriter, r *http.Request) {
	userInfo := r.URL.Path[6:]
	getUserInfo(userInfo)
	t, _ := template.ParseFiles("user.html")
	t.Execute(w, allUser)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "username", Value: "", Expires: time.Unix(0, 0), MaxAge: -1}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
	likeId := r.URL.Path[6:]
	redirect := "/info/" + likeId
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	checkLike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func dislikeHandler(w http.ResponseWriter, r *http.Request) {
	likeId := r.URL.Path[9:]
	redirect := "/info/" + likeId
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	checkDislike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}



func postHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	db := initDatabase("database/db.db/")
	titleForm := r.FormValue("inputEmail")
	contentForm := r.FormValue("inputPassword")
	user.Post = 0

	if titleForm != "" && contentForm != "" {
		insertIntoPost(db, titleForm, contentForm, username)
		db.Exec(`INSERT INTO post (date) values (?)`, time.Now())
		db.Exec(`UPDATE register SET post = post + 1 WHERE pseudo = ?`, username)

		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	t, _ := template.ParseFiles("post.html")
	t.Execute(w, nil)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	db := initDatabase("database/db.db/")
	idInfo, _ := strconv.Atoi(r.URL.Path[6:])
	contentComment := r.FormValue("commentArea")
	redirect := "/info/" + strconv.Itoa(idInfo)

	getPostDataById(idInfo)

	if len(contentComment) > 0 {
		insertIntoComment(db, idInfo, username, contentComment)
		db.Exec(`UPDATE comment SET date = ? WHERE postid = ?`, time.Now(), idInfo)
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}

	getCommentData(idInfo)

	m := map[string]interface{}{
		"Results": allResult,
		"Post":    allData,
	}
	t := template.Must(template.ParseFiles("info.html"))
	t.Execute(w, m)
}


func main() {
	fs := http.FileServer(http.Dir(""))
	http.Handle("/static/", http.StripPrefix("/static/", fs)) 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}) 
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/info/", infoHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/like/", likeHandler)
	http.HandleFunc("/dislike/", dislikeHandler)
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}