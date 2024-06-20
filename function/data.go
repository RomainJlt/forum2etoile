package forum2etoile

import (
	"log"
	
)


func GetBookLastID() int {
	db := InitDatabase("database/db.db")
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

func GetPostData() []PostData {
    db := InitDatabase("database/db.db")
    defer db.Close()
    
    var posts []PostData
    rows, err := db.Query(`SELECT id, author, date, title, content, like, dislike, filter, category FROM post`)
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
        posts = append(posts, post)
    }
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }
    
    return posts
}


func GetCommentData(idInfo int) {
	db := InitDatabase("database/db.db")
	var temp Post

	rows, _ :=
		db.Query("SELECT author, content, date FROM comment WHERE postid = ?", idInfo)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.AuthorComment, &temp.ContentComment, &temp.DateComment)
		allResult = append(allResult, temp)
	}
}

func GetPostDataById(idInfo int) {
	db := InitDatabase("database/db.db")
	var temp PostData

	rows, _ :=
		db.Query("SELECT id, author, date, title, content, like, dislike, filter FROM post WHERE id = ?", idInfo)
	allData = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allData = append(allData, temp)
	}
}

func GetPostDataByFilter(filter int) {
	db := InitDatabase("database/db.db")
	var temp Post

	rows, _ :=
		db.Query(`SELECT * FROM post WHERE filter = ?`, filter)
	allResult = nil
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allResult = append([]Post{temp}, allResult...)
	}
}

func GetUserInfo(userInfo string) {
	db := InitDatabase("database/db.db")
	var temp Login

	rows, _ :=
		db.Query(`SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?`, userInfo)
	allUser = nil
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}

func GetUserInfoByCookie(username string) {
	db := InitDatabase("database/db.db")
	var temp Login

	rows, _ :=
		db.Query(`SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?`, username)
	allUser = nil
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}