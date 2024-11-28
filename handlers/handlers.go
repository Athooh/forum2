package handlers

import (
    "fmt"
    "net/http"
    "text/template"
)

var templates *template.Template

func init() {
    // Parse all templates in one go
    templates = template.Must(template.ParseFiles(
        "templates/base.html",
        "templates/login.html",
        "templates/signup.html",
    ))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Redirect to login if accessing root directly
    if r.URL.Path == "/" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("LoginHandler invoked")

    data := map[string]interface{}{
        "Title": "Login - Forum",
    }

    err := templates.ExecuteTemplate(w, "base", data)
    if err != nil {
        fmt.Printf("Error executing login template: %v\n", err)
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
        return
    }
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("SignupHandler invoked")

    data := map[string]interface{}{
        "Title": "Sign Up - Forum",
    }

    err := templates.ExecuteTemplate(w, "base", data)
    if err != nil {
        fmt.Printf("Error executing signup template: %v\n", err)
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
        return
    }
}