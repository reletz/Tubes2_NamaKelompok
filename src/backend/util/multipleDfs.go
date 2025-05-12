package util

// MultipleRecipesResult buat nyimpen hasil dari MultipleRecipesDfs
type MultipleRecipesResult struct {
  Recipes   []map[string]Element // Kumpulan resep yang valid
  NodeCount int                  // Jumlah node/elemen yang dikunjungi
}

// MultipleDfs nyari beberapa resep valid buat elemen target dengan cara:
// 1. Nyari resep valid yang pertama
// 2. Backtracking lewat pohon resep buat nemuin variasi lain
func MultipleDfs(target string, revCombinations map[string][]Pair, tierMap map[string]int, maxRecipes int) MultipleRecipesResult {
  // Pertama, cari resep awal pake ShortestDfs biasa
  firstRecipe := ShortestDfs(target, revCombinations, tierMap)
  
  // Pantau semua elemen yang udah dikunjungi
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
  
  // Buat tiap elemen yang punya alternatif, coba bikin resep baru
  for _, element := range elementsToExplore {
    // Berhenti kalo udah nyampe max resep
    if maxRecipes > 0 && len(recipes) >= maxRecipes {
      break
    }
    
    // Ambil semua resep yang udah kita punya
    currentRecipes := make([]map[string]Element, len(recipes))
    copy(currentRecipes, recipes)
    
    // Buat tiap resep yang ada, coba bikin variasi dengan ganti elemen ini
    for _, baseRecipe := range currentRecipes {
      // Berhenti kalo udah nyampe max resep
      if maxRecipes > 0 && len(recipes) >= maxRecipes {
        break
      }
      
      // Ambil semua pasangan yang valid buat elemen ini
      pairs := revCombinations[element]
      validPairs := filterValidPairs(pairs, element, tierMap)
      
      // Ambil pasangan yang dipake di resep saat ini
      currentPair := Pair{
        First:  baseRecipe[element].Source,
        Second: baseRecipe[element].Partner,
      }
      
      // Coba tiap pasangan alternatif
      for _, pair := range validPairs {
        // Skip pasangan yang udah dipake di resep ini
        if (pair.First == currentPair.First && pair.Second == currentPair.Second) || 
        (pair.Second == currentPair.First && pair.First == currentPair.Second) {
          continue
        }
        
        // Berhenti kalo udah nyampe max resep
        if maxRecipes > 0 && len(recipes) >= maxRecipes {
          break
        }
        
        // Bikin variasi resep dengan ganti resep elemen ini
        variation := copyRecipe(baseRecipe)
        variation[element] = Element{Source: pair.First, Partner: pair.Second}
        
        // Cek apakah perubahan ini bikin resep yang valid
        valid, elementsVisited := repairRecipeAfterChange(element, variation, revCombinations, tierMap)
        
        // Update elemen yang udah dikunjungi
        for elem := range elementsVisited {
          visited[elem] = true
        }
        
        // Kalo valid dan unik, tambahin ke koleksi kita
        if valid && isUniqueRecipe(variation, recipes) {
          recipes = append(recipes, variation)
        }
      }
    }
  }
  
  return MultipleRecipesResult{
    Recipes:   recipes,
    NodeCount: len(visited),
  }
}