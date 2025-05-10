package util

// Element nyatet dua bahan ("Source" dan "Partner") yang pertama kali
// digabung buat bikin produk.
// Source = bahan yang lagi kita proses (dikeluarin dari queue)
// Partner = bahan lain (dari set yang udah dilihat) yang kita gabungin sama Source buat bikin produk.


// Elemen dasar di Little Alchemy 2
var BaseElements = []string{"Air", "Earth", "Fire", "Water"}

// ShortestBfs jalanin BFS mulai dari elemen dasar
// sampai nemuin target (atau habis opsi). Hasilnya adalah map
// dari produk ke Element (siapa Source & Partner yang ngasilin itu),
// jadi kita bisa nyusun lagi jalur resepnya nanti.
func ShortestBfs(target string, combinations map[Pair][]string) map[string]Element {
	// 1) Siapin queue yang berisi elemen dasar
	queue := make([]string, len(BaseElements))
	copy(queue, BaseElements)

	// 2) seen = set buat tandain elemen yang udah kita lihat
	seen := make(map[string]bool, len(BaseElements))
	for _, b := range BaseElements {
		seen[b] = true
	}

	// 3) prev bakal nyimpen, buat tiap produk, Element{Source, Partner}
	//    yang pertama kali nghasilin produk itu
	prev := make(map[string]Element)

	// 4) Loop BFS: selama masih ada elemen di queue (i < len(queue))
	for i := 0; i < len(queue); i++ {
		current := queue[i]

		// Kalo current udah sama dengan target, kita berhenti nyari
		if current == target {
				break
		}

		// 5) Untuk tiap bahan "partner" yang udah kita lihat,
		//    coba gabungin current + partner, pake combinations[pair] -> produk
		for partner := range seen {
			// Cari produk apa aja yang bisa dibikin dengan Pair{A: current, B: partner}
			for _, product := range combinations[Pair{First: current, Second: partner}] {
				// 6) Kalo produk baru (belum pernah dilihat), tandai dan masukin queue
				if !seen[product] {
					seen[product] = true
					// catet siapa Source & partner-nya
					prev[product] = Element{Source: current, Partner: partner}
					// masukin ke queue buat diproses nanti
					queue = append(queue, product)
				}
			}
		}
	}

	// 7) Kembalikan map prev. Dari sini kita bisa bikin jalur resep:
	//    mulai dari target, lihat prev[target], terus lihat prev[Source], dst.
	return prev
}