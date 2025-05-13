package util

import (
	"runtime"
	"sync"
	"sync/atomic" // Tambahkan import untuk atomic
	"time"
)

// DFSWorkItem represents an element to be processed in our DFS algorithm
type DFSWorkItem struct {
	Element       string               // Elemen yang sedang difokuskan
	BaseRecipes   []map[string]Element // Resep-resep yang sedang dikerjakan
	ExploredPairs map[string]bool      // Pairs yang sudah dieksplorasi untuk element ini
}

// DFSProcessingResult menyimpan hasil pemrosesan dari satu worker DFS
type DFSProcessingResult struct {
	NewRecipes      []map[string]Element // Resep baru yang ditemukan
	NewWorkItems    []DFSWorkItem        // Item kerja baru untuk diproses
	VisitedElements map[string]bool      // Elemen yang dikunjungi
}

// MultipleDfs implementasi DFS yang diparalelkan
// Menggunakan atomic counter untuk melacak jumlah resep yang dihasilkan
func MultipleDfs(target string, revCombinations map[string][]Pair, tierMap map[string]int, maxRecipes int, numWorkers int) MultipleRecipesResult {
	// Set jumlah worker optimal jika tidak ditentukan
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	// Pertama, cari resep awal pakai ShortestDfs
	firstRecipe := ShortestDfs(target, revCombinations, tierMap)

	// Pantau elemen yang sudah dikunjungi
	visited := make(map[string]bool)
	visitedMutex := &sync.Mutex{}

	// Tambahkan elemen dasar ke visited
	for _, elem := range BaseElements {
		visited[elem] = true
	}

	// Jika tidak menemukan resep awal, return hasil kosong
	if len(firstRecipe) == 0 {
		return MultipleRecipesResult{
			Recipes:   []map[string]Element{},
			NodeCount: len(visited),
		}
	}

	// Kumpulan resep, mulai dari resep pertama
	recipes := []map[string]Element{firstRecipe}
	recipesMutex := &sync.Mutex{}

	// Catat elemen di resep untuk mengukur berapa node yang dikunjungi
	for elem := range firstRecipe {
		visited[elem] = true
	}

	// Atomic counter untuk menghitung jumlah resep yang ditemukan
	// Mulai dari 1 karena kita sudah memiliki resep pertama
	recipeCounter := int32(1)

	// Track recipes yang sudah dilihat untuk menghindari duplikat
	seenRecipes := sync.Map{}

	// Catat resep pertama sebagai sudah dilihat
	seenRecipeKey := RecipeToString(firstRecipe, target)
	seenRecipes.Store(seenRecipeKey, true)

	// Cari elemen-elemen yang memiliki alternatif untuk dieksplorasi
	elementsToExplore := findElementsWithAlternatives(target, firstRecipe, revCombinations, tierMap)

	// Buat work stack awal untuk DFS
	workStack := make([]DFSWorkItem, 0, len(elementsToExplore))
	
	// Isi work stack dengan elemen-elemen yang perlu dieksplorasi
	for _, elem := range elementsToExplore {
		workStack = append(workStack, DFSWorkItem{
			Element:       elem,
			BaseRecipes:   []map[string]Element{firstRecipe},
			ExploredPairs: make(map[string]bool),
		})
	}

	// Mutex untuk mengamankan akses ke work stack
	workStackMutex := &sync.Mutex{}

	// Channel untuk memberi sinyal worker untuk berhenti
	done := make(chan struct{})

	// Gunakan sync.Once untuk memastikan kita hanya menutup channel done sekali
	var closeOnce sync.Once
	signalDone := func() {
		closeOnce.Do(func() {
			close(done)
		})
	}

	// Function untuk mendapatkan batch pekerjaan dari stack
	getWorkBatch := func(batchSize int) []DFSWorkItem {
		workStackMutex.Lock()
		defer workStackMutex.Unlock()

		if len(workStack) == 0 {
			return nil
		}

		// Ambil maksimal batchSize item, atau semua item yang tersisa jika lebih sedikit
		size := batchSize
		if size > len(workStack) {
			size = len(workStack)
		}

		// Ambil item dari atas stack (pendekatan DFS)
		startIdx := len(workStack) - size
		batch := workStack[startIdx:]
		workStack = workStack[:startIdx]

		return batch
	}

	// Function untuk menambahkan item ke stack
	addToWorkStack := func(items []DFSWorkItem) {
		if len(items) == 0 {
			return
		}

		workStackMutex.Lock()
		defer workStackMutex.Unlock()

		// Tambahkan ke atas stack untuk diproses selanjutnya (prioritas DFS)
		workStack = append(workStack, items...)
	}

	// Buat WaitGroup untuk menyinkronkan worker goroutine
	var wg sync.WaitGroup

	// Worker function yang memproses batch pekerjaan
	worker := func() {
		defer wg.Done()

		localVisited := make(map[string]bool)

		// Proses batch sampai diberi sinyal untuk berhenti
		for {
			select {
			case <-done:
				// Gabungkan localVisited ke visited global
				if len(localVisited) > 0 {
					visitedMutex.Lock()
					for elem := range localVisited {
						visited[elem] = true
					}
					visitedMutex.Unlock()
				}
				return
			default:
				// Periksa counter resep
				if maxRecipes > 0 && atomic.LoadInt32(&recipeCounter) >= int32(maxRecipes) {
					signalDone()
					return
				}

				// Ambil batch pekerjaan
				batch := getWorkBatch(5) // Process 5 items at a time (smaller batches for DFS)
				if batch == nil {
					// Tidak ada pekerjaan lagi, tapi jangan keluar dulu - bisa jadi akan ditambahkan oleh worker lain
					runtime.Gosched() // Beri kesempatan goroutine lain berjalan
					continue
				}

				result := processWorkBatchAtomic(batch, revCombinations, tierMap, &seenRecipes, localVisited, target, maxRecipes, &recipeCounter)

				// Tangani hasil pemrosesan
				if len(result.NewRecipes) > 0 {
					recipesMutex.Lock()
					recipes = append(recipes, result.NewRecipes...)
					recipesMutex.Unlock()
				}

				// Tambahkan item kerja baru ke stack
				if len(result.NewWorkItems) > 0 {
					addToWorkStack(result.NewWorkItems)
				}
			}
		}
	}

	// Mulai worker goroutine
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	// Periksa secara berkala apakah masih ada pekerjaan tersisa dan semua worker sedang idle
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				workStackMutex.Lock()
				stackEmpty := len(workStack) == 0
				workStackMutex.Unlock()

				// Periksa apakah kita sudah memiliki cukup resep
				if maxRecipes > 0 && atomic.LoadInt32(&recipeCounter) >= int32(maxRecipes) {
					signalDone() // Beri sinyal semua untuk berhenti
					return
				}

				if stackEmpty {
					// Tidak ada pekerjaan tersisa di stack
					signalDone() // Beri sinyal semua untuk berhenti
					return
				}
			}
		}
	}()

	// Tunggu semua worker selesai
	wg.Wait()

	return MultipleRecipesResult{
		Recipes:   recipes,
		NodeCount: len(visited),
	}
}

// processWorkBatchAtomic memproses batch pekerjaan DFS dan menggunakan atomic counter
// untuk melacak jumlah resep yang dihasilkan
func processWorkBatchAtomic(batch []DFSWorkItem, revCombinations map[string][]Pair,
	tierMap map[string]int, seenRecipes *sync.Map, localVisited map[string]bool,
	target string, maxRecipes int, recipeCounter *int32) DFSProcessingResult {
	
	result := DFSProcessingResult{
		NewRecipes:      make([]map[string]Element, 0),
		NewWorkItems:    make([]DFSWorkItem, 0),
		VisitedElements: make(map[string]bool),
	}

	// Helper function untuk normalisasi key pasangan bahan
	pairToString := func(a, b string) string {
		if a > b {
			return b + "+" + a
		}
		return a + "+" + b
	}

	// Proses setiap item pekerjaan dalam batch
	for _, item := range batch {
		// Periksa apakah kita sudah mencapai batas resep
		if maxRecipes > 0 && atomic.LoadInt32(recipeCounter) >= int32(maxRecipes) {
			break
		}

		element := item.Element
		baseRecipes := item.BaseRecipes
		exploredPairs := item.ExploredPairs

		// Ambil semua pasangan valid untuk elemen ini
		pairs := revCombinations[element]
		validPairs := filterValidPairs(pairs, element, tierMap)

		// Untuk setiap resep dasar, coba variasi dengan pasangan valid
		for _, baseRecipe := range baseRecipes {
			// Periksa apakah kita sudah mencapai batas resep
			if maxRecipes > 0 && atomic.LoadInt32(recipeCounter) >= int32(maxRecipes) {
				break
			}

			// Ambil resep asli untuk elemen ini
			originalSources := baseRecipe[element]
			originalPairKey := pairToString(originalSources.Source, originalSources.Partner)

			// Coba setiap pasangan alternatif
			for _, pair := range validPairs {
				// Periksa apakah kita sudah mencapai batas resep
				if maxRecipes > 0 && atomic.LoadInt32(recipeCounter) >= int32(maxRecipes) {
					break
				}

				// Buat key untuk pasangan ini
				pairKey := pairToString(pair.First, pair.Second)

				// Skip jika pasangan ini sama dengan yang asli atau sudah dieksplorasi
				if pairKey == originalPairKey || exploredPairs[pairKey] {
					continue
				}

				// Tandai pasangan ini sudah dieksplorasi
				newExploredPairs := make(map[string]bool)
				for k, v := range exploredPairs {
					newExploredPairs[k] = v
				}
				newExploredPairs[pairKey] = true

				// Buat variasi resep dengan pasangan alternatif ini
				variation := copyRecipe(baseRecipe)
				variation[element] = Element{Source: pair.First, Partner: pair.Second}

				// Pastikan semua bahan baru memiliki resep valid jika belum ada di resep kita
				allValid := true
				for _, ingredient := range []string{pair.First, pair.Second} {
					if isBaseElement(ingredient) {
						continue // Elemen dasar selalu valid
					}

					// Jika kita belum punya resep untuk bahan ini, cari resep
					if _, exists := variation[ingredient]; !exists {
						ingredientRecipe := findIngredientRecipe(ingredient, make(map[Pair]string), revCombinations, tierMap, localVisited)
						if len(ingredientRecipe) == 0 {
							allValid = false
							break
						}

						// Tambahkan resep bahan ke variasi kita
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
					continue // Skip variasi ini jika tidak dapat dilengkapi
				}

				// Pastikan variasi ini valid dengan memperhatikan constraint tier
				valid, elementsVisited := repairRecipeAfterChange(element, variation, revCombinations, tierMap)
				
				for elem := range elementsVisited {
					localVisited[elem] = true
					result.VisitedElements[elem] = true
				}

				if !valid {
					continue // Skip variasi ini jika tidak valid
				}

				// Cek apakah ini resep unik
				recipeStr := RecipeToString(variation, target)
				if _, seen := seenRecipes.LoadOrStore(recipeStr, true); !seen {
					// Tambahkan ke resep baru
					result.NewRecipes = append(result.NewRecipes, variation)
					
					// Increment atomic counter
					newCount := atomic.AddInt32(recipeCounter, 1)
					
					// Jika sudah mencapai batas, berhenti mencari lebih banyak
					if maxRecipes > 0 && newCount >= int32(maxRecipes) {
						break
					}

					// Dalam DFS, kita terus memperdalam eksplorasi untuk resep baru ini
					// Tambahkan semua elemen non-dasar di resep ini untuk eksplorasi lebih lanjut
					for elem := range variation {
						if !isBaseElement(elem) {
							// Cari alternatif untuk elemen ini jika mungkin diubah
							if len(filterValidPairs(revCombinations[elem], elem, tierMap)) > 1 {
								result.NewWorkItems = append(result.NewWorkItems, DFSWorkItem{
									Element:       elem,
									BaseRecipes:   []map[string]Element{variation},
									ExploredPairs: make(map[string]bool),
								})
							}
						}
					}
				}
			}
		}
	}

	return result
}