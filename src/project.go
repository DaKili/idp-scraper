package main

import "time"

type projectType string

const (
	TypeIDP  projectType = "IDP"
	TypePS   projectType = "Project Study"
	TypeAP   projectType = "Application Project"
	TypeMisc projectType = "Undefined"
)

type Project struct {
	Chair       string
	School      string
	Description string
	Title       string
	PdfDownload string
	Type        []projectType
	FirstSeen   time.Time
}

type Projects []Project

func (projs *Projects) Append(project Project) {
	*projs = append(*projs, project)
}

func (projs Projects) Contains(proj Project) bool {
	y1, m1, d1 := proj.FirstSeen.Date()
	for _, v := range projs {
		y2, m2, d2 := v.FirstSeen.Date()
		if v.Title == proj.Title && y1 == y2 && m1 == m2 && d1 == d2 {
			return true
		}
	}
	return false
}

// Sort by date
type DateSorter []Project

func (a DateSorter) Len() int           { return len(a) }
func (a DateSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DateSorter) Less(i, j int) bool { return a[i].FirstSeen.Unix() < a[j].FirstSeen.Unix() }
