func Home(w http.ResponseWriter, r *http.Request) {
	getPostData()
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, allResult)
}

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	db := initDatabase()
	defer db.Close()
	r.ParseForm()
	if r.Method == http.MethodPost {
		pseudo := r.FormValue("pseudo")
		email := r.FormValue("email")
		password := r.FormValue("password")
		hash, _ := HashPassword(password)
		image := "default.jpg"
		_, err := insertIntoRegister(db, pseudo, email, hash, image)
		if err == nil {
			http.Redirect(w, r, "./templates/login.html", http.StatusSeeOther)
		}
	}
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == http.MethodPost {
		pseudo := r.FormValue("pseudo")
		password := r.FormValue("password")
		if login(pseudo, password) {
			session, _ := store.Get(r, "mysession")
			session.Values["username"] = pseudo
			err := session.Save(r, w)
			if err != nil {
				log.Println("Erreur lors de la sauvegarde de la session : ", err)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
		} else {
			log.Println("Identifiants invalides pour l'utilisateur : ", pseudo)
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

}


func Index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"]
	fmt.Println(username)
	if username != nil {
		getUserInfoByCookie(username.(string))
		data := struct {
			User []Login
		}{
			allUser,
		}
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, data)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}


func Profile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	username := fmt.Sprintf("%v", session.Values["username"]) //Decrypts data of the session cookies
	getUserInfoByCookie(username)
	t, _ := template.ParseFiles("templates/profile.html")
	t.Execute(w, allUser)
}

func Post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"]
	if username != nil {
		if r.Method == http.MethodPost {
			title := r.FormValue("title")
			content := r.FormValue("content")
			insertIntoPost(initDatabase(), title, content, username.(string))
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		tmpl := template.Must(template.ParseFiles("templates/newpost.html"))
		tmpl.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}


func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	//Delete and save coockie sessions
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}