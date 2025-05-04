package main

import (
	"backend/scraper"
	"backend/util"
	"time"
	"fmt"
)

func main(){
	rawRecipe := make(map[util.Pair][]string)
	scraper.Scraper(rawRecipe, true)
	
	target := "Human"

	start1 := time.Now()
	prev1 := util.ShortestBfs(target, rawRecipe)
	tree1 := util.BuildTree(target, prev1)
	elapsed1 := time.Since(start1)

	// Save the tree as JSON
	fmt.Print(target + ", time taken: ")
	fmt.Println(elapsed1)
	util.SaveToJSON(tree1, "data/product_tree.json")


	start := time.Now()
	recipes := util.MultipleRecipeBFS(target, rawRecipe, 8)
	elapsed := time.Since(start)
	
	fmt.Printf("Found %d recipes for %s in %v\n", len(recipes), target, elapsed)
	
	// Save each recipe tree to a separate JSON file
	for i, recipe := range recipes {
			filename := fmt.Sprintf("data/product_tree_%s_recipe%d.json", target, i+1)
			util.SaveToJSON(recipe, filename)
			fmt.Printf("Recipe %d saved to %s\n", i+1, filename)
	}
}