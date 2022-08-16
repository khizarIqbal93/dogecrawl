package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	webPage := page{}
	webPage.setPageUrl("https://www.devopsplayground.co.uk/events/register")
	webPage.getLinks()
	pageJson, _ := json.Marshal(webPage)
	writeFileErr := os.WriteFile("output.json", pageJson, 0666)
	if writeFileErr != nil {
		log.Fatal("Could not create output file")
	}
}
