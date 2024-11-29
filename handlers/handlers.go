package handlers

import (
	"fmt"
	"net/http"
	"text/template"
)

var templates *template.Template

func init() {
	// Parse all templates in one go
	templates = template.Must(template.ParseGlob("templates/*.html"))

	// Print actual template names
	for _, t := range templates.Templates() {
		fmt.Println("Parsed template name:", t.Name())
	}
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

	// List all available templates
	// fmt.Println("Available templates:", templates.Templates())

	data := map[string]interface{}{
		"Title": "Login - Forum",
	}

	err := templates.ExecuteTemplate(w, "login", data)
	if err != nil {
		fmt.Printf("Error executing login template: %v\n", err)
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SignupHandler invoked")

	// List all available templates
	// fmt.Println("Available templates:", templates.Templates())

	data := map[string]interface{}{
		"Title": "Login - Forum",
	}

	err := templates.ExecuteTemplate(w, "signup", data)
	if err != nil {
		fmt.Printf("Error executing signup template: %v\n", err)
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SignupHandler invoked")

	// List all available templates
	// fmt.Println("Available templates:", templates.Templates())

	data := map[string]interface{}{
		"Title": "Dashboard - Forum",
	}

	err := templates.ExecuteTemplate(w, "dashboard", data)
	if err != nil {
		fmt.Printf("Error executing signup template: %v\n", err)
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}