package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type page struct {
	Visited   map[string]int `json:"visited"`
	PageUrl   *url.URL       `json:"pageUrl"`
	ParentUrl *url.URL       `json:"parentUrl"`
	Links     []page         `json:"links"`
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
		parsedUrl.Host = newHost
		if newHost == "" {
			parsedUrl.Host = p.ParentUrl.Host
		}
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
	For the given htmlDOM string, returns a map of type map[string]int,

where key is the url and value is the number of times it was found in the DOM
*/
func (p *page) getLinks() {
	var links = make(map[string]int)
	doc, err := html.Parse(strings.NewReader(getHtml(p.PageUrl.String())))
	if err != nil {
		panic(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					child := page{}
					child.ParentUrl = p.PageUrl
					child.setPageUrl(a.Val)
					p.Links = append(p.Links, child)
					links[a.Val]++
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	p.Visited = links

	// if len(p.Links) > 0 {
	// 	for i := 0; i < len(p.Links); i++ {
	// 		p.Links[i].getLinks()
	// 	}
	// }
}
