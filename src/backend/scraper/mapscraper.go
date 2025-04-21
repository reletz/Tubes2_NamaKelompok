package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func MapScraper() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	graph := make(map[string][][]string)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to fetch:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Failed to parse HTML:", err)
	}

	doc.Find("table.list-table.col-list.icon-hover tbody tr").Each(func(_ int, row *goquery.Selection) {
		tds := row.Find("td")
		if tds.Length() < 2 {
			return
		}

		result := strings.TrimSpace(tds.Eq(0).Find("a").Text())

		if result == "Fire" || result == "Water" || result == "Air" {
			result = "";
		}

		tds.Eq(1).Find("li").Each(func(_ int, li *goquery.Selection) {
			parts := []string{}
			li.Find("a").Each(func(_ int, a *goquery.Selection) {
				text := strings.TrimSpace(a.Text())
				if text != "" {
					parts = append(parts, text)
				}
			})

			if len(parts) == 2 && result != "" {
				graph[result] = append(graph[result], parts)
			}
		})
	})

	// Simpan graph langsung
	file, _ := os.Create("data/graph.json")
	defer file.Close()

	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		log.Fatal("Gagal format JSON:", err)
	}
	file.Write(data)

	fmt.Println("Graph langsung disimpan ke graph.json")
}
