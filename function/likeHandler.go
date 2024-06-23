package forum2etoile

import (
	"net/http"
	_ "github.com/mattn/go-sqlite3"

)

// LikeHandler gère les requêtes pour ajouter un like à un post.
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du post à partir de l'URL.
	likeId := r.URL.Path[6:]
	redirect := "/info/" + likeId
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	// Vérifie si l'utilisateur a déjà aimé le post.
	CheckLike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

// DislikeHandler gère les requêtes pour ajouter un dislike à un post.
func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du post à partir de l'URL
	likeId := r.URL.Path[9:]
	redirect := "/info/" + likeId
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	// Vérifie si l'utilisateur a déjà disliké le post.
	username := cookie.Value
	// Vérifie si l'utilisateur a déjà disliké le post.
	CheckDislike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}