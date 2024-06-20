package forum2etoile

import (
	
	_ "github.com/mattn/go-sqlite3"
)

func Login2(LogPseudo string, LogPassword string) bool {
	db := InitDatabase("database/db.db")
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

func CheckLike(username string, likeId string) {
	db := InitDatabase("database/db.db")
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
		InsertIntoLike(db, likeId, username)
		db.Exec("UPDATE post SET like = like + 1 WHERE id = ?", likeId)
		db.Exec("UPDATE like SET like = 2 WHERE author = ?", username)
	}
}


func CheckDislike(username string, likeId string) {
	db := InitDatabase("database/db.db")
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
		InsertIntoLike(db, likeId, username)
		db.Exec("UPDATE post SET dislike = dislike + 1 WHERE id = ?", likeId)
		db.Exec("UPDATE like SET dislike = 2 WHERE author = ?", username)
	}
}


func Register2(RegisterPseudo string, RegisterEmail string) bool {
	db := InitDatabase("database/db.db")

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