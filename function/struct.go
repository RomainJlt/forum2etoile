package forum2etoile

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Login est une structure qui contient les informations de l'utilisateur connecté.
type Login struct {
	Name  string
	Email string
	Image string
	Post  int
	Id    int
	Sub   int
}

// Post est une structure qui contient les informations d'un post.
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
	Category       string
	AuthorComment  string
	ContentComment string
	DateComment    string
	FormattedDate  string
	
}

// PostData est une structure qui contient les informations d'un post.
type PostData struct {
	Id       int
	Author   string
	Date     time.Time
	Title    string
	Content  string
	Like     int
	Dislike  int
	Filter   int
	Category string
	Ftime    string
}




// Login est une structure qui contient les informations de l'utilisateur connecté.
var user Login
var allUser []Login
var allResult []Post
var allData []PostData

// Register est une structure qui contient les informations d'un utilisateur.
type Register struct {
	Id       int
	Pseudo   string
	Email    string
	Password string
	Log      int
}

// Category est une structure qui contient les informations d'une catégorie.
type Category struct {
	Id   int
	Name string
}
