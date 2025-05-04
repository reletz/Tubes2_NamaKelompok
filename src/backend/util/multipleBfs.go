package util

import (
  "time"
  "sync"
)

// ConstrainedBfs melakukan BFS tetapi menghindari kombinasi tertentu
func ConstrainedBfs(target string, combinations map[Pair][]string, forbidden map[Pair]bool) map[string]Element {
  queue := make([]string, len(BaseElements))
  copy(queue, BaseElements)
  
  seen := make(map[string]bool, len(BaseElements))
  for _, b := range BaseElements {
    seen[b] = true
  }
  
  prev := make(map[string]Element)
  
  for i := 0; i < len(queue); i++ {
    current := queue[i]
    
    if current == target {
      break
    }
    
    for partner := range seen {
      pair := Pair{First: current, Second: partner}
      reversePair := Pair{First: partner, Second: current}
      
      // Lewati kombinasi yang dilarang
      if forbidden[pair] || forbidden[reversePair] {
        continue
      }
      
      for _, product := range combinations[pair] {
        if !seen[product] {
          seen[product] = true
          prev[product] = Element{Source: current, Partner: partner}
          queue = append(queue, product)
        }
      }
    }
  }
  
  return prev
}

// MultipleRecipeBFS mencari beberapa resep berbeda untuk satu target
func MultipleRecipeBFS(target string, combinations map[Pair][]string, maxRecipes int) []ResultJSON {
  var results []ResultJSON
  var mutex sync.Mutex
  var wg sync.WaitGroup
  
  // Pake jalur terpendek
  standardPath := ShortestBfs(target, combinations)
  if _, exists := standardPath[target]; exists {
    results = append(results, BuildTree(target, standardPath))
  }
  
  // Berhenti jika hanya perlu satu resep
  if maxRecipes <= 1 {
    return results
  }
  
  // Buat map buat track resep unik
  uniqueRecipes := make(map[string]bool)
  if len(results) > 0 {
    uniqueRecipes[treeSignature(results[0])] = true
  }
  
  // Channel untuk hasil
  resultChan := make(chan ResultJSON, maxRecipes*3)
  
  // Hasilkan BFS dengan batasan untuk menemukan resep yang beragam
  wg.Add(1)
  go func() {
    defer wg.Done()
    
    //Hindarin kombinasi yang terpendek
    forbidden := make(map[Pair]bool)
    current := target
    for {
      element, exists := standardPath[current]
      if !exists {
        break
      }
      
      pair1 := Pair{First: element.Source, Second: element.Partner}
      pair2 := Pair{First: element.Partner, Second: element.Source}
      forbidden[pair1] = true
      forbidden[pair2] = true
      
      current = element.Source
    }
    
    // Jalankan BFS dengan batasan menggunakan set batasan ini
    path := ConstrainedBfs(target, combinations, forbidden)
    if _, exists := path[target]; exists {
      tree := BuildTree(target, path)
      resultChan <- tree
    }
  }()
  
  // Tutup channel ketika semua goroutine selesai
  go func() {
    wg.Wait()
    close(resultChan)
  }()
  
  // Kumpulkan hasil
  timeout := time.After(5 * time.Second)
  
collectLoop:
  for len(results) < maxRecipes {
    select {
    case result, ok := <-resultChan:
      if !ok {
        break collectLoop
      }
      
      signature := treeSignature(result)
      
      mutex.Lock()
      if !uniqueRecipes[signature] {
        uniqueRecipes[signature] = true
        results = append(results, result)
      }
      mutex.Unlock()
      
    case <-timeout:
      // Batas waktu tercapai
      break collectLoop
    }
  }
  
  return results
}

func treeSignature(tree ResultJSON) string {
  // Gabungin semua nama child
  signature := tree.Name
  
  // Ekstrak dan urutkan nama child
  childNames := make([]string, 0)
  
  if children, ok := tree.Children.([]interface{}); ok {
    for _, child := range children {
      // Coba ekstrak bidang nama dari child
      if childMap, ok := child.(map[string]interface{}); ok {
        if name, ok := childMap["name"].(string); ok {
          childNames = append(childNames, name)
        }
      }
    }
  } else if children, ok := tree.Children.([]ResultJSON); ok {
    // Pemeriksaan tipe alternatif jika Children adalah slice of ResultJSON
    for _, child := range children {
      childNames = append(childNames, child.Name)
    }
  } else if children, ok := tree.Children.([]Child); ok {
    // Kemungkinan tipe lain untuk Children
    for _, child := range children {
      childNames = append(childNames, child.Name)
    }
  }
  
  // Urutkan child-child untuk tanda tangan yang konsisten
  for i := 0; i < len(childNames); i++ {
    for j := i + 1; j < len(childNames); j++ {
      if childNames[i] > childNames[j] {
        childNames[i], childNames[j] = childNames[j], childNames[i]
      }
    }
  }
  
  // Tambahkan nama child yang sudah diurutkan ke tanda tangan
  for _, name := range childNames {
    signature += ":" + name
  }
  
  return signature
}

type Child struct {
  Name string `json:"name"`
}