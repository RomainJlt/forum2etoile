package forum2etoile

import (
	"database/sql"
    "time"
    _ "github.com/mattn/go-sqlite3"
)
//Insère les données de l'utilisateur dans la table register.
func InsertIntoRegister(db *sql.DB, pseudo string, email string, password string, image string) (int64, error) {
    // Exécute une requête SQL pour insérer les données de l'utilisateur dans la table register.
	result, _ := db.Exec(`INSERT INTO register (pseudo, email, password, image, post, subscribers) values (?, ?, ?, ?, 0, 0)`, pseudo, email, password, image)
	
    return result.LastInsertId()
}

// Insère les données du post dans la table post.
func InsertIntoPost(db *sql.DB, title string, content string, author string, category string) (int64, error) {
    time := time.Now().Format("02/01/2006 15:04")
    result, err := db.Exec(`INSERT INTO post (author, date, title, content, like, dislike, filter, category) values (?, ?, ?, ?, 0, 0, 0, ?)`, author, time, title, content, category)
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

// Insère les données du commentaire dans la table comment.
func InsertIntoComment(db *sql.DB, postid int, author string, content string) (int64, error) {
    formattedDate := time.Now().Format("02/01/2006 15:04")
    result, _ := db.Exec(`INSERT INTO comment (postid, date, author, content) values (?, ?, ?, ?)`, postid, formattedDate, author, content)
    return result.LastInsertId()
}

// Insère les données du like dans la table like.
func InsertIntoLike(db *sql.DB, postid string, author string) (int64, error) {
	result, _ := db.Exec(`INSERT INTO like (postid, author, like, dislike) values (?, ?, 1, 1)`, postid, author)
	return result.LastInsertId()
}

// Insère les données de la catégorie dans la table category.
func InsertCategory(db *sql.DB, name string) (int64, error) {
    result, err := db.Exec(`INSERT INTO category (name) values (?)`, name)
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}