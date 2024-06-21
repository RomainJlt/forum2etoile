package forum2etoile

import (
	"log"
	
)

// Récupère le dernier ID de livre et retourne le prochain ID disponible
func GetBookLastID() int {
	db := InitDatabase("database/db.db") // Initialise la base de données
	defer db.Close() // Ferme la base de données après utilisation

	var id int

	// Exécute une requête SQL pour obtenir le dernier ID de la table post
	rows, err := db.Query("SELECT IFNULL(MAX(id), 0) AS id FROM post")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // Ferme les lignes après utilisation

	// Récupère l'ID maximum et l'incrémente de 1 pour obtenir le prochain ID
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
	}
	return id + 1
}

// Récupère toutes les données des posts
func GetPostData() []PostData {
	db := InitDatabase("database/db.db")
	defer db.Close() 

	var posts []PostData
	// Exécute une requête SQL pour obtenir toutes les données des posts
	rows, err := db.Query("SELECT id, author, date, title, content, like, dislike, filter, category FROM post")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // Ferme les lignes après utilisation

	// Parcourt toutes les lignes et stocke les données dans la slice posts
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

// Récupère les commentaires d'un post spécifique par ID
func GetCommentData(idInfo int) {
	db := InitDatabase("database/db.db") 
	defer db.Close() 

	var temp Post

	// Exécute une requête SQL pour obtenir les commentaires d'un post par ID
	rows, _ := db.Query("SELECT author, content, date FROM comment WHERE postid = ?", idInfo)
	allResult = nil
	// Parcourt toutes les lignes et stocke les données dans allResult
	for rows.Next() {
		rows.Scan(&temp.AuthorComment, &temp.ContentComment, &temp.DateComment)
		allResult = append(allResult, temp)
	}
}

// Récupère les données d'un post spécifique par ID
func GetPostDataById(idInfo int) {
	db := InitDatabase("database/db.db") 
	defer db.Close()

	var temp PostData

	// Exécute une requête SQL pour obtenir les données d'un post par ID
	rows, _ := db.Query("SELECT id, author, date, title, content, like, dislike, filter FROM post WHERE id = ?", idInfo)
	allData = nil
	// Parcourt toutes les lignes et stocke les données dans allData
	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		allData = append(allData, temp)
	}
}

// Récupère les posts filtrés par un critère spécifique
// Initialise la base de données.
// Ferme la base de données.
func GetPostDataByFilter(filter int) {
	db := InitDatabase("database/db.db") 
	defer db.Close() 

	var temp Post

	// Exécute une requête SQL pour obtenir les posts filtrés
	rows, _ := db.Query("SELECT * FROM post WHERE filter = ?", filter)
	//vide la slice ce que garantit allResult ne contient que les résultats de la requête actuelle.
	allResult = nil
	for rows.Next() {
		// Scanne les valeurs de la ligne courante dans les champs de la structure temp
		rows.Scan(&temp.Id, &temp.Author, &temp.Date, &temp.Title, &temp.Content, &temp.Like, &temp.Dislike, &temp.Filter)
		// Ajoute la structure temp au début de allResult
		allResult = append([]Post{temp}, allResult...)
	}
}

// Récupère les informations d'un utilisateur par pseudo
func GetUserInfo(userInfo string) {
	db := InitDatabase("database/db.db") 
	defer db.Close() 

	var temp Login

	// Exécute une requête SQL pour obtenir les informations d'un utilisateur par pseudo
	rows, _ := db.Query("SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?", userInfo)
	allUser = nil
	// Parcourt toutes les lignes et stocke les données dans allUser
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}

// Récupère les informations d'un utilisateur par cookie (pseudo)
func GetUserInfoByCookie(username string) {
	db := InitDatabase("database/db.db") 
	defer db.Close() 

	var temp Login

	rows, _ := db.Query("SELECT pseudo, email, image, post, subscribers FROM register WHERE pseudo = ?", username)
	allUser = nil
	// Parcourt toutes les lignes et stocke les données dans allUser
	for rows.Next() {
		rows.Scan(&temp.Name, &temp.Email, &temp.Image, &temp.Post, &temp.Sub)
		allUser = append(allUser, temp)
	}
}
