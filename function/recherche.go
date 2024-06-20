package forum2etoile

import (
    "log"
    "net/http"
    "html/template"

    _ "github.com/mattn/go-sqlite3"

)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.FormValue("q")
    if query == "" {
        http.Redirect(w, r, "/index", http.StatusSeeOther)
        return
    }
    results := SearchPosts(query)
    t, err := template.ParseFiles("search.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, results)
}

func SearchPosts(query string) []PostData {
    db := InitDatabase("database/db.db")
    defer db.Close()

    query = "%" + query + "%"
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
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var post PostData
        err := rows.Scan(&post.Id, &post.Author, &post.Date, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.Filter, &post.Category)
        if err != nil {
            log.Fatal(err)
        }
        results = append(results, post)
    }
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }

    return results
}