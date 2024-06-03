

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
	db := initDatabase()
	defer db.Close()
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

func getPostDataById(idInfo int) {
	db := initDatabase()
	defer db.Close()
	var temp PostData

	rows, _ := db.Query("SELECT id, author, date, title, content, like, dislike, filter FROM post WHERE id = ?", idInfo)
	allData = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allData = append(allData, temp)
	}
}

func getPostDataByFilter(filter int) {
	db := initDatabase()
	defer db.Close()
	var temp Post

	rows, _ := db.Query(`SELECT * FROM post WHERE filter = ?`, filter)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allResult = append([]Post{temp}, allResult...)
	}
}

func getUserInfo(userInfo string) {
	db := initDatabase()
	defer db.Close()
	var temp Login

	rows, _ := db.Query(`SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?`, userInfo)
	allUser = nil
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}

func getUserInfoByCookie(username string) {
	db := initDatabase()
	defer db.Close()
	var temp Login

	rows, _ := db.Query(`SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?`, username)
	allUser = nil
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}

// Vérifiez la fonction login pour vous assurer qu'elle fonctionne correctement
func login(LogPseudo string, LogPassword string) bool {
	db := initDatabase()
	var pseudo string
	var hashedPassword string
	var result = false

	rows, err := db.Query("SELECT pseudo, password FROM register WHERE pseudo = ?", LogPseudo)
	if err != nil {
		log.Println("Erreur lors de la requête : ", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&pseudo, &hashedPassword)
		if err != nil {
			log.Println("Erreur lors du scan : ", err)
			return false
		}
		if CheckPasswordHash(LogPassword, hashedPassword) { // Compare le mot de passe haché
			result = true
		}
	}
	if err = rows.Err(); err != nil {
		log.Println("Erreur après l'itération des lignes : ", err)
		return false
	}

	return result
}

func checkLike(username string, likeId string) {
	db := initDatabase()
	defer db.Close()
	var author string
	var postid int
	var like int
	var dislike int
	rows, _ := db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
	for rows.Next() {
		rows.Scan(&postid, &author, &like, &dislike)
	}
	if author != "" && postid != 0 && like != 0 {
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
	db := initDatabase()
	defer db.Close()
	var author string
	var postid int
	var like int
	var dislike int
	rows, _ := db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
	for rows.Next() {
		rows.Scan(&postid, &author, &like, &dislike)
	}
	if author != "" && postid != 0 && dislike != 0 {
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


func Home(w http.ResponseWriter, r *http.Request) {
	getPostData()
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, allResult)
}

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	db := initDatabase()
	defer db.Close()
	r.ParseForm()
	if r.Method == http.MethodPost {
		pseudo := r.FormValue("pseudo")
		email := r.FormValue("email")
		password := r.FormValue("password")
		hash, _ := HashPassword(password)
		image := "default.jpg"
		_, err := insertIntoRegister(db, pseudo, email, hash, image)
		if err == nil {
			http.Redirect(w, r, "./templates/login.html", http.StatusSeeOther)
		}
	}
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == http.MethodPost {
		pseudo := r.FormValue("pseudo")
		password := r.FormValue("password")
		if login(pseudo, password) {
			session, _ := store.Get(r, "mysession")
			session.Values["username"] = pseudo
			err := session.Save(r, w)
			if err != nil {
				log.Println("Erreur lors de la sauvegarde de la session : ", err)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
		} else {
			log.Println("Identifiants invalides pour l'utilisateur : ", pseudo)
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

}




// func Index(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	session, _ := store.Get(r, "mysession")
// 	username := session.Values["username"]
// 	fmt.Println(username)
// 	if username != nil {
// 		getUserInfoByCookie(username.(string))
// 		data := struct {
// 			User []Login
// 		}{
// 			allUser,
// 		}
// 		tmpl := template.Must(template.ParseFiles("templates/index.html"))
// 		tmpl.Execute(w, data)
// 	} else {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 	}
// }



func Index(w http.ResponseWriter, r *http.Request) {
	badgeInformatique := r.FormValue("badgeInformatique")
	badgeSport := r.FormValue("badgeSport")
	badgeMusique := r.FormValue("badgeMusique")
	badgeGame := r.FormValue("badgeGame")
	badgeFood := r.FormValue("badgeFood")
	if badgeInformatique == "1" {
		getPostDataByFilter(1)
	} else if badgeSport == "2" {
		getPostDataByFilter(2)
	} else if badgeMusique == "3" {
		getPostDataByFilter(3)
	} else if badgeGame == "4" {
		getPostDataByFilter(4)
	} else if badgeFood == "5" {
		getPostDataByFilter(5)
	} else {
		getPostData()
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, allResult) //Execute the value of all post
}

func Profile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	username := fmt.Sprintf("%v", session.Values["username"]) //Decrypts data of the session cookies
	getUserInfoByCookie(username)
	t, _ := template.ParseFiles("templates/profile.html")
	t.Execute(w, allUser)
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"]
	if username != nil {
		if r.Method == http.MethodPost {
			title := r.FormValue("title")
			content := r.FormValue("content")
			insertIntoPost(initDatabase(), title, content, username.(string))
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		tmpl := template.Must(template.ParseFiles("templates/newpost.html"))
		tmpl.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}


func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	//Delete and save coockie sessions
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}



// func PostDetail(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(r.URL.Query().Get("id"))
// 	if err == nil {
// 		getPostDataById(id)
// 		getCommentData(id)
// 		data := struct {
// 			Post    PostData
// 			Comment []Post
// 		}{
// 			allData[0],
// 			allResult,
// 		}
// 		tmpl := template.Must(template.ParseFiles("templates/post.html"))
// 		tmpl.Execute(w, data)
// 	}
// }


//HandleFunc for post.html (Get and post data)
func postHandler(w http.ResponseWriter, r *http.Request) {
	// db := initDatabase("database/db.db")
	db := initDatabase()
	titleForm := r.FormValue("inputEmail")
	contentForm := r.FormValue("inputPassword")
	session, _ := store.Get(r, "mysession")
	username := fmt.Sprintf("%v", session.Values["username"]) //Decrypts data of the session cookies
	user.Post = 0

	informatique := r.FormValue("badgeInformatique")
	sport := r.FormValue("badgeSport")
	musique := r.FormValue("badgeMusique")
	jeux := r.FormValue("badgeGame")
	food := r.FormValue("badgeFood")
	lastid := getBookLastID()

	if titleForm != "" && contentForm != "" {
		insertIntoPost(db, titleForm, contentForm, username) //insert the value send by the user in the database
		db.Exec(`INSERT INTO post (date) values (?)`, time.Now())
		db.Exec(`UPDATE register SET post = post + 1 WHERE pseudo = ?`, username) //Update the value in the database
		//Get the value of checkbox, and update the database
		if informatique == "1" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 1, lastid)
		}
		if sport == "2" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 2, lastid)
		}
		if musique == "3" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 3, lastid)
		}
		if jeux == "4" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 4, lastid)
		}
		if food == "5" {
			db.Exec(`UPDATE post SET filter = ? WHERE id = ?`, 5, lastid)
		}

		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	t, _ := template.ParseFiles("post.html")
	t.Execute(w, nil)

}



func CommentPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"]
	if username != nil {
		if r.Method == http.MethodPost {
			postid, _ := strconv.Atoi(r.FormValue("postid"))
			content := r.FormValue("content")
			insertIntoComment(initDatabase(), postid, username.(string), content)
			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postid), http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"]
	if username != nil {
		if r.Method == http.MethodPost {
			postid := r.FormValue("postid")
			checkLike(username.(string), postid)
			http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postid), http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"]
	if username != nil {
		if r.Method == http.MethodPost {
			postid := r.FormValue("postid")
			checkDislike(username.(string), postid)
			http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postid), http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

func main() {
	fs := http.FileServer(http.Dir("database"))
	http.Handle("/database/", http.StripPrefix("/database/", fs))

	http.HandleFunc("/", Home)
	http.HandleFunc("/register", RegisterPost)
	http.HandleFunc("/login", LoginPost)
	http.HandleFunc("/index", Index)
	http.HandleFunc("/profile", Profile)
	http.HandleFunc("/newpost", NewPost)
	http.HandleFunc("/post", postHandler)
	// http.HandleFunc("/post", PostDetail)
	http.HandleFunc("/comment", CommentPost)
	http.HandleFunc("/like", LikePost)
	http.HandleFunc("/dislike", DislikePost)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
