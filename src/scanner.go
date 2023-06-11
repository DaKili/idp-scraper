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
	eriDescriptionSanitizer = regexp.MustCompile(`[\n\t]+|[ ]{2,}|[ ]?\[read more\]`)
)

const EriBaseURL = "https://www.ie.mgt.tum.de/"
const EriProjectPath = "/en/ent/teaching/project-studies-and-idps/"
const FaCBaseURL = "https://www.fa.mgt.tum.de/fm/teaching/idp/"

func ScanFaC(newProjects *map[string]Project, existingProjects *map[string]Project) {
	doc := getDocumentFromURL(FaCBaseURL)

	// FaC doesn't constitently classify their IDPs/ AP.
	// Check each table element and skip the listing numbers.
	// They also don't give the date of publishment or type.
	// They offer to apply with your own topic.
	doc.Find("td").Each(func(index int, item *goquery.Selection) {
		project := Project{Chair: "TUM Financial Management and Capital Markets", School: "TUM School of Management"}
		// Create object.
		text := strings.TrimSpace(item.Text())
		// Extract title, skip entry, if name is too small -> likely just the list number.
		if len(text) < 3 {
			return
		} else {
			project.Title = text
		}

		// Break, if current project already exists.
		if _, ok := (*existingProjects)[project.Title]; ok {
			return
		}

		// FaC can be either IDP or AP.
		project.Type = []projectType{TypeIDP, TypeAP}

		// FaC doesn't give a date.. use first date of extraction. Usually they are gone in like 1-2 weeks anyway and nobody bothers to update the list.
		project.FirstSeen = time.Now()

		// FaC also doen't give a brief description. Maybe something for a text extractor.
		project.Description = ""

		// Extract PDF download link.
		link, exists := item.Find("a").Attr("href")
		if exists {
			project.PdfDownload = getAbsoluteURL(FaCBaseURL, link)
		}

		// Append to list.
		(*newProjects)[text] = project
	})
	// No pagination in FaC.
}

func ScanERI(newProjects *map[string]Project, existingProjects *map[string]Project) {
	doc := getDocumentFromURL(getAbsoluteURL(EriBaseURL, EriProjectPath))
	var lastNew bool = true

	for {
		doc.Find(".article.articletype-3").Each(func(index int, item *goquery.Selection) {
			// Create object.
			project := Project{Chair: "TUM Enterpreneurship Research Institute", School: "TUM School of Management"}
			// Extract title.
			project.Title = strings.TrimSpace(item.Find(".news-header-link").Text())

			// Break, if current project already exists.
			if _, ok := (*existingProjects)[project.Title]; ok {
				lastNew = false
				return
			}

			// Extract type.
			if strings.HasPrefix(project.Title, "IDP") {
				project.Type = []projectType{TypeIDP}
			} else if strings.HasPrefix(project.Title, "Project Study") {
				project.Type = []projectType{TypePS}
			} else {
				project.Type = []projectType{TypeMisc}
			}

			// Extract date.
			dateString := strings.TrimSpace(item.Find("time").Text())
			date, err := time.Parse("02.01.2006", dateString)
			if err != nil {
				log.Fatal(err)
			}
			project.FirstSeen = date

			// Extract description.
			description := eriDescriptionSanitizer.ReplaceAllString(item.Find("p[itemprop=description]").Text(), "")
			project.Description = description

			// Extract PDF download link.
			link, exists := item.Find(".news-header-link").Attr("href")
			if exists {
				subDoc := getDocumentFromURL(getAbsoluteURL(EriBaseURL, link))
				download, subExists := subDoc.Find(".news-related-files-link a").Attr("href")
				if subExists {
					downloadLink := getAbsoluteURL(EriBaseURL, download)
					project.PdfDownload = downloadLink
				}
			}

			// Append to list.
			(*newProjects)[project.Title] = project
		})

		if !lastNew {
			return
		}
		// Check next page.
		nextLink, exists := doc.Find("ul.f3-widget-paginator li.next a").Attr("href")
		if exists {
			fmt.Println(nextLink)
			doc = getDocumentFromURL(getAbsoluteURL(EriBaseURL, nextLink))
		} else {
			break
		}
	}
}

// Wrapper to reduce boiler plate of getting the document representation of a url.
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

func getAbsoluteURL(base string, relativeURL string) string {
	parsedBaseURL, _ := url.Parse(base)
	parsedURL, _ := url.Parse(relativeURL)
	return parsedBaseURL.ResolveReference(parsedURL).String()
}
