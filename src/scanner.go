package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	descriptionSanitizer = regexp.MustCompile(`[\n\t]+|[ ]{2,}|[ ]?\[read more\]`)
)

const EriBaseURL = "https://www.ie.mgt.tum.de/"
const EriProjectPath = "/en/ent/teaching/project-studies-and-idps/"

func ScanERI(projects *map[string]Project, existingProjects *map[string]Project) {
	doc := getDocumentFromURL(getAbsoluteURL(EriProjectPath))
	var lastNew bool = true

	for {
		doc.Find(".article.articletype-3").Each(func(index int, item *goquery.Selection) {
			// Create object
			project := Project{Chair: "TUM Enterpreneurship Research Institute", School: "TUM School of Management"}
			// Extract title
			project.Title = strings.TrimSpace(item.Find(".news-header-link").Text())

			// break, if current project already exists
			if _, ok := (*existingProjects)[project.Title]; ok {
				lastNew = false
				return
			}

			// Extract type
			if strings.HasPrefix(project.Title, "IDP") {
				project.Type = TypeIDP
			} else if strings.HasPrefix(project.Title, "Project Study") {
				project.Type = TypePS
			} else {
				project.Type = TypeMisc
			}

			// Extract date
			dateString := strings.TrimSpace(item.Find("time").Text())
			date, err := time.Parse("02.01.2006", dateString)
			if err != nil {
				log.Fatal(err)
			}
			project.Date = date

			// Extract description
			description := descriptionSanitizer.ReplaceAllString(item.Find("p[itemprop=description]").Text(), "")
			project.Description = description

			// Extract download link
			link, exists := item.Find(".news-header-link").Attr("href")
			if exists {
				subDoc := getDocumentFromURL(getAbsoluteURL(link))
				download, subExists := subDoc.Find(".news-related-files-link a").Attr("href")
				if subExists {
					downloadLink := getAbsoluteURL(download)
					project.PdfDownload = downloadLink
				}
			}

			// Append to list
			(*projects)[project.Title] = project
		})

		if !lastNew {
			return
		}
		// Check next page
		nextLink, exists := doc.Find("ul.f3-widget-paginator li.next a").Attr("href")
		if exists {
			fmt.Println(nextLink)
			doc = getDocumentFromURL(getAbsoluteURL(nextLink))
		} else {
			break
		}
	}
}

func getDocumentFromURL(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatal("Statuscode error:", res.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func getAbsoluteURL(relativeURL string) string {
	parsedBaseURL, _ := url.Parse(EriBaseURL)
	parsedURL, _ := url.Parse(relativeURL)
	return parsedBaseURL.ResolveReference(parsedURL).String()
}
