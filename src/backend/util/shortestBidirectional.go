package util

// ShortestBidirectional implementasi algoritma pencarian bidirectional (dua arah)
// yang mencari jalur terpendek untuk membuat elemen target dengan
// menjalankan BFS dari elemen dasar (maju) dan dari target (mundur) secara bersamaan
func ShortestBidirectional(target string, combinations map[Pair]string, revCombinations map[string][]Pair, tierMap map[string]int) map[string]Element {
	// Siapin queue untuk arah maju (dari elemen dasar)
	forwardQueue := make([]string, len(BaseElements))
	copy(forwardQueue, BaseElements)

	// Siapin queue untuk arah mundur (dari target)
	backwardQueue := []string{target}

	// Tandain elemen yang udah dilihat di arah maju
	forwardSeen := make(map[string]bool, len(BaseElements))
	for _, b := range BaseElements {
		forwardSeen[b] = true
	}

	// Tandain elemen yang udah dilihat di arah mundur
	backwardSeen := make(map[string]bool)
	backwardSeen[target] = true

	// Simpan resep untuk setiap elemen yang dihasilkan dari arah maju
	forwardRecipes := make(map[string]Element)

	// Simpan elemen apa aja yang bisa bikin suatu elemen dari arah mundur
	// Key: elemen yang ditemukan dari arah mundur
	// Value: pasangan bahan yang bisa membuat elemen tersebut
	backwardRecipes := make(map[string][]Pair)

	// Elemen yang jadi titik temu kedua arah pencarian
	var meetingPoint string
	
	// Flag untuk tracking apakah sudah ketemu
	found := false

	// Loop sampai ketemu titik temu atau salah satu queue kosong
	for len(forwardQueue) > 0 && len(backwardQueue) > 0 && !found {
		// ===== FORWARD SEARCH (dari elemen dasar ke target) =====
		// Jalanin satu langkah BFS dari arah maju
		if len(forwardQueue) > 0 {
			// Ambil elemen dari queue
			current := forwardQueue[0]
			forwardQueue = forwardQueue[1:]

			// Cek apakah elemen ini juga ada di backwardSeen (berarti ketemu titik temu)
			if backwardSeen[current] {
				meetingPoint = current
				found = true
				break
			}

			// Ambil tier elemen saat ini
			currentTier := tierMap[current]

			// Coba kombinasiin dengan elemen yang udah diketahui
			for partner := range forwardSeen {
				// Ambil tier partner
				partnerTier := tierMap[partner]
				
				// Coba bikin produk dari pasangan ini
				pair := Pair{First: current, Second: partner}
				if product, exists := combinations[pair]; exists {
					// Ambil tier produk
					productTier := tierMap[product]
					
					// Cek tier constraint: produk harus dari tier lebih tinggi
					if currentTier < productTier && partnerTier < productTier {
						// Kalo produk baru, tambahin ke queue maju
						if !forwardSeen[product] {
							forwardSeen[product] = true
							forwardRecipes[product] = Element{Source: current, Partner: partner}
							forwardQueue = append(forwardQueue, product)
						}
					}
				}
			}
		}

		// ===== BACKWARD SEARCH (dari target ke elemen dasar) =====
		// Jalanin satu langkah BFS dari arah mundur
		if len(backwardQueue) > 0 {
			// Ambil elemen dari queue
			current := backwardQueue[0]
			backwardQueue = backwardQueue[1:]

			// Cek apakah elemen ini juga ada di forwardSeen (berarti ketemu titik temu)
			if forwardSeen[current] {
				meetingPoint = current
				found = true
				break
			}

			// Ambil tier elemen saat ini
			currentTier := tierMap[current]

			// Cek semua pasangan yang bisa menghasilkan elemen ini
			for _, pair := range revCombinations[current] {
				// Ambil tier kedua bahan
				firstTier := tierMap[pair.First]
				secondTier := tierMap[pair.Second]
				
				// Cek tier constraint: bahan harus dari tier lebih rendah
				if firstTier < currentTier && secondTier < currentTier {
					// Catat pasangan ini sebagai pembuat elemen current
					backwardRecipes[current] = append(backwardRecipes[current], pair)
					
					// Tambahin kedua bahan ke queue mundur kalo belum pernah dilihat
					if !backwardSeen[pair.First] && !isBaseElement(pair.First) {
						backwardSeen[pair.First] = true
						backwardQueue = append(backwardQueue, pair.First)
					}
					
					if !backwardSeen[pair.Second] && !isBaseElement(pair.Second) {
						backwardSeen[pair.Second] = true
						backwardQueue = append(backwardQueue, pair.Second)
					}
				}
			}
		}
	}

	// Kalo gak ketemu titik temu, return kosong
	if !found {
		return make(map[string]Element)
	}

	// Bentuk resep lengkap dari titik temu
	// Kalo titik temu adalah target, kita udah selesai
	if meetingPoint == target {
		return forwardRecipes
	}

	// Kalo titik temu bukan target, kita perlu gabungin resep maju dan mundur
	result := make(map[string]Element)
	
	// Salin semua resep dari arah maju ke hasil
	for elem, recipe := range forwardRecipes {
		result[elem] = recipe
	}
	
	// Mulai dari titik temu, buat resep untuk semua elemen di jalur mundur
	completePath := completeBackwardPath(meetingPoint, target, backwardRecipes, forwardRecipes, revCombinations, tierMap)
	for elem, recipe := range completePath {
		result[elem] = recipe
	}
	
	return result
}

// completeBackwardPath menyelesaikan jalur mundur dari titik temu ke target
// dengan memastikan kita punya resep valid untuk semua elemen di jalur
func completeBackwardPath(meetingPoint, target string, backwardRecipes map[string][]Pair, 
						 forwardRecipes map[string]Element, revCombinations map[string][]Pair, 
						 tierMap map[string]int) map[string]Element {
	result := make(map[string]Element)
	
	// Buat rekonstruksi resep dari titik temu ke target
	toProcess := []string{target}
	processed := make(map[string]bool)
	
	// Tandai titik temu sebagai sudah diproses (kita udah punya resepnya)
	processed[meetingPoint] = true
	
	for len(toProcess) > 0 {
		current := toProcess[0]
		toProcess = toProcess[1:]
		
		if processed[current] {
			continue
		}
		
		// Kita perlu cari resep untuk current
		if pairs, exists := backwardRecipes[current]; exists && len(pairs) > 0 {
			// Ambil pasangan bahan pertama sebagai resep
			pair := pairs[0]
			result[current] = Element{Source: pair.First, Partner: pair.Second}
			
			// Tambahkan kedua bahan ke antrian untuk diproses
			// kecuali jika sudah ada di forwardRecipes atau sudah diproses
			if !processed[pair.First] && forwardRecipes[pair.First].Source == "" {
				toProcess = append(toProcess, pair.First)
			}
			
			if !processed[pair.Second] && forwardRecipes[pair.Second].Source == "" {
				toProcess = append(toProcess, pair.Second)
			}
		} else {
			// Kalo gak ada di backwardRecipes, coba cari resep dengan ShortestDfs
			// Ini bisa terjadi karena kita melompati beberapa elemen dalam pencarian mundur
			miniResult := ShortestDfs(current, revCombinations, tierMap)
			
			// Gabungkan dengan hasil kita
			for elem, recipe := range miniResult {
				if !processed[elem] {
					result[elem] = recipe
					processed[elem] = true
				}
			}
		}
		
		processed[current] = true
	}
	
	return result
}
