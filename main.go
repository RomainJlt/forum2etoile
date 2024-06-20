package forum

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
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/update-profile", updateProfileHandler)
	http.HandleFunc("/delete/", deletePostHandler)
	http.HandleFunc("/delete/confirm/", deleteConfirmationHandler)
	http.HandleFunc("/deleteAccount", deleteAccountHandler)
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
