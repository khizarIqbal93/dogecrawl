package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type page struct {
	PageUrl   *url.URL `json:"pageUrl"`
	ParentUrl *url.URL `json:"parentUrl"`
	Links     []page   `json:"links"`
}

func (p *page) setPageUrl(urlString string) {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "https"
	}
	if parsedUrl.Host == "" {
		newHost, newPath, _ := strings.Cut(parsedUrl.Path, "/")
		if newHost == "" {
			parsedUrl.Host = p.ParentUrl.Host
		}
		parsedUrl.Host = newHost
		parsedUrl.Path = "/" + newPath
	}

	p.PageUrl = parsedUrl
}

// returns the DOM of urlString as a string
func getHtml(urlString string) string {
	resp, err := http.Get(urlString)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	html := string(body)
	return html
}

/*
appends links found in p.PageUrl under p.Links and records it in visited map[string]int
*/
func (p *page) getLinks(visited map[string]int) {
	// TODO check if domain same
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					child := page{}
					child.ParentUrl = p.PageUrl
					child.setPageUrl(a.Val)
					p.Links = append(p.Links, child)
					visited[child.PageUrl.String()]++
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	if visited[p.PageUrl.String()] == 0 {
		doc, err := html.Parse(strings.NewReader(getHtml(p.PageUrl.String())))
		if err != nil {
			panic(err)
		}

		f(doc)
	}
}
