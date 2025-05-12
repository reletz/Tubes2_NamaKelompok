package util

import (
	"runtime"
	"sync"
	"time"
)

// BFSQueueItem struktur data yang disimpan dalam queue untuk BFS
type BFSQueueItem struct {
	Recipe    map[string]Element // Resep saat ini
	FocusElem string            // Elemen yang lagi difokuskan
}

// BFSWorkBatch berisi sekumpulan item yang akan diproses oleh worker
type BFSWorkBatch struct {
	Items []BFSQueueItem // Item-item yang akan diproses
}

// BFSProcessingResult menyimpan hasil pemrosesan dari satu worker
type BFSProcessingResult struct {
	NewRecipes      []map[string]Element // Resep-resep baru yang ditemukan
	NewQueueItems   []BFSQueueItem       // Item-item baru untuk dimasukkan ke queue
	VisitedElements map[string]bool      // Elemen-elemen yang dikunjungi selama pemrosesan
}

// OptimizedParallelBFS implementasi BFS yang dioptimasi dengan paralelisasi
// untuk mencari beberapa resep valid untuk elemen target
func OptimizedParallelBFS(target string, combinations map[Pair]string, revCombinations map[string][]Pair, 
                         tierMap map[string]int, maxRecipes int, numWorkers int) MultipleRecipesResult {
	// Set jumlah worker ke jumlah CPU jika tidak ditentukan
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	// Pertama, cari resep awal pake ShortestBfsFiltered
	firstRecipe := ShortestBfsFiltered(target, combinations, tierMap)
	
	// Pantau semua elemen yang udah dikunjungi
	visited := make(map[string]bool)
	visitedMutex := &sync.Mutex{}
	
	// Masukin elemen dasar ke visited
	for _, elem := range BaseElements {
		visited[elem] = true
	}
	
	// Kalo gak nemu resep awal, yaudah return hasil kosong
	if len(firstRecipe) == 0 {
		return MultipleRecipesResult{
			Recipes:   []map[string]Element{},
			NodeCount: len(visited),
		}
	}
	
	// Kumpulan resep, mulai dari resep pertama
	recipes := []map[string]Element{firstRecipe}
	recipesMutex := &sync.Mutex{}
	
	// Catat elemen di resep untuk ngukur berapa node yang dikunjungi
	for elem := range firstRecipe {
		visited[elem] = true
	}

	// Track recipes we've already seen to avoid duplicates
	// Pake concurrent map buat nyimpen resep yang udah ditemuin, biar gak duplikat
	seenRecipes := sync.Map{}
	
	// Mark first recipe as seen
	seenRecipes.Store(RecipeToString(firstRecipe, target), true)
	
	// Initial BFS queue
	initialQueue := []BFSQueueItem{}
	
	// First add target variations
	initialQueue = append(initialQueue, BFSQueueItem{Recipe: firstRecipe, FocusElem: target})
	
	// Then add component variations
	for elem := range firstRecipe {
		if !isBaseElement(elem) && elem != target {
			initialQueue = append(initialQueue, BFSQueueItem{Recipe: firstRecipe, FocusElem: elem})
		}
	}
	
	// Create a shared queue protected by a mutex
	queue := initialQueue
	queueMutex := &sync.Mutex{}
	
	// Create a WaitGroup to synchronize worker goroutines
	var wg sync.WaitGroup
	
	// Channel buat kasih sinyal worker untuk berhenti
	done := make(chan struct{})
	
	// Pake sync.Once buat mastiin kita cuma nutup channel done sekali
	var closeOnce sync.Once
	signalDone := func() {
		closeOnce.Do(func() {
			close(done)
		})
	}
	
	// Function to get a batch of work from the queue
	getBatch := func(batchSize int) []BFSQueueItem {
		queueMutex.Lock()
		defer queueMutex.Unlock()
		
		if len(queue) == 0 {
			return nil
		}
		
		// Get up to batchSize items, or all remaining items if less
		size := batchSize
		if size > len(queue) {
			size = len(queue)
		}
		
		batch := queue[:size]
		queue = queue[size:]
		
		return batch
	}
	
	// Function to add items to the queue
	addToQueue := func(items []BFSQueueItem) {
		if len(items) == 0 {
			return
		}
		
		queueMutex.Lock()
		defer queueMutex.Unlock()
		
		queue = append(queue, items...)
	}
	
	// Worker function yang memproses batch pekerjaan
	worker := func() {
		defer wg.Done()
		
		localVisited := make(map[string]bool)
		
		// Proses batch sampai diberi sinyal untuk berhenti
		for {
			select {
			case <-done:
				// Gabungin localVisited ke visited global
				if len(localVisited) > 0 {
					visitedMutex.Lock()
					for elem := range localVisited {
						visited[elem] = true
					}
					visitedMutex.Unlock()
				}
				return
			default:
				// Ambil batch kerjaan
				batch := getBatch(10) // Proses 10 item sekaligus
				if batch == nil {
					// Gak ada kerjaan lagi, tapi jangan keluar dulu - mungkin ditambah sama worker lain
					runtime.Gosched() // Kasih kesempatan goroutine lain jalan
					continue
				}
				
				result := processBatch(batch, combinations, revCombinations, tierMap, 
					&seenRecipes, localVisited, target)
				
				// Tangani hasilnya
				if len(result.NewRecipes) > 0 {
					recipesMutex.Lock()
					recipes = append(recipes, result.NewRecipes...)
					
					// Cek udah nyampe max recipes belum
					reachedMax := maxRecipes > 0 && len(recipes) >= maxRecipes
					recipesMutex.Unlock()
					
					if reachedMax {
						signalDone() // Pake fungsi signalDone yang aman
						return
					}
				}
				
				// Tambahin item baru ke queue
				if len(result.NewQueueItems) > 0 {
					addToQueue(result.NewQueueItems)
				}
			}
		}
	}
	
	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}
	
	// Periksa berkala apakah masih ada kerjaan tersisa dan semua worker sedang idle
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				queueMutex.Lock()
				queueEmpty := len(queue) == 0
				queueMutex.Unlock()
				
				if queueEmpty {
					// Gak ada kerjaan tersisa di queue - cek apa kita perlu berhenti
					recipesMutex.Lock()
					reachedMax := maxRecipes > 0 && len(recipes) >= maxRecipes
					recipesMutex.Unlock()
					
					if reachedMax || queueEmpty {
						// Kita udah punya cukup resep atau udah proses semua
						signalDone() // Pake fungsi signalDone yang aman
						return
					}
				}
			}
		}
	}()
	
	// Wait for all workers to finish
	wg.Wait()
	
	// Merge any remaining local visited maps into the global one
	// This is handled in the worker exit code
	
	return MultipleRecipesResult{
		Recipes:   recipes,
		NodeCount: len(visited),
	}
}

// processBatch handles processing a batch of queue items
// Returns new recipes, new queue items, and visited elements
func processBatch(batch []BFSQueueItem, combinations map[Pair]string, 
                 revCombinations map[string][]Pair, tierMap map[string]int,
                 seenRecipes *sync.Map, localVisited map[string]bool,
                 target string) BFSProcessingResult { // Add target parameter here
	result := BFSProcessingResult{
		NewRecipes:      make([]map[string]Element, 0),
		NewQueueItems:   make([]BFSQueueItem, 0),
		VisitedElements: make(map[string]bool),
	}
	
	// Process each queue item in the batch
	for _, current := range batch {
		// Skip base elements
		if isBaseElement(current.FocusElem) {
			continue
		}
		
		// Current recipe and focus element
		currentRecipe := current.Recipe
		focusElem := current.FocusElem
		
		// Original recipe for this element
		originalSources := currentRecipe[focusElem]
		
		// Get all valid ways to make this element
		validPairs := filterValidPairs(revCombinations[focusElem], focusElem, tierMap)
		
		// Try each alternative way to make this element
		for _, pair := range validPairs {
			// Skip the current recipe for this element
			if (pair.First == originalSources.Source && pair.Second == originalSources.Partner) ||
				(pair.First == originalSources.Partner && pair.Second == originalSources.Source) {
				continue
			}
			
			// Create a variation with this alternative
			variation := copyRecipe(currentRecipe)
			variation[focusElem] = Element{Source: pair.First, Partner: pair.Second}
			
			// Ensure the new ingredients have valid recipes if they're not already in our recipe
			allValid := true
			
			for _, ingredient := range []string{pair.First, pair.Second} {
				if isBaseElement(ingredient) {
					continue // Base elements are always valid
				}
				
				// If we don't have a recipe for this ingredient yet, find one
				if _, exists := variation[ingredient]; !exists {
					ingredientRecipe := findIngredientRecipe(ingredient, combinations, revCombinations, tierMap, localVisited)
					if len(ingredientRecipe) == 0 {
						allValid = false
						break
					}
					
					// Add the ingredient's recipe to our variation
					for elem, sources := range ingredientRecipe {
						if _, exists := variation[elem]; !exists {
							variation[elem] = sources
							localVisited[elem] = true
							result.VisitedElements[elem] = true
						}
					}
				}
			}
			
			if !allValid {
				continue // Skip this variation if we couldn't complete it
			}
			
			// Check if this is a unique recipe
			recipeStr := RecipeToString(variation, target)
			if _, seen := seenRecipes.LoadOrStore(recipeStr, true); !seen {
				// Add to new recipes
				result.NewRecipes = append(result.NewRecipes, variation)
				
				// Add variations for each component in our recipe (BFS approach)
				for elem := range variation {
					if !isBaseElement(elem) {
						result.NewQueueItems = append(result.NewQueueItems, BFSQueueItem{
							Recipe:    variation,
							FocusElem: elem,
						})
					}
				}
			}
		}
	}
	
	return result
}
