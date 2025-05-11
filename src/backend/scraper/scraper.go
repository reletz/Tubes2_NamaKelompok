package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"backend/util"

	"github.com/PuerkitoBio/goquery"
)
type RecipeJSON struct {
	Result       string      `json:"Result"`
	Asset        string      `json:"Asset"`
	Tier         int         `json:"Tier"`
	Combinations []util.Pair `json:"Combinations"`
}

// OrderedCombinations holds both the map and the order
type OrderedCombinations struct {
    Map   map[util.Pair][]string
    Order []util.Pair // Maintains insertion order
}

// OrderedRevCombinations holds both the map and the order
type OrderedRevCombinations struct {
    Map   map[string][]util.Pair
    Order []string // Maintains insertion order
}

// To preserve insertion order in maps
func Scraper(
    combinations map[util.Pair][]string,
    tierMap map[string]int,
    revCombinations map[string][]util.Pair,
    saveToFile bool,
) (OrderedCombinations, OrderedRevCombinations) {
    // Create ordered structures
    orderedCombinations := OrderedCombinations{
        Map:   combinations,
        Order: []util.Pair{},
    }
    
    orderedRevCombinations := OrderedRevCombinations{
        Map:   revCombinations,
        Order: []string{},
    }
    
    // Track seen pairs and results to maintain order
    seenPairs := make(map[util.Pair]bool)
    seenResults := make(map[string]bool)

    // Rest of your scraper code remains the same...
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
    forbiddenElements := make(map[string]bool)
    mythsAndMonstersScraper(forbiddenElements)

    var currentTier int
    // Parse the table and gather combinations

    doc.Find("h3, table.list-table.col-list.icon-hover tbody tr").Each(func(_ int, row *goquery.Selection) {
        if row.Is("h3") {
            currentTier = getTierNumber(strings.TrimSpace(row.Find("span.mw-headline").Text()))
        } else if row.Is("table.list-table.col-list.icon-hover tbody tr") {
            cols := row.Find("td")
            if cols.Length() < 2 {
                return
            }
            result := strings.TrimSpace(cols.Eq(0).Find("a").Text())
            if result == "" || result == "Time" {
                return
            }

            if forbiddenElements[result] {
                return
            }

            hrefVal, exists := cols.Eq(0).Find("span a").Attr("href")
            asset := ""
            if exists {
                asset = strings.TrimSpace(hrefVal)
            }

            if asset == "" {
                // default aja
                asset = "https://static.wikia.nocookie.net/little-alchemy/images/6/63/Time_2.svg/revision/latest?cb=20210827124225"
            }

            assetMap[result] = asset
            tierMap[result] = currentTier

            // Add result to ordered list if not seen before
            if !seenResults[result] {
                orderedRevCombinations.Order = append(orderedRevCombinations.Order, result)
                seenResults[result] = true
            }

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

                if forbiddenElements[parts[0]] || forbiddenElements[parts[1]] {
                    return
                }

                pair := util.Pair{
                    First: parts[0],
                    Second: parts[1],
                }

                // Add pair to ordered list if not seen before
                if !seenPairs[pair] {
                    orderedCombinations.Order = append(orderedCombinations.Order, pair)
                    seenPairs[pair] = true
                }

                combinations[pair] = append(combinations[pair], result)
                revCombinations[result] = append(revCombinations[result], pair)

                // Add both directions
                if parts[0] != parts[1] {
                    reversedPair := util.Pair{
                        First: parts[1],
                        Second: parts[0],
                    }
                    
                    // Add reversed pair to ordered list if not seen before
                    if !seenPairs[reversedPair] {
                        orderedCombinations.Order = append(orderedCombinations.Order, reversedPair)
                        seenPairs[reversedPair] = true
                    }
                    
                    combinations[reversedPair] = append(combinations[reversedPair], result)
                    revCombinations[result] = append(revCombinations[result], reversedPair)
                }
            })
        }
    })

    // Sort the result list by tier for more deterministic output
    sort.Slice(orderedRevCombinations.Order, func(i, j int) bool {
        a := orderedRevCombinations.Order[i]
        b := orderedRevCombinations.Order[j]
        
        // First by tier
        if tierMap[a] != tierMap[b] {
            return tierMap[a] < tierMap[b]
        }
        
        // Then alphabetically within tier
        return a < b
    })

    if saveToFile {
        f, err := os.Create("data/recipes.json")
        if err != nil {
            log.Fatal("Failed to create file:", err)
        }
        defer f.Close()
    
        var out []RecipeJSON
        // Use the ordered list of results to maintain order
        for _, result := range orderedRevCombinations.Order {
            out = append(out, RecipeJSON{
                Result:       result,
                Combinations: revCombinations[result],
                Asset:        assetMap[result],
                Tier:         tierMap[result],
            })
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
    
    return orderedCombinations, orderedRevCombinations
}

func mythsAndMonstersScraper(forbid	map[string]bool){
	url := "https://little-alchemy.fandom.com/wiki/Category:Myths_and_Monsters"
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

	doc.Find("ul li a.category-page__member-link").Each(func(_ int, row *goquery.Selection){
		element := strings.TrimSpace(row.Text())
		forbid[element] = true
	})
}

// getTierNumber converts a tier description string to its corresponding integer value.
// "Starting elements" returns 0, "Tier X elements" returns X.
func getTierNumber(currentTier string) int {
	// Check for "Starting elements"
	if currentTier == "Starting elements" {
		return 0
	}
	
	// For "Tier X elements", extract the number
	var tierNum int
	_, err := fmt.Sscanf(currentTier, "Tier %d elements", &tierNum)
	if err == nil {
		return tierNum
	}
	
	// Return -1 for invalid format
	return -1
}