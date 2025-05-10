package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ConvertToJSON(node *Node) ([]byte, error) {
  jsonData, err := json.MarshalIndent(node, "", "  ")
  if err != nil {
    return nil, fmt.Errorf("failed to marshal JSON: %v", err)
  }
  return jsonData, nil
}

// SaveToJSON saves a Node tree to a JSON file
func SaveToJSON(node *Node, filename string) error {
  // Create directory if it doesn't exist
  dir := filepath.Dir(filename)
  if err := os.MkdirAll(dir, 0755); err != nil {
    return fmt.Errorf("failed to create directory: %v", err)
  }
  
  // Marshal the node to JSON
  jsonData, err := ConvertToJSON(node)
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
