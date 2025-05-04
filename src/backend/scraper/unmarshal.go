package scraper

// import (
// 	"encoding/json"
// 	"os"

// 	"backend/util"
// )

// // UnmarshalRecipes loads JSON file and returns map[util.Pair]string
// func UnmarshalRecipes(filename string) (map[util.Pair]string, error) {
// 	data, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var recipes []RecipeJSON
// 	err = json.Unmarshal(data, &recipes)
// 	if err != nil {
// 		return nil, err
// 	}

// 	lookup := make(map[util.Pair]string)
// 	for _, recipe := range recipes {
// 		pair := util.Pair{
// 			First:  recipe.First,
// 			Second: recipe.Second,
// 		}
// 		lookup[pair] = recipe.Result
// 	}

// 	return lookup, nil
// }