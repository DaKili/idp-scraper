package main

func main() {
	// Todo: make scanner only browse pages, if no old idp is found on the current page
	// Todo: bug trying to upload an empty dictionary
	projects := make(map[string]Project)
	existingProjects := getProjects()
	ScanERI(&projects, &existingProjects)
	saveProjects(projects)
}
