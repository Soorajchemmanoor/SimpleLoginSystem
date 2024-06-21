package main

import (
	"fmt"
	"html/template"
	"log"
	"loginsystems/packages/auth"
	"loginsystems/packages/session"
	"loginsystems/packages/validation"
	"net/http"
)

func main() {
	templates := template.Must(template.ParseFiles(
		"theme/login/login.html",
		"theme/app/home.html",
		"theme/login/error404.html",
		"theme/app/about.html",
		"theme/app/services.html",
	))

	loginAssetsHandler := http.StripPrefix("/login/assets/", http.FileServer(http.Dir("theme/login/assets/")))
	http.Handle("/login/assets/", loginAssetsHandler)

	appAssetsHandler := http.StripPrefix("/app/assets/", http.FileServer(http.Dir("theme/app/assets/")))
	http.Handle("/app/assets/", appAssetsHandler)

	http.HandleFunc("/login", loginHandler(templates))
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/", rootHandler(templates))
	http.HandleFunc("/about-us", aboutHandler(templates))
	http.HandleFunc("/services", servicesHandler(templates))
	http.HandleFunc("/404", error404Handler(templates))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func rootHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			error404Handler(templates)(w, r)
			return
		}

		username, err := session.GetSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")

		data := map[string]interface{}{
			"username": username,
		}
		if err := templates.ExecuteTemplate(w, "home", data); err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func aboutHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := session.GetSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")

		data := map[string]interface{}{
			"username": username,
		}
		if err := templates.ExecuteTemplate(w, "about", data); err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func servicesHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := session.GetSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")

		data := map[string]interface{}{
			"username": username,
		}
		if err := templates.ExecuteTemplate(w, "services", data); err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func loginHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := session.GetSession(r); err == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")

		data := map[string]string{
			"username": "",
			"password": "",
			"uErr":     "",
			"pErr":     "",
			"lErr":     "",
		}

		if r.Method == "POST" {
			username, password := r.PostFormValue("username"), r.PostFormValue("password")
			data["username"] = username
			data["password"] = password

			if uErr, pErr := validation.ValidateForm(username, password); uErr == "" && pErr == "" {
				if auth.AuthLogin(username, password) {
					session.SetSession(w, username)
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				} else {
					data["lErr"] = "Invalid username or password"
				}
			} else {
				data["uErr"] = uErr
				data["pErr"] = pErr
			}
		}

		if err := templates.ExecuteTemplate(w, "login", data); err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session.ClearSession(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func error404Handler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		err := templates.ExecuteTemplate(w, "error404", nil)
		if err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
