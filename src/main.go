package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var projects = make(map[string]Project)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Loading index.html")
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, projects)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Loaded index.html")
}

func main() {
	fmt.Println("Loading projects")
	projects = getProjects()
	http.HandleFunc("/", HandleIndex)
	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8000", nil))
	// newProjects := make(map[string]Project)
	// existingProjects := getProjects()
	// ScanERI(&newProjects, &existingProjects)
	// ScanFaC(&newProjects, &existingProjects)
	// saveProjects(newProjects)
}
