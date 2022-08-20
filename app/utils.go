package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type page struct {
	PageUrl   *url.URL `json:"-"`
	Page      string   `json:"page"`
	ParentUrl *url.URL `json:"-"`
	Parent    string   `json:"parent"`
	Links     []page   `json:"links"`
}

func (p *page) setPageUrl(urlString string, isRoot bool) {
	parsedUrl, err := url.Parse(urlString)

	if err != nil {
		panic(err)
	}

	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "https"
	}
	// TODO fix this
	if parsedUrl.Host == "" {
		newHost, newPath, _ := strings.Cut(parsedUrl.Path, "/")
		parsedUrl.Host = newHost
		if newHost == "" {
			parsedUrl.Host = p.ParentUrl.Host
		}
		parsedUrl.Path = "/" + newPath
	}

	p.PageUrl = parsedUrl
	if isRoot {
		p.ParentUrl, _ = url.Parse("https://" + p.PageUrl.Host)
	}

	p.Page = p.PageUrl.String()
	p.Parent = p.ParentUrl.String()
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
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					child := page{}
					child.ParentUrl = p.PageUrl
					child.setPageUrl(a.Val, false)
					p.Links = append(p.Links, child)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	// TODO check if domain host is same
	if visited[p.PageUrl.String()] == 0 && p.ParentUrl.Host == p.PageUrl.Host {
		doc, err := html.Parse(strings.NewReader(getHtml(p.PageUrl.String())))
		if err != nil {
			panic(err)
		}
		visited[p.PageUrl.String()]++
		f(doc)
	}
}
