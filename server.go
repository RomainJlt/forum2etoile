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

func getBookLastID() int {
	db := initDatabase("database/db.db")
	defer db.Close()

	var id int

	rows, err := db.Query("select ifnull(max(id), 0) as id from post")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
	}
	return id + 1
}



func insertIntoRegister(db *sql.DB, pseudo string, email string, password string, image string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO register (pseudo, email, password, image, post, subscribers) values (?, ?, ?, ?, 0, 0)`, pseudo, email, password, image)
	return result.LastInsertId()
}

// func insertIntoPost(db *sql.DB, title string, content string, author string, category string) (int64, error) {
// 	result, err := db.Exec(`INSERT INTO post (author, date, title, content, like, dislike, filter, category) values (?, ?, ?, ?, 0, 0, 0, ?)`, author, time.Now(), title, content, category)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return result.LastInsertId()
// }


func insertIntoPost(db *sql.DB, title string, content string, author string, category string) (int64, error) {
    formattedDate := time.Now().Format("02/01/2006 15:04")
    result, err := db.Exec(`INSERT INTO post (author, date, title, content, like, dislike, filter, category) values (?, ?, ?, ?, 0, 0, 0, ?)`, author, formattedDate, title, content, category)
    if err != nil {
        return 0, err
    }
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
func insertCategory(db *sql.DB, name string) (int64, error) {
    result, err := db.Exec(`INSERT INTO category (name) values (?)`, name)
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}



// func getPostData() []PostData {
//     db := initDatabase("database/db.db")
//     defer db.Close()
    
//     var posts []PostData
//     rows, err := db.Query(`SELECT id, author, date, title, content, like, dislike, filter, category FROM post`)
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer rows.Close()

//     for rows.Next() {
//         var post PostData
//         err := rows.Scan(&post.Id, &post.Author, &post.Date, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.Filter, &post.Category)
//         if err != nil {
//             log.Fatal(err)
//         }
//         posts = append(posts, post)
//     }
//     if err = rows.Err(); err != nil {
//         log.Fatal(err)
//     }
    
//     return posts
// }
func getPostData() []PostData {
    db := initDatabase("database/db.db")
    defer db.Close()
    
    var posts []PostData
    rows, err := db.Query(`SELECT id, author, date, title, content, like, dislike, filter, category FROM post`)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var post PostData
        var postDate string
        err := rows.Scan(&post.Id, &post.Author, &postDate, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.Filter, &post.Category)
        if err != nil {
            log.Fatal(err)
        }

        postTime, _ := time.Parse("02/01/2006 15:04", postDate)
        post.Date = timeAgo(postTime)
        posts = append(posts, post)
    }
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }
    
    return posts
}



// func getCommentData(idInfo int) {
// 	db := initDatabase("database/db.db")
// 	var temp Post

// 	rows, _ :=
// 		db.Query("SELECT author, content, date FROM comment WHERE postid = ?", idInfo)
// 	allResult = nil
// 	for rows.Next() {
// 		rows.Scan(&temp.AuthorComment, &temp.ContentComment, &temp.DateComment)
// 		allResult = append(allResult, temp)
// 	}
// }

func getCommentData(idInfo int) {
    db := initDatabase("database/db.db")
    var temp Post

    rows, _ := db.Query("SELECT author, content, date FROM comment WHERE postid = ?", idInfo)
    allResult = nil
    for rows.Next() {
        var commentDate string
        rows.Scan(&temp.AuthorComment, &temp.ContentComment, &commentDate)
        commentTime, _ := time.Parse("2006-01-02 15:04", commentDate)
        temp.DateComment = timeAgo(commentTime)
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
    posts := getPostData()
    t, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, posts)
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
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	db := initDatabase("database/db.db")
	titleForm := r.FormValue("inputTitle")
	contentForm := r.FormValue("inputContent")
	categoryForm := r.FormValue("category")
	user.Post = 0

	if titleForm != "" && contentForm != "" && categoryForm != "" {
		category := categoryForm
		fmt.Println(category)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		_, err = insertIntoPost(db, titleForm, contentForm, username, category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db.Exec(`UPDATE register SET post = post + 1 WHERE pseudo = ?`, username)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
		return
	}

	t, _ := template.ParseFiles("post.html")
	categories := getCategories(db)
	t.Execute(w, map[string]interface{}{
		"Categories": categories,
	})
}


func getCategories(db *sql.DB) []Category {
	rows, err := db.Query("SELECT id, name FROM category")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			log.Fatal(err)
		}
		categories = append(categories, category)
	}
	return categories
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


func searchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.FormValue("q")
    if query == "" {
        http.Redirect(w, r, "/index", http.StatusSeeOther)
        return
    }
    results := searchPosts(query)
    t, err := template.ParseFiles("search.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, results)
}


func searchPosts(query string) []PostData {
    db := initDatabase("database/db.db")
    defer db.Close()

    query = "%" + query + "%"
    var results []PostData
    rows, err := db.Query(`SELECT id, author, date, title, content, like, dislike, filter, category 
                           FROM post 
                           WHERE lower(title) LIKE lower(?) 
                           OR lower(content) LIKE lower(?) 
                           OR lower(date) LIKE lower(?) 
                           OR like = ? 
                           OR dislike = ? 
                           OR lower(author) LIKE lower(?) 
                           OR lower(category) LIKE lower(?)`, 
                           query, query, query, query, query, query, query)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var post PostData
        err := rows.Scan(&post.Id, &post.Author, &post.Date, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.Filter, &post.Category)
        if err != nil {
            log.Fatal(err)
        }
        results = append(results, post)
    }
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }

    return results
}


func timeAgo(t time.Time) string {
    now := time.Now()
    duration := now.Sub(t)
    
    if duration.Hours() < 24 {
        if duration.Hours() < 1 {
            if duration.Minutes() < 1 {
                return "à l'instant"
            }
            return fmt.Sprintf("il y a %.0f minutes", duration.Minutes())
        }
        return fmt.Sprintf("il y a %.0f heures", duration.Hours())
    } else {
        days := int(duration.Hours() / 24)
        return fmt.Sprintf("il y a %d jours", days)
    }
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
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
