package util
import "hash/fnv"

func HashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// func BuildHash(inventory map[string]bool) uint64 {
// 	var hash uint64 = 0
// 	for elem := range inventory {
// 		hash ^= HashString(elem)
// 	}
// 	return hash
// }

func BuildHashFromSlice(arr []string) uint64 {
	var hash uint64 = 0
	for _, elem := range arr {
		hash ^= HashString(elem)
	}
	return hash
}