package util

var BaseElements = []string{"Air", "Earth", "Fire", "Water"}

// BuildTree builds a tree from a recipe map iteratively to avoid stack overflow
func BuildTree(element string, recipeMap map[string]Element) (*Node, int) {
  // Create a map to store nodes we've already built
  nodeMap := make(map[string]*Node)
  visited := len(recipeMap)
  
  // Create nodes for base elements first
  for _, base := range BaseElements {
    nodeMap[base] = &Node{Name: base}
  }
  
  // Create a queue of elements to process
  queue := []string{element}
  processed := make(map[string]bool)
  
  // Process elements in the queue
  for len(queue) > 0 {
    // Get the next element
    current := queue[0]
    queue = queue[1:]
    
    // Skip if we've already processed this node
    if processed[current] {
      continue
    }
    processed[current] = true
    
    // Create the current node if it doesn't exist
    if _, exists := nodeMap[current]; !exists {
      nodeMap[current] = &Node{Name: current, Children: []*Node{}}
    }
    
    // Get the recipe for this element
    recipe, exists := recipeMap[current]
    if !exists || (recipe.Source == "" && recipe.Partner == "") {
      // This is either a base element or has no recipe
      continue
    }
    
    // Add children to queue if they haven't been processed
    if !processed[recipe.Source] {
      queue = append(queue, recipe.Source)
    }
    
    if !processed[recipe.Partner] {
      queue = append(queue, recipe.Partner)
    }
    
    // Create child nodes if they don't exist
    if _, exists := nodeMap[recipe.Source]; !exists {
      nodeMap[recipe.Source] = &Node{Name: recipe.Source, Children: []*Node{}}
    }
    
    if _, exists := nodeMap[recipe.Partner]; !exists {
      nodeMap[recipe.Partner] = &Node{Name: recipe.Partner, Children: []*Node{}}
    }
    
    // Add children to the current node
    currentNode := nodeMap[current]
    currentNode.Children = append(currentNode.Children, nodeMap[recipe.Source])
    currentNode.Children = append(currentNode.Children, nodeMap[recipe.Partner])
  }
  
  return nodeMap[element], visited
}

// BuildMultipleTrees builds trees for each recipe in MultipleRecipesResult
func BuildMultipleTrees(element string, result MultipleRecipesResult) ([]*Node, int) {
  var trees []*Node
  
  // Build a tree for each recipe
  for _, recipe := range result.Recipes {
    tree, _ := BuildTree(element, recipe)
    trees = append(trees, tree)
  }
  
  // Return the trees and the NodeCount from the result
  return trees, result.NodeCount
}

// Helper function to check if an element is a base element
func isBaseElement(element string) bool {
	BaseElements := map[string]bool{
		"Fire": true,
		"Water": true,
		"Air": true,
		"Earth": true,
	}
	return BaseElements[element]
}
