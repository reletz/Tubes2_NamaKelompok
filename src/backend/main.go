package main

import (
	"backend/scraper"
	"backend/util"
	"fmt"
	"time"
)

func main(){
	rawRecipe := make(map[util.Pair]string)
	reversedRawRecipe := make(map[string][]util.Pair)
	ingredientsTier := make(map[string]int)
	scraper.Scraper(rawRecipe, ingredientsTier, reversedRawRecipe, true)
	
	target := "Cloud"

	// Standard BFS
	start1 := time.Now()
	prev1 := util.ShortestBfs(target, rawRecipe)
	tree1, visited := util.BuildTree(target, prev1)
	elapsed1 := time.Since(start1)
	fmt.Print("BFS for " + target + ", time taken: ")
	fmt.Println(elapsed1)
	util.SaveToJSON([]*util.Node{tree1}, "data/product_tree.json", visited, time.Since(start1))

	// Standard DFS
	start2 := time.Now()
	prev2 := util.ShortestDfs(target, reversedRawRecipe, ingredientsTier)
	tree2, visited := util.BuildTree(target, prev2)
	elapsed2 := time.Since(start2)
	fmt.Print("DFS for " + target + ", time taken: ")
	fmt.Println(elapsed2)
	util.SaveToJSON([]*util.Node{tree2}, "data/product_tree2.json", visited, time.Since(start2))
	
	// MultipleBFS demonstration
	start3 := time.Now()
	multiBfsResult, _ := util.MultipleBfs(target, rawRecipe, 10, ingredientsTier)
	elapsed3 := time.Since(start3)
	util.SaveToJSON(multiBfsResult.Trees, "data/multi_bfs_results.json", multiBfsResult.VisitedNodes, elapsed3)
	fmt.Println(len(multiBfsResult.Trees))
	
	// MultipleDFS demonstration
	start4 := time.Now()
	multiDfsResult := util.MultipleDfs(target, reversedRawRecipe, ingredientsTier, 10)
	elapsed4 := time.Since(start4)
	tree4, visited := util.BuildMultipleTrees(target, multiDfsResult)
	util.SaveToJSON(tree4, "data/multi_dfs_results.json", visited, elapsed4)
	fmt.Println(len(multiDfsResult.Recipes))

	start5 := time.Now()
	multicDfsResult := util.MultipleParallelDfs(target, reversedRawRecipe, ingredientsTier, 10, 10)
	elapsed5 := time.Since(start5)
	tree5, visited := util.BuildMultipleTrees(target, multicDfsResult)
	util.SaveToJSON(tree5, "data/multi_dfs_results2.json", visited, elapsed5)
	fmt.Println(len(multicDfsResult.Recipes))
}