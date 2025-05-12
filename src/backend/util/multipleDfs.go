package util

import (
	"strings"
)

// MultipleRecipesResult holds the results of MultipleRecipesDfs
type MultipleRecipesResult struct {
	Recipes   []map[string]Element // Collection of valid recipes
	NodeCount int                  // Number of nodes/elements visited
}

// MultipleRecipesDfs finds multiple valid recipe paths for a target element by:
// 1. Finding the first valid recipe
// 2. Backtracking through the recipe tree to find variations
func MultipleDfs(target string, revCombinations map[string][]Pair, tierMap map[string]int, maxRecipes int) MultipleRecipesResult {
	// First, find the initial recipe using the standard ShortestDfs
	firstRecipe := ShortestDfs(target, revCombinations, tierMap)
	
	// Track all elements visited during exploration
	visited := make(map[string]bool)
	
	// Initialize visited with base elements
	baseElements := []string{"Fire", "Water", "Air", "Earth"}
	for _, elem := range baseElements {
		visited[elem] = true
	}
	
	// If no initial recipe found, return empty result
	if len(firstRecipe) == 0 {
		return MultipleRecipesResult{
			Recipes:   []map[string]Element{},
			NodeCount: len(visited),
		}
	}
	
	// Collection of recipes, starting with the first one
	recipes := []map[string]Element{firstRecipe}
	
	// Track elements in the recipes to measure visited nodes
	for elem := range firstRecipe {
		visited[elem] = true
	}
	
	// Find all elements in the recipe tree that have alternative recipes
	// Start from target and traverse down through each ingredient
	elementsToExplore := findElementsWithAlternatives(target, firstRecipe, revCombinations, tierMap)
	
	// For each element with alternatives, try to generate new recipes
	for _, element := range elementsToExplore {
		// Stop if we've reached max recipes
		if maxRecipes > 0 && len(recipes) >= maxRecipes {
			break
		}
		
		// Get all the recipes we have so far
		currentRecipes := make([]map[string]Element, len(recipes))
		copy(currentRecipes, recipes)
		
		// For each existing recipe, try to create variations by changing this element
		for _, baseRecipe := range currentRecipes {
			// Stop if we've reached max recipes
			if maxRecipes > 0 && len(recipes) >= maxRecipes {
				break
			}
			
			// Get all valid pairs for this element
			pairs := revCombinations[element]
			validPairs := filterValidPairs(pairs, element, tierMap)
			
			// Get the pair used in the current recipe
			currentPair := Pair{
				First:  baseRecipe[element].Source,
				Second: baseRecipe[element].Partner,
			}
			
			// Try each alternative pair
			for _, pair := range validPairs {
				// Skip the pair already used in this recipe
				if (pair.First == currentPair.First && pair.Second == currentPair.Second) || 
				(pair.Second == currentPair.First && pair.First == currentPair.Second) {
					continue
				}
				
				// Stop if we've reached max recipes
				if maxRecipes > 0 && len(recipes) >= maxRecipes {
					break
				}
				
				// Create a variation of the recipe by replacing this element's recipe
				variation := copyRecipe(baseRecipe)
				variation[element] = Element{Source: pair.First, Partner: pair.Second}
				
				// Check if this change creates a valid complete recipe
				valid, elementsVisited := repairRecipeAfterChange(element, variation, revCombinations, tierMap)
				
				// Update visited elements
				for elem := range elementsVisited {
					visited[elem] = true
				}
				
				// If valid and unique, add to our collection
				if valid && isUniqueRecipe(variation, recipes) {
					recipes = append(recipes, variation)
				}
			}
		}
	}
	
	return MultipleRecipesResult{
		Recipes:   recipes,
		NodeCount: len(visited),
	}
}

// findElementsWithAlternatives identifies elements in the recipe tree that have multiple valid recipes
// Returns elements sorted by their position in the recipe tree (leaf to root)
func findElementsWithAlternatives(target string, recipe map[string]Element, revCombinations map[string][]Pair, tierMap map[string]int) []string {
	result := []string{}
	processed := make(map[string]bool)
	
	// Recursive function to explore the recipe tree
	var explore func(element string)
	explore = func(element string) {
		// Skip if already processed or is a base element
		if processed[element] || isBaseElement(element) {
			return
		}
		processed[element] = true
		
		// Check if this element has alternative recipes
		pairs := revCombinations[element]
		validPairs := filterValidPairs(pairs, element, tierMap)
		if len(validPairs) > 1 {
			result = append(result, element)
		}
		
		// Explore ingredients
		elemRecipe, exists := recipe[element]
		if exists && elemRecipe.Source != "" && elemRecipe.Partner != "" {
			explore(elemRecipe.Source)
			explore(elemRecipe.Partner)
		}
	}
	
	// Start exploration from the target
	explore(target)
	
	return result
}

// repairRecipeAfterChange ensures a recipe is still valid after changing one element's recipe
// Returns whether the repair was successful and a map of elements visited during repair
func repairRecipeAfterChange(changedElement string, recipe map[string]Element, revCombinations map[string][]Pair, tierMap map[string]int) (bool, map[string]bool) {
	visited := make(map[string]bool)
	
	// Mark base elements as visited
	for _, elem := range []string{"Fire", "Water", "Air", "Earth"} {
		visited[elem] = true
	}
	
	// Collect all elements that need to be checked/repaired
	// Starting from the changed element's ingredients
	elementsToCheck := []string{}
	changedRecipe := recipe[changedElement]
	
	// Add the changed element's ingredients to check list
	if !isBaseElement(changedRecipe.Source) {
		elementsToCheck = append(elementsToCheck, changedRecipe.Source)
	}
	if !isBaseElement(changedRecipe.Partner) {
		elementsToCheck = append(elementsToCheck, changedRecipe.Partner)
	}
	
	// For each element to check
	for len(elementsToCheck) > 0 {
		// Get the next element
		element := elementsToCheck[0]
		elementsToCheck = elementsToCheck[1:]
		
		// Skip if already processed
		if visited[element] {
			continue
		}
		visited[element] = true
		
		// Check if this element has a valid recipe in the current state
		// If not, find one using ShortestDfs
		if _, exists := recipe[element]; !exists || recipe[element].Source == "" || recipe[element].Partner == "" {
			// Run ShortestDfs just for this element
			miniResult := ShortestDfs(element, revCombinations, tierMap)
			
			// If no recipe found, repair fails
			if len(miniResult) == 0 || miniResult[element].Source == "" || miniResult[element].Partner == "" {
				return false, visited
			}
			
			// Add this recipe to our current recipe map
			recipe[element] = miniResult[element]
			
			// Add all elements from miniResult to our recipe map
			for elem, r := range miniResult {
				if elem != element {
					recipe[elem] = r
					visited[elem] = true
				}
			}
		}
		
		// Add this element's ingredients to the check list if they're not base elements
		elemRecipe := recipe[element]
		if !isBaseElement(elemRecipe.Source) {
			elementsToCheck = append(elementsToCheck, elemRecipe.Source)
		}
		if !isBaseElement(elemRecipe.Partner) {
			elementsToCheck = append(elementsToCheck, elemRecipe.Partner)
		}
	}
	
	return true, visited
}

// copyRecipe makes a deep copy of a recipe map
func copyRecipe(original map[string]Element) map[string]Element {
	copy := make(map[string]Element)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

// isUniqueRecipe checks if a recipe is not already in the collection
func isUniqueRecipe(newRecipe map[string]Element, collection []map[string]Element) bool {
	// Generate a signature for the new recipe
	newSig := generateRecipeSignature(newRecipe)
	
	// Check against existing recipes
	for _, existingRecipe := range collection {
		existingSig := generateRecipeSignature(existingRecipe)
		if newSig == existingSig {
			return false
		}
	}
	
	return true
}

// generateRecipeSignature creates a string signature that uniquely identifies a recipe path
func generateRecipeSignature(recipe map[string]Element) string {
	// We'll use a map to track which elements we've processed
	processed := make(map[string]bool)
	
	// Create a signature starting from non-base elements
	var generateElemSignature func(elem string) string
	generateElemSignature = func(elem string) string {
		// Base elements have fixed signatures
		if isBaseElement(elem) {
			return elem
		}
		
		// Avoid processing the same element twice
		if processed[elem] {
			return elem
		}
		processed[elem] = true
		
		// Get recipe for this element
		elemRecipe, exists := recipe[elem]
		if !exists || elemRecipe.Source == "" || elemRecipe.Partner == "" {
			return elem
		}
		
		// Recursively generate signatures for ingredients
		// Ensure consistent ordering by sorting
		first := generateElemSignature(elemRecipe.Source)
		second := generateElemSignature(elemRecipe.Partner)
		
		// Sort for consistency (so A+B == B+A)
		if first > second {
			first, second = second, first
		}
		
		return elem + "(" + first + "+" + second + ")"
	}
	
	// Start from all non-base elements in the recipe
	var elements []string
	for elem := range recipe {
		if !isBaseElement(elem) && (recipe[elem].Source != "" || recipe[elem].Partner != "") {
			elements = append(elements, elem)
		}
	}
	
	// Generate signature starting from each non-base element
	var signatures []string
	for _, elem := range elements {
		signatures = append(signatures, generateElemSignature(elem))
	}
	
	// Join all signatures
	return strings.Join(signatures, ",")
}