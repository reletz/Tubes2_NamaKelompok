package main

import (
	"backend/scraper"
	"backend/util"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SearchRequest struct {
	NamaResep      string `json:"namaResep"`
	MaksimalResep  int    `json:"maksimalResep"`
	Algoritma      string `json:"algoritma"`
	ModePencarian  string `json:"modePencarian"`
}

type TreeNode struct {
	Name     string     `json:"name"`
	Children []TreeNode `json:"children,omitempty"`
}

type TreeResponse struct {
	TreeData    []TreeNode `json:"treeData"`
	TimeTaken   string     `json:"timetaken"`
	NodeVisited int        `json:"node_visited"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("❌ JSON decode error:", err)
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	log.Printf("Received search: NamaResep=%s, MaksimalResep=%d, Algoritma=%s, ModePencarian=%s",
		req.NamaResep, req.MaksimalResep, req.Algoritma)

	// scrape
	rawRecipe := make(map[util.Pair]string)
	reversedRawRecipe := make(map[string][]util.Pair)
	ingredientsTier := make(map[string]int)
	scraper.Scraper(rawRecipe, ingredientsTier, reversedRawRecipe, true)

	var result util.MultipleRecipesResult
	start := time.Now()

	switch req.Algoritma {
	case "BFS":
		result = util.MultipleBfs(req.NamaResep, rawRecipe, reversedRawRecipe, ingredientsTier, req.MaksimalResep, 4)
	case "DFS":
		result = util.MultipleDfs(req.NamaResep, reversedRawRecipe, ingredientsTier, req.MaksimalResep, 4)
	case "Bi-BFS":
		result = util.MultipleBidirectional(req.NamaResep, rawRecipe, reversedRawRecipe, ingredientsTier, req.MaksimalResep, 4)
	default:
		http.Error(w, "Unsupported algorithm", http.StatusBadRequest)
		return
	}

	elapsed := time.Since(start)

	trees, nodeVisited := util.BuildMultipleTrees(req.NamaResep, result)

	// nyesuaiin jumlah resep
	if len(trees) > req.MaksimalResep {
		trees = trees[:req.MaksimalResep]
	}

	// convert util.Node to TreeNode
	var convert func(n *util.Node) TreeNode
	convert = func(n *util.Node) TreeNode {
		children := []TreeNode{}
		for _, c := range n.Children {
			children = append(children, convert(c))
		}
		return TreeNode{
			Name:     n.Name,
			Children: children,
		}
	}

	var treeData []TreeNode
	for _, t := range trees {
		treeData = append(treeData, convert(t))
	}

	response := TreeResponse{
		TreeData:    treeData,
		TimeTaken:   elapsed.String(),
		NodeVisited: nodeVisited,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/search", searchHandler)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}