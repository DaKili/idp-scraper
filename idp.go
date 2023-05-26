package main

import "time"

type projectType int

const (
	TypeIDP  = "IDP"
	TypePS   = "Project Study"
	TypeMisc = "Undefined"
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
