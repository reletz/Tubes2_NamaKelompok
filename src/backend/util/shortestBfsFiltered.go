package util

// ShortestBfsFiltered implementasi algoritma BFS dengan batasan tingkatan
// buat nyari jalur terpendek bikin elemen target.
// Cuma mempertimbangkan kombinasi dimana kedua bahan dari tier lebih rendah dari produk.
func ShortestBfsFiltered(target string, combinations map[Pair]string, tierMap map[string]int) map[string]Element {
	// Siapin queue dengan elemen dasar
	queue := make([]string, len(BaseElements))
	copy(queue, BaseElements)

	// Tandain elemen yang udah dilihat
	seen := make(map[string]bool, len(BaseElements))
	for _, b := range BaseElements {
		seen[b] = true
	}

	// Simpan resep untuk setiap elemen yang dihasilkan
	prev := make(map[string]Element)

	// Loop BFS
	for i := 0; i < len(queue); i++ {
		current := queue[i]

		// Kalo udah ketemu target, berhenti pencarian
		if current == target {
			break
		}

		// Ambil tier elemen saat ini
		currentTier := tierMap[current]

		// Coba kombinasiin elemen saat ini dengan semua elemen yang udah dilihat
		for partner := range seen {
			// Ambil tier partner
			partnerTier := tierMap[partner]
			
			// Coba bikin produk dari pasangan ini
			pair := Pair{First: current, Second: partner}
			if product, exists := combinations[pair]; exists {
				// Ambil tier produk
				productTier := tierMap[product]
				
				// Cek apakah ini kombinasi tier yang valid:
				// Produk harus tier lebih tinggi dari kedua bahan
				if currentTier < productTier && partnerTier < productTier {
					// Kalo produk baru (belum pernah dilihat), tambahin ke queue
					if !seen[product] {
						seen[product] = true
						prev[product] = Element{Source: current, Partner: partner}
						queue = append(queue, product)
					}
				}
			}
		}
	}

	return prev
}