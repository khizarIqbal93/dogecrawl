package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var visited map[string]int = make(map[string]int)

func main() {
	webPage := page{}
	webPage.setPageUrl("https://www.devopsplayground.co.uk/events/register")
	webPage.getLinks(visited)
	pageJson, _ := json.Marshal(webPage)
	writeFileErr := os.WriteFile("output.json", pageJson, 0666)
	if writeFileErr != nil {
		log.Fatal("Could not create output file")
	}
	fmt.Println(visited)
}
