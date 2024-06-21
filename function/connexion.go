package forum2etoile

import (
	
	_ "github.com/mattn/go-sqlite3"
)

// Vérifie les informations de connexion de l'utilisateur
func Login2(LogPseudo string, LogPassword string) bool {
	db := InitDatabase("database/db.db") // Initialise la base de données
	var pseudo string
	var password string
	var result = false
	
	// Exécute une requête SQL pour obtenir les pseudos et mots de passe enregistrés
	rows, _ := db.Query("SELECT pseudo, password FROM register")
	for rows.Next() {
		rows.Scan(&pseudo, &password)
		// Vérifie si le pseudo et le mot de passe correspondent
		if LogPseudo == pseudo && CheckPasswordHash(LogPassword, password) { 
			result = true
		}
	}
	return result // Retourne vrai si les informations de connexion sont correctes, sinon faux
}

// Vérifie si un utilisateur a aimé un post et met à jour en conséquence
func CheckLike(username string, likeId string) {
	db := InitDatabase("database/db.db") // Initialise la base de données
	var author string
	var postid int
	var like int
	var dislike int
	
	// Exécute une requête SQL pour obtenir les informations sur le like
	rows, _ := db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
	for rows.Next() {
		rows.Scan(&postid, &author, &like, &dislike)
	}
	
	// Vérifie si l'utilisateur a déjà aimé ce post
	if author != "" && postid != 0 && like != 0 {
		if like == 1 && dislike == 1 {
			// Met à jour le nombre de likes et l'état du like
			db.Exec("UPDATE post SET like = like + 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET like = 2 WHERE author = ?", username)
		}
		if like == 2 {
			// Réduit le nombre de likes et réinitialise l'état du like
			db.Exec("UPDATE post SET like = like - 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET like = 1 WHERE author = ?", username)
		}
	} else {
		// Insère un nouveau like et met à jour le post
		InsertIntoLike(db, likeId, username)
		db.Exec("UPDATE post SET like = like + 1 WHERE id = ?", likeId)
		db.Exec("UPDATE like SET like = 2 WHERE author = ?", username)
	}
}

// Vérifie si un utilisateur a disliké un post et met à jour en conséquence
func CheckDislike(username string, likeId string) {
	db := InitDatabase("database/db.db") // Initialise la base de données
	var author string
	var postid int
	var like int
	var dislike int
	
	// Exécute une requête SQL pour obtenir les informations sur le dislike
	rows, _ := db.Query("SELECT postid, author, like, dislike FROM like WHERE author = ? and postid = ?", username, likeId)
	for rows.Next() {
		rows.Scan(&postid, &author, &like, &dislike)
	}
	
	// Vérifie si l'utilisateur a déjà disliké ce post
	if author != "" && postid != 0 && dislike != 0 {
		if like == 1 && dislike == 1 {
			// Met à jour le nombre de dislikes et l'état du dislike
			db.Exec("UPDATE post SET dislike = dislike + 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET dislike = 2 WHERE author = ?", username)
		}
		if dislike == 2 {
			// Réduit le nombre de dislikes et réinitialise l'état du dislike
			db.Exec("UPDATE post SET dislike = dislike - 1 WHERE id = ?", likeId)
			db.Exec("UPDATE like SET dislike = 1 WHERE author = ?", username)
		}
	} else {
		// Insère un nouveau dislike et met à jour le post
		InsertIntoLike(db, likeId, username)
		db.Exec("UPDATE post SET dislike = dislike + 1 WHERE id = ?", likeId)
		db.Exec("UPDATE like SET dislike = 2 WHERE author = ?", username)
	}
}

// Enregistre un utilisateur après vérification si le pseudo ou l'email existe déjà
func Register2(RegisterPseudo string, RegisterEmail string) bool {
	db := InitDatabase("database/db.db") 

	var pseudo string
	var email string
	var result = true
	
	// Exécute une requête SQL pour obtenir les pseudos et emails enregistrés
	rows, _ := db.Query("SELECT pseudo, email FROM register")
	for rows.Next() {
		rows.Scan(&pseudo, &email)
		// Vérifie si le pseudo ou l'email existe déjà
		if RegisterPseudo == pseudo || RegisterEmail == email { 
			result = false
		}
	}
	return result // Retourne vrai si le pseudo et l'email sont uniques, sinon faux
}
