package main

func main() {

	projects := make(map[string]Project)
	existingProjects := getProjects()
	ScanERI(&projects, &existingProjects)
	saveProjects(projects)
}
