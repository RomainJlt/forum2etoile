package forum2etoile

import (
	"net/http"
	_ "github.com/mattn/go-sqlite3"

)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	likeId := r.URL.Path[6:]
	redirect := "/info/" + likeId
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	CheckLike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	likeId := r.URL.Path[9:]
	redirect := "/info/" + likeId
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	CheckDislike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}