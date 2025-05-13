package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func ConvertToJSON(nodes []*Node, visited int, timetaken time.Duration) ([]byte, error) {
  result := struct {
    Recipes     []*Node       `json:"recipes"`
    TimeTaken   string        `json:"timetaken"`
    NodeVisited int           `json:"node_visited"`
  }{
    Recipes:     nodes,
    TimeTaken:   timetaken.String(),
    NodeVisited: visited,
  }
  
  jsonData, err := json.MarshalIndent(result, "", "  ")
  if err != nil {
    return nil, fmt.Errorf("failed to marshal JSON: %v", err)
  }
  return jsonData, nil
}

// SaveToJSON saves Node trees to a JSON file
func SaveToJSON(nodes []*Node, filename string, visited int, timetaken time.Duration) error {
  // Create directory if it doesn't exist
  dir := filepath.Dir(filename)
  if err := os.MkdirAll(dir, 0755); err != nil {
    return fmt.Errorf("failed to create directory: %v", err)
  }
  
  // Marshal the nodes to JSON
  jsonData, err := ConvertToJSON(nodes, visited, timetaken)
  if err != nil {
    return fmt.Errorf("failed to marshal JSON: %v", err)
  }
  
  // Write to file
  err = os.WriteFile(filename, jsonData, 0644)
  if err != nil {
    return fmt.Errorf("failed to write JSON file: %v", err)
  }
  
  fmt.Printf("Successfully saved to %s\n", filename)
  return nil
}

// Node represents a node in the recipe tree
type Node struct {
  Name     string  `json:"name"`
  Children []*Node `json:"children,omitempty"`
}

// Element represents a recipe with two ingredients
type Element struct {
  Source  string
  Partner string
}

// Pair represents a combination of two elements
type Pair struct {
  First  string
  Second string
}

// To save to JSON
type ResultJSON struct {
	Name     string      `json:"name"`
	Children interface{} `json:"children,omitempty"`
}
