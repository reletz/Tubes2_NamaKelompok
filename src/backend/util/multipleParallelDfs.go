package util

import (
	"strings"
	"sync"
)

// MultipleParallelDfs nyari banyak resep valid dengan cara paralel
// Implementasi ini ngikutin cara kerja MultipleDfs yang asli tapi pake multithreading
func MultipleParallelDfs(target string, revCombinations map[string][]Pair, tierMap map[string]int, maxRecipes int, numWorkers int) MultipleRecipesResult {
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

// findElementsWithAlternatives nyari elemen di pohon resep yang punya banyak resep valid
// Ngereturn elemen berurutan dari posisinya di pohon resep (dari daun ke akar)
func findElementsWithAlternatives(target string, recipe map[string]Element, revCombinations map[string][]Pair, tierMap map[string]int) []string {
  result := []string{}
  processed := make(map[string]bool)
  
  // Fungsi rekursif buat jelajahin pohon resep
  var explore func(element string)
  explore = func(element string) {
    // Skip kalo udah diproses atau elemen dasar
    if processed[element] || isBaseElement(element) {
      return
    }
    processed[element] = true
    
    // Cek apakah elemen ini punya resep alternatif
    pairs := revCombinations[element]
    validPairs := filterValidPairs(pairs, element, tierMap)
    if len(validPairs) > 1 {
      result = append(result, element)
    }
    
    // Jelajahi bahan-bahannya
    elemRecipe, exists := recipe[element]
    if exists && elemRecipe.Source != "" && elemRecipe.Partner != "" {
      explore(elemRecipe.Source)
      explore(elemRecipe.Partner)
    }
  }
  
  // Mulai jelajah dari target
  explore(target)
  
  return result
}

// repairRecipeAfterChange mastiin resep masih valid setelah ganti resep satu elemen
// Ngereturn apakah perbaikan berhasil dan map elemen yang dikunjungi selama perbaikan
func repairRecipeAfterChange(changedElement string, recipe map[string]Element, revCombinations map[string][]Pair, tierMap map[string]int) (bool, map[string]bool) {
  visited := make(map[string]bool)
  
  // Tandai elemen dasar sebagai visited
  for _, elem := range BaseElements {
    visited[elem] = true
  }
  
  // Kumpulin semua elemen yang perlu dicek/diperbaiki
  // Mulai dari bahan-bahan elemen yang diubah
  elementsToCheck := []string{}
  changedRecipe := recipe[changedElement]
  
  // Tambah bahan-bahan elemen yang diubah ke list cek
  if !isBaseElement(changedRecipe.Source) {
    elementsToCheck = append(elementsToCheck, changedRecipe.Source)
  }
  if !isBaseElement(changedRecipe.Partner) {
    elementsToCheck = append(elementsToCheck, changedRecipe.Partner)
  }
  
  // Buat tiap elemen yang mau dicek
  for len(elementsToCheck) > 0 {
    // Ambil elemen berikutnya
    element := elementsToCheck[0]
    elementsToCheck = elementsToCheck[1:]
    
    // Skip kalo udah diproses
    if visited[element] {
      continue
    }
    visited[element] = true
    
    // Cek apakah elemen ini punya resep valid di kondisi saat ini
    // Kalo gak, cari pake ShortestDfs
    if _, exists := recipe[element]; !exists || recipe[element].Source == "" || recipe[element].Partner == "" {
      // Jalanin ShortestDfs cuma buat elemen ini
      miniResult := ShortestDfs(element, revCombinations, tierMap)
      
      // Kalo gak nemu resep, perbaikan gagal
      if len(miniResult) == 0 || miniResult[element].Source == "" || miniResult[element].Partner == "" {
        return false, visited
      }
      
      // Tambahin resep ini ke map resep kita
      recipe[element] = miniResult[element]
      
      // Tambahin semua elemen dari miniResult ke map resep kita
      for elem, r := range miniResult {
        if elem != element {
          recipe[elem] = r
          visited[elem] = true
        }
      }
    }
    
    // Tambahin bahan-bahan elemen ini ke list cek kalo bukan elemen dasar
    elemRecipe := recipe[element]
    if !isBaseElement(elemRecipe.Source) {
      elementsToCheck = append(elementsToCheck, elemRecipe.Source)
    }
    if !isBaseElement(elemRecipe.Partner) {
      elementsToCheck = append(elementsToCheck, elemRecipe.Partner)
    }
  }
  
  return true, visited
}

// copyRecipe bikin salinan dalam dari map resep
func copyRecipe(original map[string]Element) map[string]Element {
  copy := make(map[string]Element)
  for k, v := range original {
    copy[k] = v
  }
  return copy
}

// isUniqueRecipe ngecek apakah resep belum ada di koleksi
func isUniqueRecipe(newRecipe map[string]Element, collection []map[string]Element) bool {
  // Bikin tanda tangan buat resep baru
  newSig := generateRecipeSignature(newRecipe)
  
  // Cek lawan resep yang udah ada
  for _, existingRecipe := range collection {
    existingSig := generateRecipeSignature(existingRecipe)
    if newSig == existingSig {
      return false
    }
  }
  
  return true
}

// generateRecipeSignature bikin string tanda tangan yang unik buat ngenali satu jalur resep
func generateRecipeSignature(recipe map[string]Element) string {
  // Kita pake map buat ngetrack elemen yang udah diproses
  processed := make(map[string]bool)
  
  // Bikin tanda tangan mulai dari elemen non-dasar
  var generateElemSignature func(elem string) string
  generateElemSignature = func(elem string) string {
    // Elemen dasar punya tanda tangan tetap
    if isBaseElement(elem) {
      return elem
    }
    
    // Hindari memproses elemen yang sama dua kali
    if processed[elem] {
      return elem
    }
    processed[elem] = true
    
    // Ambil resep buat elemen ini
    elemRecipe, exists := recipe[elem]
    if !exists || elemRecipe.Source == "" || elemRecipe.Partner == "" {
      return elem
    }
    
    // Rekursif bikin tanda tangan buat bahan-bahannya
    // Pastiin urutan konsisten dengan mengurutkan
    first := generateElemSignature(elemRecipe.Source)
    second := generateElemSignature(elemRecipe.Partner)
    
    // Urutkan buat konsistensi (supaya A+B == B+A)
    if first > second {
      first, second = second, first
    }
    
    return elem + "(" + first + "+" + second + ")"
  }
  
  // Mulai dari semua elemen non-dasar di resep
  var elements []string
  for elem := range recipe {
    if !isBaseElement(elem) && (recipe[elem].Source != "" || recipe[elem].Partner != "") {
      elements = append(elements, elem)
    }
  }
  
  // Bikin tanda tangan mulai dari tiap elemen non-dasar
  var signatures []string
  for _, elem := range elements {
    signatures = append(signatures, generateElemSignature(elem))
  }
  
  // Gabungin semua tanda tangan
  return strings.Join(signatures, ",")
}