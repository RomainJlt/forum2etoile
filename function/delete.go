package forum2etoile

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

)


func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	db := InitDatabase("database/db.db")
	defer db.Close()

	postIdStr := r.URL.Path[len("/delete/"):]
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if the post belongs to the logged-in user
	var author string
	err = db.QueryRow("SELECT author FROM post WHERE id = ?", postId).Scan(&author)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if author != username {
		http.Error(w, "You can only delete your own posts", http.StatusForbidden)
		return
	}

	// Delete the post
	_, err = db.Exec("DELETE FROM post WHERE id = ?", postId)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/index", http.StatusSeeOther)
}


func DeleteConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	db := InitDatabase("database/db.db")
	defer db.Close()

	postIdStr := r.URL.Path[len("/delete/confirm/"):]
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var post Post
	err = db.QueryRow("SELECT id, title, author, content, date FROM post WHERE id = ?", postId).Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.Date)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if post.Author != username {
		http.Error(w, "You can only delete your own posts", http.StatusForbidden)
		return
	}

	t := template.Must(template.ParseFiles("templates/delete.html"))
	t.Execute(w, map[string]interface{}{
		"Post": post,
	})
}

func DeleteAccount(db *sql.DB, username string) error {
    // Début de la transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    // Supprimer les likes et dislikes
    _, err = tx.Exec(`DELETE FROM like WHERE author = ?`, username)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Supprimer les commentaires
    _, err = tx.Exec(`DELETE FROM comment WHERE author = ?`, username)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Supprimer les posts
    _, err = tx.Exec(`DELETE FROM post WHERE author = ?`, username)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Supprimer l'utilisateur
    _, err = tx.Exec(`DELETE FROM register WHERE pseudo = ?`, username)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Commit de la transaction
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	db := InitDatabase("database/db.db")
	defer db.Close()

	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	err = DeleteAccount(db, username)
	if err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	cookie = &http.Cookie{
		Name:   "username",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}
