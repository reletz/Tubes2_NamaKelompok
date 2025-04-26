package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"backend/util"

	"github.com/PuerkitoBio/goquery"
)

type RecipeJSON struct {
	Pair1  string `json:"pair1"`
	Pair2  string `json:"pair2"`
	Result string `json:"result"`
}

func Scraper(combinations *(map[util.Pair]string), saveToFile bool) {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

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

	lookup := make(map[util.Pair]string)

	doc.Find("table.list-table.col-list.icon-hover tbody tr").Each(func(i int, row *goquery.Selection) {
		tds := row.Find("td")
		if tds.Length() < 2 {
			return
		}

		result := strings.TrimSpace(tds.Eq(0).Find("a").Text())

		tds.Eq(1).Find("li").Each(func(_ int, li *goquery.Selection) {
			parts := []string{}
			li.Find("a").Each(func(_ int, a *goquery.Selection) {
				text := strings.TrimSpace(a.Text())
				if text != "" {
					parts = append(parts, text)
				}
			})

			if len(parts) == 2 && result != "" {
				pair := util.Pair{
					First: parts[0],
					Second: parts[1],
				}

				lookup[pair] = result;
				
				// kalau parts[0] != parts[1], tambahin duplikat kebalikannya
				if parts[0] != parts[1] {
					pair = util.Pair{
						First: parts[1],
						Second: parts[0],
					}

					lookup[pair] = result
				}
			}
		})
	})

	if saveToFile {
		file, err := os.Create("data/recipes.json")
		if err != nil {
			log.Fatal("Failed to create file", err)
		}
		defer file.Close()

		// CONVERT map[Pair]string --> []RecipeJSON
		var recipes []RecipeJSON
		for k, v := range lookup {
			recipes = append(recipes, RecipeJSON{
				Pair1:  k.First,
				Pair2:  k.Second,
				Result: v,
			})
		}

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		if err := enc.Encode(recipes); err != nil {
			log.Fatal("Failed to encode JSON", err)
		}
		fmt.Println("Done. Saved to recipes.json")
	} else {
		fmt.Print("Done. ")
	}

	fmt.Println("Lookup available in memory.")

	*combinations = lookup
}