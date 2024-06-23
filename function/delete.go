package forum2etoile

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

)

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
    // Récupère le cookie contenant le nom d'utilisateur
	cookie, err := r.Cookie("username")
	if err != nil {
        // Redirige vers la page de connexion si le cookie n'existe pas
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
    // Stocke la valeur du cookie (nom d'utilisateur) dans une variable
	username := cookie.Value
    // Initialise la base de données
	db := InitDatabase("database/db.db")
    // Ferme la base de données à la fin de la fonction
	defer db.Close()

    // Récupère l'ID du post depuis l'URL
	postIdStr := r.URL.Path[len("/delete/"):]
    // Convertit l'ID du post en entier
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
        // Retourne une erreur si l'ID du post n'est pas valide
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

    // Vérifie si le post appartient à l'utilisateur connecté
	var author string
	err = db.QueryRow("SELECT author FROM post WHERE id = ?", postId).Scan(&author)
	if err != nil {
        // Retourne une erreur si le post n'est pas trouvé
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

    // Vérifie si l'auteur du post est le même que l'utilisateur connecté
	if author != username {
        // Retourne une erreur si l'utilisateur essaie de supprimer un post qui ne lui appartient pas
		http.Error(w, "You can only delete your own posts", http.StatusForbidden)
		return
	}

    // Supprime le post de la base de données
	_, err = db.Exec("DELETE FROM post WHERE id = ?", postId)
	if err != nil {
        // Retourne une erreur si la suppression échoue
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

    // Redirige vers la page d'accueil après suppression
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

func DeleteConfirmationHandler(w http.ResponseWriter, r *http.Request) {
    // Récupère le cookie contenant le nom d'utilisateur
	cookie, err := r.Cookie("username")
	if err != nil {
        // Redirige vers la page de connexion si le cookie n'existe pas
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
    // Stocke la valeur du cookie (nom d'utilisateur) dans une variable
	username := cookie.Value
    // Initialise la base de données
	db := InitDatabase("database/db.db")
    // Ferme la base de données à la fin de la fonction
	defer db.Close()

    // Récupère l'ID du post depuis l'URL
	postIdStr := r.URL.Path[len("/delete/confirm/"):]
    // Convertit l'ID du post en entier
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
        // Retourne une erreur si l'ID du post n'est pas valide
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
    // Initialise une variable pour stocker les détails du post
	var post Post
    // Récupère les détails du post depuis la base de données
	err = db.QueryRow("SELECT id, title, author, content, date FROM post WHERE id = ?", postId).Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.Date)
	if err != nil {
        // Retourne une erreur si le post n'est pas trouvé
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

    // Vérifie si l'auteur du post est le même que l'utilisateur connecté
	if post.Author != username {
        // Retourne une erreur si l'utilisateur essaie de supprimer un post qui ne lui appartient pas
		http.Error(w, "You can only delete your own posts", http.StatusForbidden)
		return
	}

    // Charge et exécute le template de confirmation de suppression
	t := template.Must(template.ParseFiles("delete.html"))
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

    // Supprime les likes de l'utilisateur
    _, err = tx.Exec(`DELETE FROM like WHERE author = ?`, username)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Supprime les commentaires 
    _, err = tx.Exec(`DELETE FROM comment WHERE author = ?`, username)
    if err != nil {
    // Annule la transaction en cas d'erreur.
        tx.Rollback()
        return err
    }

    // Supprime les posts de l'utilisateur.
    _, err = tx.Exec(`DELETE FROM post WHERE author = ?`, username)
    if err != nil {
        tx.Rollback()
        return err
    }
 	// Supprime l'utilisateur de la table "register" en utilisant son pseudo.
    _, err = tx.Exec(`DELETE FROM register WHERE pseudo = ?`, username)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Valide la transaction si toutes les opérations ont réussi.
    err = tx.Commit()
	// Retourne une erreur si l'operation échoue.
    if err != nil {
        return err
    }

    return nil
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	db := InitDatabase("database/db.db")
	defer db.Close()

    // Récupère le cookie contenant le nom d'utilisateur.
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
    // Appelle la fonction pour supprimer le compte de l'utilisateur.
	err = DeleteAccount(db, username)
	if err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

    // Supprime le cookie après la suppression du compte.
	cookie = &http.Cookie{
		Name:   "username",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}
