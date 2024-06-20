package main

import (
	"fmt"
	"log"
	"net/http"
	forum2etoile "forum2etoile/function"

	_ "github.com/mattn/go-sqlite3"
	
)


func main() {
	fs := http.FileServer(http.Dir(""))
	http.Handle("/static/", http.StripPrefix("/static/", fs)) 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}) 
	http.HandleFunc("/register", forum2etoile.RegisterHandler)
	http.HandleFunc("/index", forum2etoile.IndexHandler)
	http.HandleFunc("/profile", forum2etoile.ProfileHandler)
	http.HandleFunc("/user/", forum2etoile.UserHandler)
	http.HandleFunc("/post", forum2etoile.PostHandler)
	http.HandleFunc("/info/", forum2etoile.InfoHandler)
	http.HandleFunc("/logout", forum2etoile.LogoutHandler)
	http.HandleFunc("/like/", forum2etoile.LikeHandler)
	http.HandleFunc("/dislike/", forum2etoile.DislikeHandler)
	http.HandleFunc("/search", forum2etoile.SearchHandler)
	http.HandleFunc("/update-profile", forum2etoile.UpdateProfileHandler)
	http.HandleFunc("/delete/", forum2etoile.DeletePostHandler)
	http.HandleFunc("/delete/confirm/", forum2etoile.DeleteConfirmationHandler)
	http.HandleFunc("/deleteAccount", forum2etoile.DeleteAccountHandler)
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
