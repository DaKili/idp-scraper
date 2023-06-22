package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
)

var projects Projects

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Loading index.html")
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	projectArray := projects
	x := map[string][]Project{
		"Projects": projectArray,
	}
	err = tmpl.Execute(w, x)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Loaded index.html")
}

func HandleSortProjectsByDate(w http.ResponseWriter, r *http.Request) {
	log.Println("Sorting table by date")

	sortedProjects := projects
	sort.Stable(DateSorter(sortedProjects))
	tmpl, err := template.ParseFiles("templates/table_rows.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sortedProjectsMap := map[string][]Project{
		"Projects": sortedProjects,
	}
	err = tmpl.Execute(w, sortedProjectsMap)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Loaded new table body sorted by date")
}

func HandleClearDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Clearing database")
	delProjects()
	projects = Projects{}

	tmpl, err := template.ParseFiles("templates/table_rows.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	projectsMap := map[string][]Project{
		"Projects": projects,
	}

	err = tmpl.Execute(w, projectsMap)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("Cleared database")
}

func HandleUpdateDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating database")
	var newProjects Projects
	ScanERI(&newProjects, &projects)
	saveProjects(&newProjects)
	projects = append(projects[:], newProjects...)

	tmpl, err := template.ParseFiles("templates/table_rows.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	projectsMap := map[string][]Project{
		"Projects": projects,
	}
	err = tmpl.Execute(w, projectsMap)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Updated table")
}

func main() {
	fmt.Println("Loading projects")
	projects = getProjects()
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/sort_date/", HandleSortProjectsByDate)
	http.HandleFunc("/clear_database/", HandleClearDatabase)
	http.HandleFunc("/update_database/", HandleUpdateDatabase)
	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8000", nil))
	// newProjects := make(map[string]Project)
	// existingProjects := getProjects()
	// ScanERI(&newProjects, &existingProjects)
	// ScanFaC(&newProjects, &existingProjects)
	// saveProjects(newProjects)
}
