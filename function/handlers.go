package forum2etoile

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
	_ "github.com/mattn/go-sqlite3"
)
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    posts := GetPostData()
    t, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, posts)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    pseudoForm := r.FormValue("pseudoCreate")
    emailForm := r.FormValue("emailCreate")
    passwordForm := r.FormValue("passwordCreate")
    imageForm := r.FormValue("imageCreate")
    pseudoLog := r.FormValue("pseudoLog")
    passwordLog := r.FormValue("passwordLog")

    db := InitDatabase("database/db.db")

    if r.Method == "POST" && r.FormValue("deleteAccount") == "true" {
        err := DeleteAccount(db, pseudoLog)
        if err != nil {
            http.Error(w, "Failed to delete account", http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/index", http.StatusSeeOther)
        return
    }

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

    if Login2(pseudoLog, passwordLog) {
        user.Name = pseudoLog
        expiration := time.Now().Add(24 * time.Hour)
        cookie := http.Cookie{Name: "username", Value: pseudoLog, Expires: expiration}
        http.SetCookie(w, &cookie)
        http.Redirect(w, r, "/index", http.StatusSeeOther)
        return
    }

    t, _ := template.ParseFiles("register.html")
    t.Execute(w, nil)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	GetUserInfoByCookie(username)
	t, _ := template.ParseFiles("profile.html")
	t.Execute(w, allUser)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
    userInfo := r.URL.Path[len("/user/"):]
    GetUserInfo(userInfo)
    t, err := template.ParseFiles("user.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, allUser)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "username", Value: "", Expires: time.Unix(0, 0), MaxAge: -1}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}


func PostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	db := InitDatabase("database/db.db")
	titleForm := r.FormValue("inputTitle")
	contentForm := r.FormValue("inputContent")
	categoryForm := r.FormValue("category")
	user.Post = 0

	if titleForm != "" && contentForm != "" && categoryForm != "" {
		category := categoryForm
		
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		_, err = InsertIntoPost(db, titleForm, contentForm, username, category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db.Exec(`UPDATE register SET post = post + 1 WHERE pseudo = ?`, username)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
		return
	}

	t, _ := template.ParseFiles("post.html")
	categories := GetCategories(db)
	t.Execute(w, map[string]interface{}{
		"Categories": categories,
	})
}

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

	if len(contentComment) > 0 {
		InsertIntoComment(db, idInfo, username, contentComment)
		db.Exec(`UPDATE comment SET date = ? WHERE postid = ?`, time.Now(), idInfo)
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}

	GetCommentData(idInfo)

	m := map[string]interface{}{
		"Results": allResult,
		"Post":    allData,
	}
	t := template.Must(template.ParseFiles("info.html"))
	t.Execute(w, m)
}


func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("username")
    if err != nil {
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }
    username := cookie.Value
    db := InitDatabase("database/db.db")
    defer db.Close()

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

        if newPseudo != "" {
            _, err := db.Exec(`UPDATE register SET pseudo = ? WHERE pseudo = ?`, newPseudo, username)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            // Update the cookie with the new pseudo
            cookie := http.Cookie{Name: "username", Value: newPseudo, Expires: time.Now().Add(24 * time.Hour)}
            http.SetCookie(w, &cookie)
        }

        http.Redirect(w, r, "/profile", http.StatusSeeOther)
        return
    }

    // If the method is not POST, display the profile page with the update form
    t, err := template.ParseFiles("update_profile.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, nil)
}