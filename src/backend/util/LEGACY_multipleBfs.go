package util

// MultipleBfs nyari beberapa resep valid buat elemen target dengan cara:
// 1. Nyari resep valid yang pertama dengan BFS
// 2. Mengeksplorasi alternatif resep secara melebar (BFS) untuk menemukan variasi lain
// 3. Mencari variasi bukan hanya di level teratas, tapi juga komponen-komponen di dalamnya
func Legacy_MultipleBfs(target string, combinations map[Pair]string, revCombinations map[string][]Pair, tierMap map[string]int, maxRecipes int) MultipleRecipesResult {
  // Pertama, cari resep awal pake ShortestBfsFiltered
  firstRecipe := ShortestBfs(target, combinations, tierMap)
  
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

  // Track recipes we've already seen to avoid duplicates
  // Pake map untuk nyimpen resep yang udah kita temuin, biar gak duplikat
  seenRecipes := make(map[string]bool)
  
  // Mark first recipe as seen
  seenRecipes[RecipeToString(firstRecipe, target)] = true
  
  // Queue for BFS - we'll store recipe variations with the element we're focusing on
  // Queue untuk BFS - kita simpan variasi resep beserta elemen yang lagi kita fokuskan
  type QueueItem struct {
    Recipe      map[string]Element  // Resep saat ini
    FocusElem   string              // Elemen yang lagi kita coba variasikan
  }
  
  queue := []QueueItem{}
  
  // First add target variations
  // Pertama, tambahkan variasi untuk target
  queue = append(queue, QueueItem{Recipe: firstRecipe, FocusElem: target})
  
  // Then add component variations - cari variasi untuk semua komponen non-base
  // Kemudian tambahkan semua elemen non-dasar ke queue untuk divariasikan
  for elem := range firstRecipe {
    if !isBaseElement(elem) && elem != target {
      queue = append(queue, QueueItem{Recipe: firstRecipe, FocusElem: elem})
    }
  }
  
  // BFS to explore different recipe variations at all levels
  // BFS untuk menjelajahi variasi resep di semua level
  for len(queue) > 0 && (maxRecipes <= 0 || len(recipes) < maxRecipes) {
    // Get next item from queue
    current := queue[0]
    queue = queue[1:]
    
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
    // Dapatkan semua cara valid untuk membuat elemen ini
    validPairs := filterValidPairs(revCombinations[focusElem], focusElem, tierMap)
    
    // Try each alternative way to make this element
    // Coba setiap alternatif cara membuat elemen ini
    for _, pair := range validPairs {
      // Skip the current recipe for this element
      if (pair.First == originalSources.Source && pair.Second == originalSources.Partner) ||
         (pair.First == originalSources.Partner && pair.Second == originalSources.Source) {
        continue
      }
      
      // Create a variation with this alternative
      // Buat variasi resep dengan alternatif ini
      variation := copyRecipe(currentRecipe)
      variation[focusElem] = Element{Source: pair.First, Partner: pair.Second}
      
      // Ensure the new ingredients have valid recipes if they're not already in our recipe
      // Pastikan ingredient baru punya resep valid kalo belum ada di resep kita
      allValid := true
      for _, ingredient := range []string{pair.First, pair.Second} {
        if isBaseElement(ingredient) {
          continue // Base elements are always valid
        }
        
        // If we don't have a recipe for this ingredient yet, find one
        if _, exists := variation[ingredient]; !exists {
          ingredientRecipe := findIngredientRecipe(ingredient, combinations, revCombinations, tierMap, visited)
          if len(ingredientRecipe) == 0 {
            allValid = false
            break
          }
          
          // Add the ingredient's recipe to our variation
          for elem, sources := range ingredientRecipe {
            if _, exists := variation[elem]; !exists {
              variation[elem] = sources
              visited[elem] = true // Mark as visited
            }
          }
        }
      }
      
      if !allValid {
        continue // Skip this variation if we couldn't complete it
      }
      
      // Check if this is a unique recipe
      // Cek apakah ini resep unik yang belum pernah kita temuin
      recipeStr := RecipeToString(variation, target)
      if !seenRecipes[recipeStr] {
        seenRecipes[recipeStr] = true
        recipes = append(recipes, variation)
        
        // Check if we have enough recipes
        if maxRecipes > 0 && len(recipes) >= maxRecipes {
          break
        }
        
        // Add variations for each component in our recipe (BFS approach)
        // Tambahkan variasi untuk setiap komponen dalam resep (pendekatan BFS)
        for elem := range variation {
          if !isBaseElement(elem) {
            queue = append(queue, QueueItem{
              Recipe:    variation,
              FocusElem: elem,
            })
          }
        }
      }
    }
  }
  
  return MultipleRecipesResult{
    Recipes:   recipes,
    NodeCount: len(visited),
  }
}