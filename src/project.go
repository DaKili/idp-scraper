package main

import "time"

type projectType string

const (
	TypeIDP  projectType = "IDP"
	TypePS   projectType = "Project Study"
	TypeMisc projectType = "Undefined"
)

type Project struct {
	Chair       string
	School      string
	Description string
	Title       string
	PdfDownload string
	Type        projectType
	Date        time.Time
}
