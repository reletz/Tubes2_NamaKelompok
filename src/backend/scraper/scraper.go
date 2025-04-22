package scraper

import (
	"fmt"
	"net/http"
	"log"
	"strings"
	"os"
	"encoding/json"
	
	"github.com/PuerkitoBio/goquery"
)

type Recipe struct {
	Result     string   `json:"result"`
	Components []string `json:"components"`
}

func Scraper() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

	// Ambil HTML-nya
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to connect to the target page", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("HTTP Error %d: %s", resp.StatusCode, resp.Status)
 	}

  
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Failed to parse the HTML document", err)
	}

	var recipes []Recipe;

	// Each iterates over a Selection object, executing a function for each matched element. 
	// It returns the current Selection object. 
	// The function f is called for each element in the selection with the index of the element in that selection starting at 0, 
	// and a *Selection that contains only that element.
	doc.Find("table.list-table.col-list.icon-hover tbody tr").Each(func(i int, row *goquery.Selection) {
		tds := row.Find("td")
		if tds.Length() < 2 {
			return
		}
		
		// Eq reduces the set of matched elements to the one at the specified index. 
		// If a negative index is given, it counts backwards starting at the end of the set. 
		// It returns a new Selection object, and an empty Selection object if the index is invalid.

		// Eq(0) = ambil kolom 1 (0 kalo zero-indexing)
		result := strings.TrimSpace(tds.Eq(0).Find("a").Text())

		// Exclude element dasar dari recipe
		if result == "Fire" || result == "Water" || result == "Air" {
			result = "";
		}
	
		// Eq(1) = ambil kolom 0 (1 kalo zero-indexing)
		tds.Eq(1).Find("li").Each(func(_ int, li *goquery.Selection) {
			parts := []string{}
			li.Find("a").Each(func(_ int, a *goquery.Selection) {
				text := strings.TrimSpace(a.Text())
				if text != "" {
					parts = append(parts, text)
				}
			})
		
			if len(parts) == 2 && result != "" {
				r := Recipe{
					Result:     strings.TrimSpace(result),
					Components: []string{
						strings.TrimSpace(parts[0]),
						strings.TrimSpace(parts[1]),
					},
				}
				recipes = append(recipes, r)
			}
		})		
	})

	file, _ := os.Create("data/recipes.json");
	defer file.Close();

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(recipes)

	fmt.Println("Done. Saved to recipes.json")
}