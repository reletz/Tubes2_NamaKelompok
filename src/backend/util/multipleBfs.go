package util

import (
	"runtime"
	"sync"
)

// MultiTreeResult stores the result of multiple recipe trees
type MultiTreeResult struct {
	Trees        []*Node  `json:"recipes"`
	VisitedNodes int      `json:"node_visited"`
}

// Menyimpan task kombinasi 
type levelTask struct {
	Pair    Pair
	Product string
}

// Menyimpan result dari task yang berhasil
type levelResult struct {
	Product string
	Node    *Node
}

// MultiBFS finds multiple recipes for a target element using breadth-first search
func MultipleBfs(target string, recipeMap map[Pair]string, maxRecipes int, tierMap map[string]int) (MultiTreeResult, error) {
	startingElements := BaseElements
	existing := sync.Map{}
	visitedCombo := sync.Map{}
	visitedElem := sync.Map{}

	// Inisialisasi dari starting element
	for _, e := range startingElements {
		existing.Store(e, &Node{Name: e, Children: []*Node{}})
		visitedElem.Store(e, true)
	}

	// Cek target starting element apa bukan
	for _, e := range startingElements {
		if e == target {
			n := &Node{Name: e, Children: []*Node{}}
			return MultiTreeResult{
				Trees:        []*Node{n},
				VisitedNodes: 1,
			}, nil
		}
	}

	foundRecipes := []*Node{}
	tier := 0
	numWorkers := runtime.NumCPU() * 2

	// Proses Multiple BFS
	for len(foundRecipes) < maxRecipes {
		tasks := make(chan levelTask, 1000)
		results := make(chan levelResult, 1000)
		var wg sync.WaitGroup

		// Worker pool multithreading
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for task := range tasks {
					pair := task.Pair
					product := task.Product

					// Pengecekan untuk resep harus di bawah elemen yang akan dibentuk
					productTier := tierMap[product]
					valid := true
					// Check both ingredients from the Pair
					first, second := pair.First, pair.Second
					if tierMap[first] >= productTier || tierMap[second] >= productTier {
						valid = false
					}
					if !valid {
						continue
					}

					// Pengecekan elemen dari recipe
					n1Raw, ok1 := existing.Load(first)
					n2Raw, ok2 := existing.Load(second)
					if !ok1 || !ok2 {
						continue
					}

					// Menghindari duplikasi kalo ada kombinasi yang sama
					comboKey := first + "+" + second + ">" + product
					if _, dup := visitedCombo.LoadOrStore(comboKey, true); dup {
						continue
					}

					newNode := &Node{Name: product, Children: []*Node{n1Raw.(*Node), n2Raw.(*Node)}}
					results <- levelResult{Product: product, Node: newNode}
				}
			}(i)
		}

		// Mengirim semua kombinasi ke worker
		// Cari semua kombinasi pasangan elemen yang ada
		var existingElements []string
		visitedElem.Range(func(key, _ interface{}) bool {
			existingElements = append(existingElements, key.(string))
			return true
		})

		// Generate all possible combinations of elements we've seen
		for i := 0; i < len(existingElements); i++ {
			for j := i; j < len(existingElements); j++ {
				first := existingElements[i]
				second := existingElements[j]
				
				// Try pair in both orders
				pairs := []Pair{
					{First: first, Second: second},
					{First: second, Second: first},
				}
				
				for _, pair := range pairs {
					if product, exists := recipeMap[pair]; exists {
						if _, seen := visitedElem.Load(product); !seen {
							tasks <- levelTask{Pair: pair, Product: product}
						}
					}
				}
			}
		}
		
		close(tasks)

		go func() {
			wg.Wait()
			close(results)
		}()

		nextCount := 0
		for res := range results {
			existing.Store(res.Product, res.Node)
			visitedElem.Store(res.Product, true)
			nextCount++
			
			if res.Product == target {
				cloned := deepCopyTree(res.Node)
				foundRecipes = append(foundRecipes, cloned)
				if len(foundRecipes) >= maxRecipes {
					break
				}
			}
		}

		if nextCount == 0 {
			break
		}
		tier++
	}

	nodeCount := 0
	visitedElem.Range(func(_, _ any) bool {
		nodeCount++
		return true
	})

	return MultiTreeResult{
		Trees:        foundRecipes,
		VisitedNodes: nodeCount,
	}, nil
}

func deepCopyTree(node *Node) *Node {
	if node == nil {
		return nil
	}
	copy := &Node{
		Name:     node.Name,
		Children: []*Node{},
	}
	for _, child := range node.Children {
		copy.Children = append(copy.Children, deepCopyTree(child))
	}
	return copy
}