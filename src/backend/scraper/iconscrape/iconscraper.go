package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// struktur data untuk menyimpan nama elemen beserta recipe penyusunnya
type Element struct {
	Name    string     `json:"element"`
	Recipes [][]string `json:"recipes"`
}

// fungsi untuk donlot file, di sini dipake untuk donlot SVG
func donlotFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	os.Mkdir("icons", 0755)

	// get request ke halaman wiki Little Alchemy 2
	res, err := http.Get("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// parsing HTML-nya
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var elements []Element
	var wg sync.WaitGroup
	var mu sync.Mutex

	inElements := make(map[string]bool)

	// ambil semua elemen dari tabel
	doc.Find("table.list-table tr").Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			return
		}

		// ambil nama elemen di kolom pertama
		element := strings.TrimSpace(row.Find("td:nth-child(1) a").Text())
		if element == "" {
			return
		}

		// cek apakah elemen sudah diproses
		mu.Lock()
		if inElements[element] {
			mu.Unlock()
			return
		}
		inElements[element] = true
		mu.Unlock()

		// ambil URL dari SVG di kolom pertama
		svgURL := ""
		imgTag := row.Find("td:nth-child(1) .icon-hover a").First()
		if imgHref, exists := imgTag.Attr("href"); exists {
			svgURL = imgHref
		}

		// ambil semua recipe dari kolom kedua
		var recipes [][]string
		row.Find("td:nth-child(2) li").Each(func(j int, li *goquery.Selection) {
			recipe := strings.TrimSpace(li.Text())
			parts := strings.Split(recipe, " + ")
			if len(parts) == 2 {
				recipes = append(recipes, []string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])})
			}
		})

		wg.Add(1)

		go func(element, svgURL string, recipes [][]string) {
			defer wg.Done()
			// donlot SVG
			svgPath := filepath.Join("icons", element+".svg")
			if err := donlotFile(svgURL, svgPath); err != nil {
				log.Printf("Failed to download SVG for %s: %v", element, err)
				return
			}

			// tambahin elemen ke slice
			mu.Lock()
			elements = append(elements, Element{
				Name:    element,
				Recipes: recipes,
			})
			mu.Unlock()
		}(element, svgURL, recipes)
	})

	// tunggu semua goroutine selesai
	wg.Wait()

	// simpen data ke file JSON
	file, err := os.Create("elements.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(elements); err != nil {
		log.Fatal(err)
	}
}
