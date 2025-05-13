package util

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// BidirQueueItem struktur data buat menyimpan item kerja di queue
type BidirQueueItem struct {
	Recipe    map[string]Element // Resep saat ini
	FocusElem string             // Elemen yang lagi kita fokuskan
}

// BidirProcessingResult nyimpen hasil dari worker yang memproses batch
type BidirProcessingResult struct {
	NewRecipes      []map[string]Element // Resep baru yang ditemukan
	NewQueueItems   []BidirQueueItem     // Item baru untuk dimasukin ke queue
	VisitedElements map[string]bool      // Elemen yang dikunjungi selama pemrosesan
}

// MultipleBidirectional nyari banyak resep dengan metode bidirectional
// yang diparalelkan untuk mempercepat proses pencarian
func MultipleBidirectional(target string, combinations map[Pair]string, 
	revCombinations map[string][]Pair, tierMap map[string]int, 
	maxRecipes int, numWorkers int) MultipleRecipesResult {
	
	// Set jumlah worker optimal kalo gak ditentuin
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}
	
	// Pertama, cari resep awal pake ShortestBidirectional
	firstRecipe := ShortestBidirectional(target, combinations, revCombinations, tierMap)
	
	// Pantau elemen yang udah dikunjungi, pake mutex biar aman
	visited := make(map[string]bool)
	visitedMutex := &sync.Mutex{}
	
	// Masukin elemen dasar ke visited
	for _, elem := range BaseElements {
		visited[elem] = true
	}
	
	// Kalo gak nemu resep awal, return hasil kosong
	if len(firstRecipe) == 0 {
		return MultipleRecipesResult{
			Recipes:   []map[string]Element{},
			NodeCount: len(visited),
		}
	}
	
	// Kumpulan resep, mulai dari resep pertama
	recipes := []map[string]Element{firstRecipe}
	recipesMutex := &sync.Mutex{}
	
	// Atomic counter buat ngitung jumlah resep yang udah ditemuin
	recipeCounter := int32(1) // Mulai dari 1 karena udah ada resep pertama
	
	// Catat elemen di resep untuk ngukur berapa node yang dikunjungi
	for elem := range firstRecipe {
		visited[elem] = true
	}
	
	// Pake concurrent map buat nyimpen resep yang udah ditemuin
	seenRecipes := sync.Map{}
	
	// Tandain resep pertama udah diliat
	seenRecipeKey := RecipeToString(firstRecipe, target)
	seenRecipes.Store(seenRecipeKey, true)
	
	// Bikin queue awal
	initialQueue := []BidirQueueItem{}
	
	// Tambahin variasi elemen target ke queue
	initialQueue = append(initialQueue, BidirQueueItem{Recipe: firstRecipe, FocusElem: target})
	
	// Tambahin juga variasi komponen lainnya
	for elem := range firstRecipe {
		if !isBaseElement(elem) && elem != target {
			initialQueue = append(initialQueue, BidirQueueItem{Recipe: firstRecipe, FocusElem: elem})
		}
	}
	
	// Queue yang dishare antar worker, dilindungi mutex
	queue := initialQueue
	queueMutex := &sync.Mutex{}
	
	// Channel buat ngasih sinyal worker buat berhenti
	done := make(chan struct{})
	
	// Pake sync.Once buat mastiin channel cuma ditutup sekali
	var closeOnce sync.Once
	signalDone := func() {
		closeOnce.Do(func() {
			close(done)
		})
	}
	
	// Fungsi buat ngambil batch kerjaan dari queue
	getBatch := func(batchSize int) []BidirQueueItem {
		queueMutex.Lock()
		defer queueMutex.Unlock()
		
		if len(queue) == 0 {
			return nil
		}
		
		// Ambil maksimal batchSize item, atau semua item tersisa kalo lebih sedikit
		size := batchSize
		if size > len(queue) {
			size = len(queue)
		}
		
		batch := queue[:size]
		queue = queue[size:]
		
		return batch
	}
	
	// Fungsi buat nambahin item ke queue
	addToQueue := func(items []BidirQueueItem) {
		if len(items) == 0 {
			return
		}
		
		queueMutex.Lock()
		defer queueMutex.Unlock()
		
		queue = append(queue, items...)
	}
	
	// Buat WaitGroup buat sinkronisasi worker goroutine
	var wg sync.WaitGroup
	
	// Fungsi worker yang memproses batch kerjaan
	worker := func() {
		defer wg.Done()
		
		localVisited := make(map[string]bool)
		
		// Proses batch sampai disuruh berhenti
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
				// Cek counter resep
				if maxRecipes > 0 && atomic.LoadInt32(&recipeCounter) >= int32(maxRecipes) {
					signalDone()
					return
				}
				
				// Ambil batch kerjaan
				batch := getBatch(10) // Proses 10 item sekaligus
				if batch == nil {
					// Gak ada kerjaan lagi, tapi jangan keluar dulu
					runtime.Gosched() // Kasih kesempatan goroutine lain jalan
					continue
				}
				
				result := processBidirBatch(batch, combinations, revCombinations, tierMap, 
					&seenRecipes, localVisited, target, maxRecipes, &recipeCounter)
				
				// Tangani hasilnya
				if len(result.NewRecipes) > 0 {
					recipesMutex.Lock()
					recipes = append(recipes, result.NewRecipes...)
					recipesMutex.Unlock()
				}
				
				// Tambahin item baru ke queue
				if len(result.NewQueueItems) > 0 {
					addToQueue(result.NewQueueItems)
				}
			}
		}
	}
	
	// Jalanin worker goroutine
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}
	
	// Periksa berkala apakah masih ada kerjaan tersisa
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
				
				// Periksa apakah kita udah punya cukup resep
				if maxRecipes > 0 && atomic.LoadInt32(&recipeCounter) >= int32(maxRecipes) {
					signalDone() // Pake fungsi signalDone yang aman
					return
				}
				
				if queueEmpty {
					// Gak ada kerjaan tersisa di queue
					signalDone() // Pake fungsi signalDone yang aman
					return
				}
			}
		}
	}()
	
	// Tunggu sampe semua worker selesai
	wg.Wait()
	
	return MultipleRecipesResult{
		Recipes:   recipes,
		NodeCount: len(visited),
	}
}

// processBidirBatch ngolah satu batch dari queue
// Ngehasilin resep baru, item queue baru, dan elemen yang dikunjungi
func processBidirBatch(batch []BidirQueueItem, combinations map[Pair]string, 
	revCombinations map[string][]Pair, tierMap map[string]int,
	seenRecipes *sync.Map, localVisited map[string]bool,
	target string, maxRecipes int, recipeCounter *int32) BidirProcessingResult {
	
	result := BidirProcessingResult{
		NewRecipes:      make([]map[string]Element, 0),
		NewQueueItems:   make([]BidirQueueItem, 0),
		VisitedElements: make(map[string]bool),
	}
	
	// Proses tiap item dalam batch
	for _, current := range batch {
		// Periksa apakah sudah mencapai batas resep
		if maxRecipes > 0 && atomic.LoadInt32(recipeCounter) >= int32(maxRecipes) {
			break
		}
		
		// Skip elemen dasar
		if isBaseElement(current.FocusElem) {
			continue
		}
		
		// Resep dan elemen fokus saat ini
		currentRecipe := current.Recipe
		focusElem := current.FocusElem
		
		// Resep asli buat elemen ini
		originalSources := currentRecipe[focusElem]
		
		// Cari semua cara valid dengan bidirectional search
		// Kita gunakan gabungan forward dan backward search
		validPairs := filterValidPairs(revCombinations[focusElem], focusElem, tierMap)
		
		// Coba tiap alternatif cara
		for _, pair := range validPairs {
			// Periksa apakah kita sudah mencapai batas resep
			if maxRecipes > 0 && atomic.LoadInt32(recipeCounter) >= int32(maxRecipes) {
				break
			}
			
			// Skip resep yang sama dengan yang ada sekarang
			if (pair.First == originalSources.Source && pair.Second == originalSources.Partner) ||
			   (pair.First == originalSources.Partner && pair.Second == originalSources.Source) {
				continue
			}
			
			// Bikin variasi dengan alternatif ini
			variation := copyRecipe(currentRecipe)
			variation[focusElem] = Element{Source: pair.First, Partner: pair.Second}
			
			// Pastiin bahan baru punya resep valid
			allValid := true
			
			for _, ingredient := range []string{pair.First, pair.Second} {
				if isBaseElement(ingredient) {
					continue // Elemen dasar selalu valid
				}
				
				// Kalo belum punya resep buat bahan ini, cari pake bidirectional
				if _, exists := variation[ingredient]; !exists {
					// Cari resep dengan cara bikin minimap dari ingredient ke elemen dasar
					// Ini mirip dengan ShortestBidirectional tapi dengan scope lebih kecil
					ingredientRecipe := findIngredientRecipeBidir(ingredient, combinations, revCombinations, tierMap, localVisited)
					if len(ingredientRecipe) == 0 {
						allValid = false
						break
					}
					
					// Tambahin resep bahan ke variasi kita
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
				continue // Skip variasi ini kalo gak bisa dilengkapin
			}
			
			// Cek apakah ini resep unik
			recipeStr := RecipeToString(variation, target)
			if _, seen := seenRecipes.LoadOrStore(recipeStr, true); !seen {
				// Tambahin ke resep baru
				result.NewRecipes = append(result.NewRecipes, variation)
				
				// Increment atomic counter
				newCount := atomic.AddInt32(recipeCounter, 1)
				
				// Kalo udah nyampe batas, berhenti nyari resep lagi
				if maxRecipes > 0 && newCount >= int32(maxRecipes) {
					break
				}
				
				// Nambahin variasi untuk tiap komponen dalam resep
				for elem := range variation {
					if !isBaseElement(elem) {
						result.NewQueueItems = append(result.NewQueueItems, BidirQueueItem{
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

// findIngredientRecipeBidir nyari resep untuk suatu bahan pakai pencarian bidirectional
func findIngredientRecipeBidir(ingredient string, combinations map[Pair]string, 
	revCombinations map[string][]Pair, tierMap map[string]int, 
	visited map[string]bool) map[string]Element {
	
	// Kalo udah elemen dasar, gak perlu resep
	if isBaseElement(ingredient) {
		return map[string]Element{}
	}
	
	// Cari resep yang valid dengan ShortestBidirectional
	miniResult := ShortestBidirectional(ingredient, combinations, revCombinations, tierMap)
	
	// Kalo gak ketemu resep, return kosong
	if len(miniResult) == 0 {
		return map[string]Element{}
	}
	
	// Tandai semua elemen dalam resep sebagai visited
	for elem := range miniResult {
		visited[elem] = true
	}
	
	return miniResult
}
