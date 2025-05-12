package util

import (
	"runtime"
	"sync"
)

// ElementProcessingResult isinya hasil dari pemrosesan elemen-elemen
type ElementProcessingResult struct {
	Recipes         []map[string]Element
	VisitedElements map[string]bool
}

// OptimizedParallelDfs adalah versi yang lebih kenceng dari MultipleParallelDfs buat nyari banyak resep valid
// Perbaikan utamanya:
// - Ngurangin perebutan mutex
// - Pemrosesan berkelompok
// - Paralelisasi tingkat elemen yang lebih efisien
// - Cek keunikan yang lebih optimal
func OptimizedParallelDfs(target string, revCombinations map[string][]Pair, tierMap map[string]int, maxRecipes int, numWorkers int) MultipleRecipesResult {
	// Set jumlah worker optimal kalo gak ditentuin
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	// Pertama, cari resep awal pake ShortestDfs
	firstRecipe := ShortestDfs(target, revCombinations, tierMap)

	// Pantau elemen yang udah dikunjungi pake map terpusat
	visited := make(map[string]bool)

	// Tambahin elemen dasar ke visited
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

	// Inisialisasi koleksi resep dengan resep pertama
	recipes := []map[string]Element{firstRecipe}

	// Catat elemen di resep buat ngitung node
	for elem := range firstRecipe {
		visited[elem] = true
	}

	// Cari semua elemen di pohon resep yang punya alternatif
	elementsToExplore := findElementsWithAlternatives(target, firstRecipe, revCombinations, tierMap)

	// Proses elemen secara paralel, bukan berurutan
	elementResults := processElementsInParallel(elementsToExplore, recipes, revCombinations, tierMap, maxRecipes, numWorkers)

	// Tambahin semua resep baru ke koleksi kita
	recipes = append(recipes, elementResults.Recipes...)

	// Gabungin elemen-elemen yang udah dikunjungi
	for elem := range elementResults.VisitedElements {
		visited[elem] = true
	}

	return MultipleRecipesResult{
		Recipes:   recipes,
		NodeCount: len(visited),
	}
}

// Proses elemen secara paralel
func processElementsInParallel(elements []string, baseRecipes []map[string]Element,
	revCombinations map[string][]Pair, tierMap map[string]int,
	maxRecipes int, numWorkers int) ElementProcessingResult {

	// Bikin map hasil bersama dengan mutex buat sinkronisasi
	var resultMu sync.Mutex
	result := ElementProcessingResult{
		Recipes:         make([]map[string]Element, 0),
		VisitedElements: make(map[string]bool),
	}

	// Bikin antrian kerjaan buat elemen
	type elementWork struct {
		element string
		recipes []map[string]Element
	}

	workQueue := make(chan elementWork, len(elements))
	var wg sync.WaitGroup

	// Jalanin worker-worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for work := range workQueue {
				// Proses elemen ini dengan semua resep
				localResult := processElementLocally(work.element, work.recipes,
					revCombinations, tierMap, baseRecipes)

				// Satu kuncian buat update hasil bersama dengan semua hasil lokal
				if len(localResult.Recipes) > 0 {
					resultMu.Lock()
					// Tambahin resep baru
					result.Recipes = append(result.Recipes, localResult.Recipes...)
					// Tambahin elemen yang udah dikunjungi
					for elem := range localResult.VisitedElements {
						result.VisitedElements[elem] = true
					}

					// Cek udah nyampe max recipes belum
					reachedMax := maxRecipes > 0 && len(result.Recipes) >= maxRecipes
					resultMu.Unlock()

					if reachedMax {
						break
					}
				}
			}
		}()
	}

	// Antriin kerjaan - satu job per elemen dengan semua resep saat ini
	for _, element := range elements {
		workQueue <- elementWork{
			element: element,
			recipes: baseRecipes,
		}
	}

	// Tutup channel dan tunggu
	close(workQueue)
	wg.Wait()

	return result
}

// Proses satu elemen secara lokal, ngumpulin hasil sebelum sinkronisasi
func processElementLocally(element string, recipes []map[string]Element,
	revCombinations map[string][]Pair, tierMap map[string]int,
	baseRecipes []map[string]Element) ElementProcessingResult {

	localResult := ElementProcessingResult{
		Recipes:         make([]map[string]Element, 0),
		VisitedElements: make(map[string]bool),
	}

	// Ambil pasangan valid buat elemen ini
	pairs := revCombinations[element]
	validPairs := filterValidPairs(pairs, element, tierMap)

	// Proses tiap resep
	for _, baseRecipe := range recipes {
		// Ambil pasangan yang dipake di resep saat ini
		currentPair := Pair{
			First:  baseRecipe[element].Source,
			Second: baseRecipe[element].Partner,
		}

		// Coba setiap pasangan alternatif
		for _, pair := range validPairs {
			// Skip pasangan yang udah dipake di resep ini
			if (pair.First == currentPair.First && pair.Second == currentPair.Second) ||
				(pair.Second == currentPair.First && pair.First == currentPair.Second) {
				continue
			}

			// Bikin variasi resep dengan sumber elemen yang berbeda
			variation := copyRecipe(baseRecipe)
			variation[element] = Element{Source: pair.First, Partner: pair.Second}

			// Cek apakah perubahan ini bikin resep yang valid
			valid, elementsVisited := repairRecipeAfterChange(element, variation, revCombinations, tierMap)

			if valid {
				// Tambahin elemen yang dikunjungi ke tracking lokal
				for elem := range elementsVisited {
					localResult.VisitedElements[elem] = true
				}

				// Cek keunikan resep secara lokal dulu
				unique := true
				for _, existingRecipe := range localResult.Recipes {
					if areRecipesEqual(variation, existingRecipe) {
						unique = false
						break
					}
				}

				if unique && isUniqueRecipe(variation, baseRecipes) {
					localResult.Recipes = append(localResult.Recipes, variation)
				}
			}
		}
	}

	return localResult
}

// Fungsi helper buat ngecek apakah dua resep sama
func areRecipesEqual(a, b map[string]Element) bool {
	if len(a) != len(b) {
		return false
	}

	for elem, aElement := range a {
		bElement, exists := b[elem]
		if !exists {
			return false
		}

		// Cek apakah source dan partner sama (dalam urutan apapun)
		if (aElement.Source != bElement.Source || aElement.Partner != bElement.Partner) &&
			(aElement.Source != bElement.Partner || aElement.Partner != bElement.Source) {
			return false
		}
	}

	return true
}