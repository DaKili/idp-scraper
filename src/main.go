package main

func main() {
	newProjects := make(map[string]Project)
	existingProjects := getProjects()
	ScanERI(&newProjects, &existingProjects)
	ScanFaC(&newProjects, &existingProjects)
	saveProjects(newProjects)
}
