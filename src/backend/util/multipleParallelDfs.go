package util

import (
	"context"
	"sync"
	"sync/atomic"
)

// RecipeTask represents a task for worker goroutines
type RecipeTask struct {
	Element       string             // Element to find alternatives for
	BaseRecipe    map[string]Element // Recipe to create variation from
}

// MultipleParallelDfs finds multiple valid recipe paths using parallel processing
func MultipleParallelDfs(target string, revCombinations map[string][]Pair, tierMap map[string]int, maxRecipes int, numWorkers int) MultipleRecipesResult {
	// If numWorkers is not specified, use a reasonable default
	if numWorkers <= 0 {
		numWorkers = 4 // Default to 4 workers
	}

	// First, find the initial recipe using the standard ShortestDfs
	firstRecipe := ShortestDfs(target, revCombinations, tierMap)
	
	// Create a context that can be canceled when max recipes is reached
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Shared synchronized data
	var mu sync.Mutex
	visited := make(map[string]bool)
	recipes := []map[string]Element{}
	var recipeCount int32 = 0
	
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
	
	// Add the first recipe and update visited elements
	mu.Lock()
	recipes = append(recipes, firstRecipe)
	atomic.AddInt32(&recipeCount, 1)
	for elem := range firstRecipe {
		visited[elem] = true
	}
	mu.Unlock()
	
	// Find all elements in the recipe tree that have alternative recipes
	elementsToExplore := findElementsWithAlternatives(target, firstRecipe, revCombinations, tierMap)
	
	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	
	// Channel for distributing tasks to workers
	taskChan := make(chan RecipeTask, len(elementsToExplore)*10) // Buffer size is a heuristic
	
	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for {
				select {
				case <-ctx.Done():
					// Context was canceled, stop processing
					return
				case task, ok := <-taskChan:
					if !ok {
						// Channel closed, no more tasks
						return
					}
					
					// Process the task - try to create variations
					processElementVariations(
						ctx, task.Element, task.BaseRecipe, 
						revCombinations, tierMap, maxRecipes,
						&mu, &recipes, &visited, &recipeCount, cancel)
				}
			}
		}()
	}
	
	// Create initial tasks - one for each element with alternatives paired with the first recipe
	for _, element := range elementsToExplore {
		// Only create initial tasks if we haven't reached maxRecipes
		if maxRecipes <= 0 || atomic.LoadInt32(&recipeCount) < int32(maxRecipes) {
			taskChan <- RecipeTask{
				Element:    element,
				BaseRecipe: firstRecipe,
			}
		} else {
			cancel() // Signal all workers to stop
			break
		}
	}
	
	// Close the task channel when no more initial tasks need to be added
	close(taskChan)
	
	// Wait for all workers to finish
	wg.Wait()
	
	// Get the final result under a lock
	mu.Lock()
	finalRecipes := make([]map[string]Element, len(recipes))
	copy(finalRecipes, recipes)
	finalVisitedCount := len(visited)
	mu.Unlock()
	
	return MultipleRecipesResult{
		Recipes:   finalRecipes,
		NodeCount: finalVisitedCount,
	}
}

// processElementVariations explores all variations for an element and adds valid recipes to the result
func processElementVariations(
	ctx context.Context,
	element string, 
	baseRecipe map[string]Element, 
	revCombinations map[string][]Pair, 
	tierMap map[string]int, 
	maxRecipes int,
	mu *sync.Mutex,
	recipes *[]map[string]Element,
	visited *map[string]bool,
	recipeCount *int32,
	cancel context.CancelFunc,
) {
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
		// Check if context is done periodically
		select {
		case <-ctx.Done():
			return
		default:
			// Continue processing
		}
		
		// Skip the pair already used in this recipe
		if (pair.First == currentPair.First && pair.Second == currentPair.Second) || 
		   (pair.Second == currentPair.First && pair.First == currentPair.Second) {
			continue
		}
		
		// Check if we've reached maxRecipes
		if maxRecipes > 0 && atomic.LoadInt32(recipeCount) >= int32(maxRecipes) {
			cancel() // Signal all workers to stop
			return
		}
		
		// Create a variation of the recipe by replacing this element's recipe
		variation := copyRecipe(baseRecipe)
		variation[element] = Element{Source: pair.First, Partner: pair.Second}
		
		// Check if this change creates a valid complete recipe
		valid, elementsVisited := repairRecipeAfterChange(element, variation, revCombinations, tierMap)
		
		if valid {
			// Need to check if this recipe is unique under a lock
			mu.Lock()
			
			// Update visited elements
			for elem := range elementsVisited {
				(*visited)[elem] = true
			}
			
			// If valid and unique, add to our collection
			isUnique := isUniqueRecipe(variation, *recipes)
			if isUnique {
				*recipes = append(*recipes, variation)
				currentCount := atomic.AddInt32(recipeCount, 1)
				
				// Check if we've reached maxRecipes after adding
				if maxRecipes > 0 && currentCount >= int32(maxRecipes) {
					mu.Unlock()
					cancel() // Signal all workers to stop
					return
				}
			}
			
			mu.Unlock()
		}
	}
}
