package forum2etoile

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func GetCategories(db *sql.DB) []Category {
	// Exécute une requête SQL pour sélectionner les colonnes id et name de la table category
	rows, err := db.Query("SELECT id, name FROM category")
	if err != nil {
		log.Fatal(err)
	}
	// Ferme les lignes après que tout soit terminé
	defer rows.Close() 

	var categories []Category
	for rows.Next() { // Itère sur chaque ligne résultante
		var category Category
		// Scanne les valeurs de la ligne courante dans les champs de la structure Category
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			log.Fatal(err)
		}
		// Ajoute la catégorie à la slice de catégories
		categories = append(categories, category)
	}
	return categories // Retourne la slice de catégories
}
