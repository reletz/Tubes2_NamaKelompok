package util

// RecipeToString menghasilkan representasi string unik dari sebuah resep
// Berfungsi sebagai "fingerprint" resep untuk deteksi duplikat
func RecipeToString(recipe map[string]Element, target string) string {
	// Track all elements we've seen so far
	processed := make(map[string]bool)
	var result string
	
	// Recursive function to build string representation
	// Fungsi rekursif untuk membuat representasi string dari resep,
	// mempertimbangkan struktur lengkap (bukan cuma ingredient top-level)
	var processElement func(elem string)
	processElement = func(elem string) {
		if processed[elem] || isBaseElement(elem) {
			return
		}
		
		processed[elem] = true
		sources := recipe[elem]
		
		// Normalize ingredient order
		first, second := sources.Source, sources.Partner
		if first > second {
			first, second = second, first
		}
		
		// Add this element's recipe to the string
		result += elem + ":" + first + "+" + second + "|"
		
		// Process ingredients recursively
		if !isBaseElement(first) {
			processElement(first)
		}
		if !isBaseElement(second) {
			processElement(second)
		}
	}
	
	// Start with the target
	processElement(target)
	return result
}

// NormalizeIngredients menormalkan urutan dua bahan sehingga A+B = B+A
func NormalizeIngredients(a, b string) string {
	if a > b {
		return b + "+" + a
	}
	return a + "+" + b
}


// Helper function untuk mengecek apakah sebuah resep unik dibandingkan dengan daftar resep yang ada
// Perlu memastikan bahwa kombinasi A+B dianggap sama dengan B+A
func isUniqueRecipe(recipe map[string]Element, existingRecipes []map[string]Element) bool {
  for _, existing := range existingRecipes {
    // Perlu tracking elemen yang sudah dibandingkan
    match := true
    
    // Bandingkan semua elemen di resep dengan resep yang sudah ada
    for elem, sources := range recipe {
      if existingSources, hasElem := existing[elem]; hasElem {
        // Normalisasi urutan bahan untuk perbandingan yang adil
        recipeFirst, recipeSecond := sources.Source, sources.Partner
        existingFirst, existingSecond := existingSources.Source, existingSources.Partner
        
        // Urutkan bahan berdasarkan string agar A+B = B+A
        if recipeFirst > recipeSecond {
          recipeFirst, recipeSecond = recipeSecond, recipeFirst
        }
        if existingFirst > existingSecond {
          existingFirst, existingSecond = existingSecond, existingFirst
        }
        
        // Bandingkan dengan urutan yang sudah dinormalisasi
        if recipeFirst != existingFirst || recipeSecond != existingSecond {
          match = false
          break
        }
      } else {
        // Elemen tidak ada di resep existing
        match = false
        break
      }
    }
    
    // Pastikan kedua resep memiliki jumlah elemen yang sama
    if match && len(recipe) == len(existing) {
      return false // Ditemukan resep yang sama
    }
  }
  
  return true // Tidak ada resep yang sama
}

// findIngredientRecipe mencari resep valid untuk suatu ingredient
// Fungsi ini memastikan kita bisa membuat ingredient dengan aturan tiering
func findIngredientRecipe(ingredient string, combinations map[Pair]string, revCombinations map[string][]Pair, tierMap map[string]int, visited map[string]bool) map[string]Element {
  // Kalo udah elemen dasar, gak perlu resep
  if isBaseElement(ingredient) {
    return map[string]Element{}
  }
  
  // Cari semua pasangan valid berdasarkan tier
  pairs := revCombinations[ingredient]
  result := make(map[string]Element)
  
  // Coba setiap pasangan valid
  for _, pair := range pairs {
    // Cek validitas tier
    productTier := tierMap[ingredient]
    if tierMap[pair.First] >= productTier || tierMap[pair.Second] >= productTier {
      continue
    }
    
    // Catat resep untuk ingredient ini
    result[ingredient] = Element{Source: pair.First, Partner: pair.Second}
    
    // Coba cari resep untuk tiap bahan
    validRecipe := true
    
    for _, source := range []string{pair.First, pair.Second} {
      if isBaseElement(source) {
        continue
      }
      
      // Cari resep untuk bahan secara rekursif
      sourceRecipe := findIngredientRecipe(source, combinations, revCombinations, tierMap, visited)
      if len(sourceRecipe) == 0 {
        validRecipe = false
        break
      }
      
      // Tambahkan resep bahan ke hasil
      for elem, elemSources := range sourceRecipe {
        result[elem] = elemSources
      }
    }
    
    if validRecipe {
      // Update elemen yang udah dikunjungi
      for elem := range result {
        visited[elem] = true
      }
      return result
    }
  }
  
  return map[string]Element{}
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

// Helper function to check if an element is a base element
func isBaseElement(element string) bool {
	BaseElements := map[string]bool{
		"Fire": true,
		"Water": true,
		"Air": true,
		"Earth": true,
	}
	return BaseElements[element]
}

// MultipleRecipesResult buat nyimpen hasil dari MultipleRecipesDfs
type MultipleRecipesResult struct {
  Recipes   []map[string]Element // Kumpulan resep yang valid
  NodeCount int                  // Jumlah node/elemen yang dikunjungi
}