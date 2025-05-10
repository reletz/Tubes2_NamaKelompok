package main

import (
	"backend/scraper"
	"backend/util"
	"fmt"
	"time"
)

func main(){
	rawRecipe := make(map[util.Pair][]string)
	reversedRawRecipe := make(map[string][]util.Pair)
	ingredientsTier := make(map[string]int)
	scraper.Scraper(rawRecipe, ingredientsTier, reversedRawRecipe, true)
	
	target := "Grilled cheese"

	start1 := time.Now()
	prev1 := util.ShortestBfs(target, rawRecipe)
	tree1 := util.BuildTree(target, prev1)
	elapsed1 := time.Since(start1)

	// Save the tree as JSON
	fmt.Print("BFS for " + target + ", time taken: ")
	fmt.Println(elapsed1)
	util.SaveToJSON(tree1, "data/product_tree.json")

	start2 := time.Now()
	recipes := util.ShortestDfs(target, reversedRawRecipe, ingredientsTier)
	
	tree2 := util.BuildTree(target, recipes)
	elapsed2 := time.Since(start2)
	
	util.SaveToJSON(tree2, "data/product_tree2.json")
	fmt.Print("DFS for " + target + ", time taken: ")
	fmt.Println(elapsed2)
}