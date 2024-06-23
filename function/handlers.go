package forum2etoile

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
	_ "github.com/mattn/go-sqlite3"
)
    // affichage de la page d'accueil.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    posts := GetPostData()
    t, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, posts)
}

    // affichage de la page de connexion.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    // Récupère le cookie contenant le nom d'utilisateur.
    pseudoForm := r.FormValue("pseudoCreate")
    emailForm := r.FormValue("emailCreate")
    passwordForm := r.FormValue("passwordCreate")
    imageForm := r.FormValue("imageCreate")
    pseudoLog := r.FormValue("pseudoLog")
    passwordLog := r.FormValue("passwordLog")

    db := InitDatabase("database/db.db")
    if r.Method == "POST" && r.FormValue("deleteAccount") == "true" {
        err := DeleteAccount(db, pseudoLog)
    // si la suppression échoue, retourne une erreur.
    // autrement redirige vers la page d'accueil.
        if err != nil {
            http.Error(w, "Failed to delete account", http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/index", http.StatusSeeOther)
        return
    }
    // si le formulaire de création de compte est soumis.
    //vérifie les champs et insère les données dans la base de données.
    hash, _ := HashPassword(passwordForm)
    if pseudoForm != "" && emailForm != "" && passwordForm != "" {
        if Register2(pseudoForm, emailForm) {
            if imageForm != "" {
                InsertIntoRegister(db, pseudoForm, emailForm, hash, imageForm)
            } else {
                InsertIntoRegister(db, pseudoForm, emailForm, hash, "http://marclimoservices.com/wp-content/uploads/2017/05/facebook-default.png")
            }
        }
    }
    // verification si le formulaire de connexion est soumis.
    if Login2(pseudoLog, passwordLog) {
        user.Name = pseudoLog
        expiration := time.Now().Add(24 * time.Hour)
    // création d'un cookie pour stocker le nom d'utilisateur.
        cookie := http.Cookie{Name: "username", Value: pseudoLog, Expires: expiration}
    // définit le cookie dans la réponse.
        http.SetCookie(w, &cookie)
        http.Redirect(w, r, "/index", http.StatusSeeOther)
        return
    }

    t, _ := template.ParseFiles("templates/register.html")
    t.Execute(w, nil)
}
    // affichage de la page de connexion.
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	GetUserInfoByCookie(username)
	t, _ := template.ParseFiles("templates/profile.html")
	t.Execute(w, allUser)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
    // Récupère le pseudo de l'utilisateur depuis l'URL
    userInfo := r.URL.Path[len("/user/"):]
    GetUserInfo(userInfo)
    t, err := template.ParseFiles("templates/user.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, allUser)
}
    // affichage de la page logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "username", Value: "", Expires: time.Unix(0, 0), MaxAge: -1}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

    // affichage de la page de post
func PostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
    // Redirige vers la page de connexion si le cookie n'existe pas
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	db := InitDatabase("database/db.db")
	titleForm := r.FormValue("inputTitle")
	contentForm := r.FormValue("inputContent")
	categoryForm := r.FormValue("category")
    // Récupère les informations de l'utilisateur connecté
	user.Post = 0
    // Vérifie si le formulaire de post est soumis
	if titleForm != "" && contentForm != "" && categoryForm != "" {
		category := categoryForm
	// Vérifie si la catégorie existe dans la base de données
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		_, err = InsertIntoPost(db, titleForm, contentForm, username, category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
    // Incrémente le nombre de posts de l'utilisateur
		db.Exec(`UPDATE register SET post = post + 1 WHERE pseudo = ?`, username)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
		return
	}

	t, _ := template.ParseFiles("templates/post.html")
	categories := GetCategories(db)
	t.Execute(w, map[string]interface{}{
		"Categories": categories,
	})
}
// affiche la page info.
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
    
	username := cookie.Value
	db := InitDatabase("database/db.db/")

	idInfo, _ := strconv.Atoi(r.URL.Path[6:])
	contentComment := r.FormValue("commentArea")
	redirect := "/info/" + strconv.Itoa(idInfo)

	GetPostDataById(idInfo)
// Vérifie si le formulaire de commentaire est soumis.
	if len(contentComment) > 0 {
		InsertIntoComment(db, idInfo, username, contentComment)
        // Met à jour la date du post pour le trier par date de commentaire
		db.Exec(`UPDATE comment SET date = ? WHERE postid = ?`, time.Now(), idInfo)
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}

	GetCommentData(idInfo)
// Crée une carte pour stocker les résultats de la requête.
// et les données du post pour les afficher dans la page info.
	m := map[string]interface{}{
		"Results": allResult,
		"Post":    allData,
	}
	t := template.Must(template.ParseFiles("templates/info.html"))
	t.Execute(w, m)
}
// affiche la page de mise à jour de profil.
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("username")
    if err != nil {
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }
    // recupère le nom d'utilisateur du cookie
    username := cookie.Value
    db := InitDatabase("database/db.db")
    defer db.Close()
    // Vérifie si le formulaire de mise à jour est soumis
    if r.Method == http.MethodPost {
        newPassword := r.FormValue("newPassword")
        newPseudo := r.FormValue("newPseudo")

        if newPassword != "" {
            hash, _ := HashPassword(newPassword)
            _, err := db.Exec(`UPDATE register SET password = ? WHERE pseudo = ?`, hash, username)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }
    // Vérifie si le nouveau pseudo est différent du pseudo actuel.
        if newPseudo != "" {
            _, err := db.Exec(`UPDATE register SET pseudo = ? WHERE pseudo = ?`, newPseudo, username)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            // Met à jour le cookie avec le nouveau pseudo
            cookie := http.Cookie{Name: "username", Value: newPseudo, Expires: time.Now().Add(24 * time.Hour)}
            http.SetCookie(w, &cookie)
        }

        http.Redirect(w, r, "/profile", http.StatusSeeOther)
        return
    }

    // If the method is not POST, display the profile page with the update form
    t, err := template.ParseFiles("templates/update_profile.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, nil)
}

func RegulationsHandler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("templates/regulations.html")
    if err != nil {
        http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
        return
    }
    err = t.Execute(w, nil)
    if err != nil { 
        http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
    }
}