package main

import (
	"backend/scraper"
	"backend/util"
	"time"
	"fmt"
)

func main(){
	rawRecipe := make(map[util.Pair][]string)
	reversedRawRecipe := make(map[string][]util.Pair)
	scraper.Scraper(rawRecipe, reversedRawRecipe, true)
	
	target := "Human"

	start1 := time.Now()
	prev1 := util.ShortestBfs(target, rawRecipe)
	fmt.Println(prev1["Acid rain"])
	tree1 := util.BuildTree(target, prev1)
	elapsed1 := time.Since(start1)

	// Save the tree as JSON
	fmt.Print(target + ", time taken: ")
	fmt.Println(elapsed1)
	util.SaveToJSON(tree1, "data/product_tree.json")


	// start := time.Now()
	// recipes := util.MultipleBfs(target, reversedRawRecipe, 3)
	// path2 := util.ConvertPathsToTrees(recipes, rawRecipe)
	// elapsed := time.Since(start)
	
	// fmt.Printf("Found %d recipes for %s in %v\n", len(recipes), target, elapsed)
	// util.SaveMultipleToJSON(path2, "data/awokawok.json")
	
}