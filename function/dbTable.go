package forum2etoile

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	
)

// InitDatabase initialise la base de données SQLite avec les tables nécessaires
// et retourne l'objet *sql.DB représentant la connexion à la base de données.
func InitDatabase(database string) *sql.DB {
	// Ouvre la connexion à la base de données SQLite spécifiée
	db, err := sql.Open("sqlite3", "database/db.db")
	if err != nil {
		// Arrête le programme si une erreur survient lors de l'ouverture de la base de données
		log.Fatal(err)
	}
	// Créer Sql tables si elles n'existent pas déjà
	sqlStmt := `
				CREATE TABLE IF NOT EXISTS register (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					pseudo TEXT NOT NULL, 
					email TEXT NOT NULL, 
					password TEXT NOT NULL,
					image TEXT NOT NULL,
					post INT NOT NULL,
					subscribers INT NOT NULL
				);

				CREATE TABLE IF NOT EXISTS post (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					author TEXT NOT NULL,
					date TEXT NOT NULL,
					title TEXT NOT NULL,
					content TEXT NOT NULL,
					like INT NOT NULL,
					dislike INT NOT NULL,
					filter INT NOT NULL,
					category TEXT NOT NULL
					
				);

				CREATE TABLE IF NOT EXISTS comment (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					postid INT NOT NULL,
					date TEXT NOT NULL,
					author TEXT NOT NULL,
					content TEXT NOT NULL
				);

				CREATE TABLE IF NOT EXISTS like (
					postid INTEGER NOT NULL,
					author TEXT NOT NULL,
					like INT NOT NULL,
					dislike INT NOT NULL,
					PRIMARY KEY (postid, author)
				);

					 CREATE TABLE IF NOT EXISTS category (
   					 id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
   					 name TEXT NOT NULL
				);
				`
				// Exécute les instructions SQL pour créer les tables dans la base de données
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	return db
	
}
