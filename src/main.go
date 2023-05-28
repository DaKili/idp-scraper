package main

func main() {
	newProjects := make(map[string]Project)
	existingProjects := getProjects()
	ScanERI(&newProjects, &existingProjects)
	saveProjects(newProjects)
}
