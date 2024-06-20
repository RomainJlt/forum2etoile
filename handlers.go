func indexHandler(w http.ResponseWriter, r *http.Request) {
    posts := getPostData()
    t, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, posts)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    pseudoForm := r.FormValue("pseudoCreate")
    emailForm := r.FormValue("emailCreate")
    passwordForm := r.FormValue("passwordCreate")
    imageForm := r.FormValue("imageCreate")
    pseudoLog := r.FormValue("pseudoLog")
    passwordLog := r.FormValue("passwordLog")

    db := initDatabase("database/db.db")

    if r.Method == "POST" && r.FormValue("deleteAccount") == "true" {
        err := deleteAccount(db, pseudoLog)
        if err != nil {
            http.Error(w, "Failed to delete account", http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/index", http.StatusSeeOther)
        return
    }

    hash, _ := HashPassword(passwordForm)
    if pseudoForm != "" && emailForm != "" && passwordForm != "" {
        if register(pseudoForm, emailForm) {
            if imageForm != "" {
                insertIntoRegister(db, pseudoForm, emailForm, hash, imageForm)
            } else {
                insertIntoRegister(db, pseudoForm, emailForm, hash, "http://marclimoservices.com/wp-content/uploads/2017/05/facebook-default.png")
            }
        }
    }

    if login(pseudoLog, passwordLog) {
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

func profileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	getUserInfoByCookie(username)
	t, _ := template.ParseFiles("profile.html")
	t.Execute(w, allUser)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    userInfo := r.URL.Path[len("/user/"):]
    getUserInfo(userInfo)
    t, err := template.ParseFiles("user.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, allUser)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "username", Value: "", Expires: time.Unix(0, 0), MaxAge: -1}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
	likeId := r.URL.Path[6:]
	redirect := "/info/" + likeId
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	checkLike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func dislikeHandler(w http.ResponseWriter, r *http.Request) {
	likeId := r.URL.Path[9:]
	redirect := "/info/" + likeId
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	checkDislike(username, likeId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	db := initDatabase("database/db.db")
	titleForm := r.FormValue("inputTitle")
	contentForm := r.FormValue("inputContent")
	categoryForm := r.FormValue("category")
	user.Post = 0

	if titleForm != "" && contentForm != "" && categoryForm != "" {
		category := categoryForm
		fmt.Println(category)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		_, err = insertIntoPost(db, titleForm, contentForm, username, category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db.Exec(`UPDATE register SET post = post + 1 WHERE pseudo = ?`, username)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
		return
	}

	t, _ := template.ParseFiles("post.html")
	categories := getCategories(db)
	t.Execute(w, map[string]interface{}{
		"Categories": categories,
	})
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value
	db := initDatabase("database/db.db/")
	idInfo, _ := strconv.Atoi(r.URL.Path[6:])
	contentComment := r.FormValue("commentArea")
	redirect := "/info/" + strconv.Itoa(idInfo)

	getPostDataById(idInfo)

	if len(contentComment) > 0 {
		insertIntoComment(db, idInfo, username, contentComment)
		db.Exec(`UPDATE comment SET date = ? WHERE postid = ?`, time.Now(), idInfo)
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}

	getCommentData(idInfo)

	m := map[string]interface{}{
		"Results": allResult,
		"Post":    allData,
	}
	t := template.Must(template.ParseFiles("info.html"))
	t.Execute(w, m)
}


func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("username")
    if err != nil {
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }
    username := cookie.Value
    db := initDatabase("database/db.db")
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