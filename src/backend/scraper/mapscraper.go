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

func GraphScraper() {
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

	// Bikin ReverseGraph
	// Generate sementara dalam map[string]map[string]bool
	tempReverse := make(map[string]map[string]bool)

	for result, pairs := range graph {
		for _, pair := range pairs {
			if len(pair) == 2 {
				for _, comp := range pair {
					if _, exists := tempReverse[comp]; !exists {
						tempReverse[comp] = make(map[string]bool)
					}
					tempReverse[comp][result] = true
				}
			}
		}
	}

	// Convert to final map[string][]string
	reverseGraph := make(map[string][]string)
	for comp, resultSet := range tempReverse {
		for result := range resultSet {
			reverseGraph[comp] = append(reverseGraph[comp], result)
		}
	}


	// Simpan graph langsung
	file, _ := os.Create("data/graph.json")
	defer file.Close()

	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		log.Fatal("Gagal format JSON:", err)
	}
	file.Write(data)

	fmt.Println("Graph langsung disimpan ke graph.json")

	// Simpan ReverseGraph
	out, err := json.MarshalIndent(reverseGraph, "", "  ")
	if err != nil {
		log.Fatal("Failed to convert to JSON:", err)
	}

	if err := os.WriteFile("data/reverseGraph.json", out, 0644); err != nil {
		log.Fatal("Failed to save file:", err)
	}

	fmt.Println("File reverseGraph.json berhasil dibuat")
}