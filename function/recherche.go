package forum2etoile

import (
    "log"
    "net/http"
    "html/template"

    _ "github.com/mattn/go-sqlite3"

)
// SearchHandler gère les requêtes de recherche de posts.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
    // Récupère la valeur de la requête de recherche.
    //cette valeur est stockée dans le paramètre "q" de la requête GET.
    query := r.FormValue("q")
    if query == "" {
        // Redirige vers la page d'accueil si la requête est vide.
        http.Redirect(w, r, "/index", http.StatusSeeOther)
        return
    }
    // Recherche les posts correspondant à la requêtevet appel la function searchPost.
    results := SearchPosts(query)
    // Affiche les résultats de la recherche dans le modèle de recherche.
    t, err := template.ParseFiles("templates/search.html")
    // Retourne une erreur si le modèle n'est pas trouvé.
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Affiche les résultats de la recherche.
    t.Execute(w, results)
}

// SearchPosts recherche les posts correspondant à la requête spécifiée.
func SearchPosts(query string) []PostData {
    db := InitDatabase("database/db.db")
    defer db.Close()
    // Ajoute des caractères de joker (%) autour de la requête pour rechercher des correspondances partielles.
    query = "%" + query + "%"
    // Requête SQL pour rechercher des posts basés sur le titre, le contenu ...
    var results []PostData
    rows, err := db.Query(`SELECT id, author, date, title, content, like, dislike, filter, category 
                           FROM post 
                           WHERE lower(title) LIKE lower(?) 
                           OR lower(content) LIKE lower(?) 
                           OR lower(date) LIKE lower(?) 
                           OR like = ? 
                           OR dislike = ? 
                           OR lower(author) LIKE lower(?) 
                           OR lower(category) LIKE lower(?)`, 
                           query, query, query, query, query, query, query)
    // Retourne une erreur si la requête échoue.
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
// Parcourt les résultats de la requête et stocke les données dans la slice results.
    for rows.Next() {
        var post PostData
        err := rows.Scan(&post.Id, &post.Author, &post.Date, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.Filter, &post.Category)
        if err != nil {
            log.Fatal(err)
        }
        results = append(results, post)
    }
    // Retourne une erreur si la requête échoue.
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }

    return results
}