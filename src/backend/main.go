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
	
	target := "Picnic"

	// Standard BFS
	start1 := time.Now()
	prev1 := util.ShortestBfsFiltered(target, rawRecipe, ingredientsTier)
	tree1, visited := util.BuildTree(target, prev1)
	elapsed1 := time.Since(start1)
	fmt.Print("BFS for " + target + ", time taken: ")
	fmt.Println(elapsed1)
	util.SaveToJSON([]*util.Node{tree1}, "data/product_tree.json", visited, time.Since(start1))

	// Standard Bidirectional
	start2 := time.Now()
	prev2 := util.ShortestBidirectional(target, rawRecipe, reversedRawRecipe, ingredientsTier)
	tree2, visited := util.BuildTree(target, prev2)
	elapsed2 := time.Since(start2)
	fmt.Print("Bidirectional for " + target + ", time taken: ")
	fmt.Println(elapsed2)
	util.SaveToJSON([]*util.Node{tree2}, "data/product_tree2.json", visited, time.Since(start2))
	
	// MultipleBFS demonstration
	start3 := time.Now()
	multiBfsResult := util.MultipleBidirectional(target, rawRecipe, reversedRawRecipe, ingredientsTier, 10000, 4)
	tree3, visited := util.BuildMultipleTrees(target, multiBfsResult)
	elapsed3 := time.Since(start3)
	util.SaveToJSON(tree3, "data/multi_bidir_results.json", visited, elapsed3)
	fmt.Println(len(multiBfsResult.Recipes))
	fmt.Println(elapsed3.Seconds())
	// // MultipleBFS demonstration
	// start3 := time.Now()
	// multiBfsResult := util.MultipleBfs(target, rawRecipe, reversedRawRecipe, ingredientsTier, 10)
	// tree3, visited := util.BuildMultipleTrees(target, multiBfsResult)
	// elapsed3 := time.Since(start3)
	// util.SaveToJSON(tree3, "data/multi_bfs_results.json", visited, elapsed3)
	// fmt.Println(len(multiBfsResult.Recipes))

	// // MultipleBFS demonstration
	// start8 := time.Now()
	// multiBfsResult8 := util.OptimizedParallelBFS(target, rawRecipe, reversedRawRecipe, ingredientsTier, 10, 4)
	// tree8, visited := util.BuildMultipleTrees(target, multiBfsResult8)
	// elapsed8 := time.Since(start8)
	// util.SaveToJSON(tree8, "data/multi_bfs_results_optimize.json", visited, elapsed8)
	// fmt.Println(len(multiBfsResult.Recipes))
	
	// // MultipleDFS demonstration
	// start4 := time.Now()
	// multiDfsResult := util.MultipleDfs(target, reversedRawRecipe, ingredientsTier, 10)
	// elapsed4 := time.Since(start4)
	// tree4, visited := util.BuildMultipleTrees(target, multiDfsResult)
	// util.SaveToJSON(tree4, "data/multi_dfs_results.json", visited, elapsed4)
	// fmt.Println(len(multiDfsResult.Recipes))

	// start5 := time.Now()
	// multicDfsResult := util.MultipleParallelDfs(target, reversedRawRecipe, ingredientsTier, 10, 8)
	// elapsed5 := time.Since(start5)
	// tree5, visited := util.BuildMultipleTrees(target, multicDfsResult)
	// util.SaveToJSON(tree5, "data/multi_dfs_results_parallel.json", visited, elapsed5)
	// fmt.Println(len(multicDfsResult.Recipes))

	// start6 := time.Now()
	// multic1DfsResult := util.OptimizedParallelDfs(target, reversedRawRecipe, ingredientsTier, 10, 4)
	// elapsed6 := time.Since(start6)
	// tree6, visited := util.BuildMultipleTrees(target, multic1DfsResult)
	// util.SaveToJSON(tree6, "data/multi_dfs_results_optimize.json", visited, elapsed6)
	// fmt.Println(len(multic1DfsResult.Recipes))
}