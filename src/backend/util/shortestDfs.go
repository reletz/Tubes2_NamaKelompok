package util

// NodeState tracks the exploration state for each element
type NodeState struct {
	CurrentPairIndex int
	Pairs            []Pair
	ValidPairs       []Pair
	Visited          bool
}

func ShortestDfs(target string, revCombinations map[string][]Pair, tierMap map[string]int) map[string]Element {
  // Initialize result map: element -> its recipe
  result := make(map[string]Element)
  
  // Check if the target exists in the combinations
  if _, exists := revCombinations[target]; !exists {
    return result
  }
  
  // Map to track exploration state for each element
  nodeStates := make(map[string]*NodeState)
  
  // Add base elements to node states with visited=true
  baseElements := []string{"Fire", "Water", "Air", "Earth"}
  for _, elem := range baseElements {
    nodeStates[elem] = &NodeState{Visited: true}
    result[elem] = Element{} // Mark base elements with empty recipes
  }
  
  // Track the elements that we're currently trying to resolve
  inProgress := make(map[string]bool)
  
  // Use a recursive helper function for DFS
  var explore func(element string) bool
  explore = func(element string) bool {
    // Base elements are already solved
    if isBaseElement(element) {
      return true
    }
    
    // Skip if we've already found a solution for this element
    if state := nodeStates[element]; state != nil && state.Visited {
      return true
    }
    
    // Detect cycles - if we're already trying to resolve this element in the current path
    if inProgress[element] {
      return false
    }
    
    // Mark as in progress
    inProgress[element] = true
    defer func() { inProgress[element] = false }()
    
    // Get or create node state
    state := nodeStates[element]
    if state == nil {
      pairs := revCombinations[element]
      validPairs := filterValidPairs(pairs, element, tierMap)
      state = &NodeState{
        CurrentPairIndex: 0,
        Pairs:            pairs,
        ValidPairs:       validPairs,
        Visited:          false,
      }
      nodeStates[element] = state
    }
    
    // Try each valid recipe
    for i := 0; i < len(state.ValidPairs); i++ {
      pair := state.ValidPairs[i]
      
      // Record this recipe tentatively
      result[element] = Element{
        Source:  pair.First,
        Partner: pair.Second,
      }
      
      // Try to resolve both ingredients
      firstResolved := isBaseElement(pair.First) || explore(pair.First)
      if !firstResolved {
        continue // Try next recipe if first ingredient can't be resolved
      }
      
      secondResolved := isBaseElement(pair.Second) || explore(pair.Second)
      if !secondResolved {
        continue // Try next recipe if second ingredient can't be resolved
      }
      
      // Both ingredients resolved - we found a valid recipe
      state.Visited = true
      state.CurrentPairIndex = i + 1 // Remember which recipe we used
      return true
    }
    
    // If we get here, no valid recipe was found
    delete(result, element) // Remove any tentative recipe
    return false
  }
  
  // Start exploration from the target
  explore(target)
  
  // Clean up the result map - remove any entries with empty recipes that aren't base elements
  for key, elem := range result {
    if !isBaseElement(key) && (elem.Source == "" || elem.Partner == "") {
      delete(result, key)
    }
  }
  
  return result
}

// Improved filterValidPairs to prioritize lower tier ingredients
func filterValidPairs(pairs []Pair, element string, tierMap map[string]int) []Pair {
    currentTier := tierMap[element]
    var validPairs []Pair
    
    for _, pair := range pairs {
        // Only consider pairs where both ingredients are from lower tier
        // This prevents using equal or higher tier elements to create lower tier ones
        if tierMap[pair.First] < currentTier && tierMap[pair.Second] < currentTier {
            validPairs = append(validPairs, pair)
        }
    }
    
    return validPairs
}