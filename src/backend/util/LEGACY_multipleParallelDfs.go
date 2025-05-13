package util

import (
	"sync"
)

// MultipleParallelDfs nyari banyak resep valid dengan cara paralel
// Implementasi ini ngikutin cara kerja MultipleDfs yang asli tapi pake multithreading
func Legacy_MultipleParallelDfs(target string, revCombinations map[string][]Pair, tierMap map[string]int, maxRecipes int, numWorkers int) MultipleRecipesResult {
  // Kalo numWorkers gak diisi, kita pake nilai default aja
  if numWorkers <= 0 {
    numWorkers = 4 // Default pake 4 worker
  }

  // Pertama, cari resep awal pake ShortestDfs biasa
  firstRecipe := ShortestDfs(target, revCombinations, tierMap)
  
  // Pantau semua elemen yang udah dikunjungi, pake mutex biar aman
  var mu sync.Mutex
  visited := make(map[string]bool)
  
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
  
  // Catat elemen di resep untuk ngukur berapa node yang dikunjungi
  for elem := range firstRecipe {
    visited[elem] = true
  }
  
  // Cari semua elemen di pohon resep yang punya resep alternatif
  // Mulai dari target terus telusurin ke bawah lewat bahan-bahannya
  elementsToExplore := findElementsWithAlternatives(target, firstRecipe, revCombinations, tierMap)
  
  // Bikin wait group buat proses paralel
  var wg sync.WaitGroup
  
  // Fungsi buat proses variasi untuk satu elemen dalam resep
  processElementInRecipe := func(element string, baseRecipe map[string]Element) []map[string]Element {
    newRecipes := []map[string]Element{}
    
    // Ambil semua pasangan yang valid buat elemen ini
    pairs := revCombinations[element]
    validPairs := filterValidPairs(pairs, element, tierMap)
    
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
      
      // Bikin variasi resep dengan ganti resep elemen ini
      variation := copyRecipe(baseRecipe)
      variation[element] = Element{Source: pair.First, Partner: pair.Second}
      
      // Cek apakah perubahan ini bikin resep yang valid
      valid, elementsVisited := repairRecipeAfterChange(element, variation, revCombinations, tierMap)
      
      if valid {
        // Perlu cek apakah resep ini unik
        mu.Lock()
        // Update elemen yang udah dikunjungi
        for elem := range elementsVisited {
          visited[elem] = true
        }
        
        // Kalo valid dan unik, tambahin ke hasil lokal kita
        if isUniqueRecipe(variation, recipes) {
          newRecipes = append(newRecipes, variation)
        }
        mu.Unlock()
      }
    }
    
    return newRecipes
  }
  
  // Loop semua elemen yang mau dieksplorasi
  for _, element := range elementsToExplore {
    // Ambil semua resep yang udah kita punya
    mu.Lock()
    currentRecipes := make([]map[string]Element, len(recipes))
    copy(currentRecipes, recipes)
    mu.Unlock()
    
    // Bikin channel buat komunikasi worker
    type workItem struct {
      recipeIndex int
      recipe      map[string]Element
    }
    workChan := make(chan workItem, len(currentRecipes))
    resultChan := make(chan []map[string]Element, len(currentRecipes))
    
    // Jalanin goroutine worker
    for i := 0; i < numWorkers; i++ {
      wg.Add(1)
      go func() {
        defer wg.Done()
        for work := range workChan {
          // Proses resep ini
          results := processElementInRecipe(element, work.recipe)
          resultChan <- results
        }
      }()
    }
    
    // Kirim kerjaan ke worker
    for i, recipe := range currentRecipes {
      workChan <- workItem{
        recipeIndex: i,
        recipe:      recipe,
      }
    }
    close(workChan)
    
    // Tunggu callback (masih dalam loop elemen)
    var allNewRecipes []map[string]Element
    for i := 0; i < len(currentRecipes); i++ {
      newRecipes := <-resultChan
      allNewRecipes = append(allNewRecipes, newRecipes...)
    }
    
    // Tambahin semua resep baru ke koleksi utama
    mu.Lock()
    for _, newRecipe := range allNewRecipes {
      // Double-check keunikan sebelum nambah
      if isUniqueRecipe(newRecipe, recipes) {
        recipes = append(recipes, newRecipe)
        
        // Cek udah nyampe maxRecipes belum
        if maxRecipes > 0 && len(recipes) >= maxRecipes {
          break
        }
      }
    }
    mu.Unlock()
    
    // Cek udah nyampe maxRecipes belum
    if maxRecipes > 0 && len(recipes) >= maxRecipes {
      break
    }
  }
  
  // Tunggu sampe semua worker selesai
  wg.Wait()
  
  // Return hasilnya
  return MultipleRecipesResult{
    Recipes:   recipes,
    NodeCount: len(visited),
  }
}