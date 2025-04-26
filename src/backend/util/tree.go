package util

type ElType []string

type Address *TreeNode

type TreeNode struct {
	Value      ElType
	Hash       uint64
	FirstChild Address
	NextSibling Address
}

// Constructor NewTreeNode dari awal (tanpa parent, full build)
func NewTreeNode(val ElType) Address {
	return &TreeNode{
		Value:      val,
		Hash:       BuildHashFromSlice(val),
		FirstChild: nil,
		NextSibling: nil,
	}
}

// Constructor NewChildNode (dari parent + satu elemen baru)
func NewChildNode(parent Address, newElement string) Address {
	newVal := make([]string, len(parent.Value))
	copy(newVal, parent.Value)
	newVal = append(newVal, newElement)

	return &TreeNode{
		Value:      newVal,
		Hash:       parent.Hash ^ HashString(newElement),
		FirstChild: nil,
		NextSibling: nil,
	}
}

// CreateTree initializes an empty tree
func CreateTree() Address {
	return nil
}

// IsEmpty checks if the tree is empty
func IsEmpty(t Address) bool {
	return t == nil
}

// FindNode searches for a node by hash
func FindNode(t Address, val ElType) Address {
	if t == nil {
		return nil
	}

	targetHash := BuildHashFromSlice(val)

	if t.Hash == targetHash {
		return t
	}

	child := FindNode(t.FirstChild, val)
	if child != nil {
		return child
	}
	return FindNode(t.NextSibling, val)
}

// InsertChild inserts a new child under the parent with new element
func InsertChild(t *Address, parentVal ElType, newElement string) {
	parent := FindNode(*t, parentVal)
	if parent == nil {
		return
	}
	newNode := NewChildNode(parent, newElement)

	if parent.FirstChild == nil {
		parent.FirstChild = newNode
	} else {
		current := parent.FirstChild
		for current.NextSibling != nil {
			current = current.NextSibling
		}
		current.NextSibling = newNode
	}
}

// GetNbElmt returns the number of nodes in the tree
func GetNbElmt(t Address) int {
	if t == nil {
		return 0
	}
	return 1 + GetNbElmt(t.FirstChild) + GetNbElmt(t.NextSibling)
}

// GetNbLeaf returns the number of leaves in the tree
func GetNbLeaf(t Address) int {
	if t == nil {
		return 0
	}
	if t.FirstChild == nil {
		return 1 + GetNbLeaf(t.NextSibling)
	}
	return GetNbLeaf(t.FirstChild) + GetNbLeaf(t.NextSibling)
}

// GetHeight returns the height of the tree
func GetHeight(t Address) int {
	if t == nil {
		return 0
	}
	heightChild := GetHeight(t.FirstChild)
	heightSibling := GetHeight(t.NextSibling)
	if heightChild+1 > heightSibling {
		return heightChild + 1
	}
	return heightSibling
}

// FreeTree releases the tree (in Go handled by GC, but can nullify references)
func FreeTree(t *Address) {
	*t = nil
}