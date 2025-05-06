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
	Result       string     `json:"Result"`
	Asset				 string			`json:"Asset"`
	Combinations util.Pair  `json:"Combinations"`
}

func Scraper(
	combinations map[util.Pair][]string,
	revCombinations map[string][]util.Pair,
	saveToFile bool,
) {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to connect to the target page:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("HTTP Error %d: %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Failed to parse the HTML document:", err)
	}

	assetMap := make(map[string]string)

	// Parse the table and gather combinations
	doc.Find("table.list-table.col-list.icon-hover tbody tr").Each(func(_ int, row *goquery.Selection) {
		cols := row.Find("td")
		if cols.Length() < 2 {
			return
		}
		result := strings.TrimSpace(cols.Eq(0).Find("a").Text())
		if result == "" || result == "Time" {
			return
		}

		hrefVal, exists := cols.Eq(0).Find("span a").Attr("href")
		asset := ""
		if exists {
			asset = strings.TrimSpace(hrefVal)
		}
		
		if asset == "" {
			asset = "https://static.wikia.nocookie.net/little-alchemy/images/6/63/Time_2.svg/revision/latest?cb=20210827124225"
		}

		assetMap[result] = asset

		// Process combinations from the second column
		cols.Eq(1).Find("li").Each(func(_ int, li *goquery.Selection) {
			parts := []string{}
			li.Find("a").Each(func(_ int, a *goquery.Selection) {
				txt := strings.TrimSpace(a.Text())
				if (txt != "" && txt != "Time") {
					parts = append(parts, txt)
				}
			})
			// Ensure we have two parts for a valid combination
			if len(parts) != 2 {
				return
			}

			pair := util.Pair{
				First: parts[0],
				Second: parts[1],
			}

			combinations[pair] = append(combinations[pair], result)
			revCombinations[result] = append(revCombinations[result], pair)

			// Add both directions
			if parts[0] != parts[1] {
				reversedPair := util.Pair{
					First: parts[1],
					Second: parts[0],
				}
				combinations[reversedPair] = append(combinations[reversedPair], result)
				revCombinations[result] = append(revCombinations[result], reversedPair)
			}
		})
	})

	if saveToFile {
		f, err := os.Create("data/recipes.json")
		if err != nil {
			log.Fatal("Failed to create file:", err)
		}
		defer f.Close()
	
		var out []RecipeJSON
		for pair, resultList := range combinations {
			for _, result := range resultList {
				out = append(out, RecipeJSON{
					Result:       result,
					Combinations: pair,
					Asset:				assetMap[result],
				})
			}
		}

		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			log.Fatal("Failed to write JSON:", err)
		}
		fmt.Println("Done. Saved to data/recipes.json")
	} else {
		fmt.Print("Done. ")
	}

	fmt.Println("Lookup available in memory.")
}