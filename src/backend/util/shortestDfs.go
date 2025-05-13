package util

// NodeState buat ngetracking status eksplorasi tiap elemen
type NodeState struct {
	CurrentPairIndex int
	Pairs            []Pair
	ValidPairs       []Pair
	Visited          bool
}

func ShortestDfs(target string, revCombinations map[string][]Pair, tierMap map[string]int) map[string]Element {
  // Inisialisasi map hasil: elemen -> resepnya
  result := make(map[string]Element)
  
  // Cek apakah target ada di kombinasi
  if _, exists := revCombinations[target]; !exists {
    return result
  }
  
  // Map buat ngetracking status eksplorasi tiap elemen
  nodeStates := make(map[string]*NodeState)
  
  // Tambahin elemen dasar ke node states dengan visited=true
  for _, elem := range BaseElements {
    nodeStates[elem] = &NodeState{Visited: true}
    result[elem] = Element{} // Tandain elemen dasar dengan resep kosong
  }
  
  // Tracking elemen yang lagi kita coba resolve
  inProgress := make(map[string]bool)
  
  // Pake fungsi rekursif sebagai helper buat DFS
  var explore func(element string) bool
  explore = func(element string) bool {
    // Elemen dasar udah selesai
    if isBaseElement(element) {
      return true
    }
    
    // Skip kalo kita udah nemu solusi buat elemen ini
    if state := nodeStates[element]; state != nil && state.Visited {
      return true
    }
    
    // Deteksi siklus - kalo kita udah coba resolve elemen ini di jalur saat ini
    if inProgress[element] {
      return false
    }
    
    // Tandain sebagai sedang diproses
    inProgress[element] = true
    defer func() { inProgress[element] = false }()
    
    // Ambil atau bikin state node
    state := nodeStates[element]
    if state == nil {
      pairs := revCombinations[element]
      validPairs := filterValidPairs(pairs, element, tierMap)
      state = &NodeState{
        CurrentPairIndex: 0,
        Pairs:            pairs,
        ValidPairs:       validPairs,
        Visited:          false,
      }
      nodeStates[element] = state
    }
    
    // Coba tiap resep yang valid
    for i := 0; i < len(state.ValidPairs); i++ {
      pair := state.ValidPairs[i]
      
      // Catat resep ini sementara
      result[element] = Element{
        Source:  pair.First,
        Partner: pair.Second,
      }
      
      // Coba resolve kedua bahan
      firstResolved := isBaseElement(pair.First) || explore(pair.First)
      if !firstResolved {
        continue // Coba resep berikutnya kalo bahan pertama gak bisa diresolved
      }
      
      secondResolved := isBaseElement(pair.Second) || explore(pair.Second)
      if !secondResolved {
        continue // Coba resep berikutnya kalo bahan kedua gak bisa diresolved
      }
      
      // Kedua bahan resolved - kita nemu resep valid
      state.Visited = true
      state.CurrentPairIndex = i + 1 // Inget resep mana yang kita pake
      return true
    }
    
    // Kalo sampe sini, gak ada resep valid yang ketemu
    delete(result, element) // Hapus resep sementara
    return false
  }
  
  // Mulai eksplorasi dari target
  explore(target)
  
  // Bersihin map hasil - hapus entri dengan resep kosong yang bukan elemen dasar
  for key, elem := range result {
    if !isBaseElement(key) && (elem.Source == "" || elem.Partner == "") {
      delete(result, key)
    }
  }
  
  return result
}

// Improved filterValidPairs to prioritize lower tier ingredients
func filterValidPairs(pairs []Pair, element string, tierMap map[string]int) []Pair {
    currentTier := tierMap[element]
    var validPairs []Pair
    
    for _, pair := range pairs {
        // Only consider pairs where both ingredients are from lower tier
        // This prevents using equal or higher tier elements to create lower tier ones
        if tierMap[pair.First] < currentTier && tierMap[pair.Second] < currentTier {
            validPairs = append(validPairs, pair)
        }
    }
    
    return validPairs
}