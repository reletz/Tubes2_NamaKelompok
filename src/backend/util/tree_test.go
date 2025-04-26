package util

import (
	"reflect"
	"testing"
)

func TestTreeOperations(t *testing.T) {
	// 1. Test NewTreeNode
	starting := ElType{"fire", "water"}
	root := NewTreeNode(starting)
	if root == nil {
		t.Fatalf("NewTreeNode returned nil")
	}
	if !reflect.DeepEqual(root.Value, starting) {
		t.Errorf("Expected value %v, got %v", starting, root.Value)
	}

	// 2. Test InsertChild
	InsertChild(&root, starting, "steam") // fire + water = steam
	if root.FirstChild == nil {
		t.Fatalf("FirstChild is nil after InsertChild")
	}
	if !contains(root.FirstChild.Value, "steam") {
		t.Errorf("Expected child to contain 'steam', got %v", root.FirstChild.Value)
	}

	// 3. Test NewChildNode (via InsertChild)
	if root.FirstChild.Hash != root.Hash^HashString("steam") {
		t.Errorf("Child hash incorrect. Expected %v, got %v", root.Hash^HashString("steam"), root.FirstChild.Hash)
	}

	// 4. Test FindNode
	found := FindNode(root, root.FirstChild.Value)
	if found == nil {
		t.Fatalf("FindNode failed to find inserted child")
	}
	if found != root.FirstChild {
		t.Errorf("FindNode returned wrong node")
	}

	// 5. Test GetNbElmt
	nb := GetNbElmt(root)
	if nb != 2 {
		t.Errorf("Expected 2 elements, got %d", nb)
	}

	// 6. Test GetNbLeaf
	leaves := GetNbLeaf(root)
	if leaves != 1 {
		t.Errorf("Expected 1 leaf, got %d", leaves)
	}

	// 7. Test GetHeight
	height := GetHeight(root)
	if height != 2 {
		t.Errorf("Expected height 2, got %d", height)
	}
}

// Helper function untuk cek apakah ElType berisi elemen tertentu
func contains(slice []string, element string) bool {
	for _, s := range slice {
		if s == element {
			return true
		}
	}
	return false
}
