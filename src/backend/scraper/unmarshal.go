package scraper

import (
	"encoding/json"
	"os"

	"backend/util"
)

// UnmarshalRecipes loads recipe data from a JSON file and populates the specified maps
// Returns the populated maps and any error that occurred
func UnmarshalRecipes(filename string) (map[util.Pair]string, map[string][]util.Pair, map[string]int, error) {
	// Initialize the maps
	rawRecipe := make(map[util.Pair]string)
	reversedRawRecipe := make(map[string][]util.Pair)
	ingredientsTier := make(map[string]int)

	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, err
	}

	// Unmarshal the JSON data into a slice of RecipeJSON objects
	var recipes []RecipeJSON
	err = json.Unmarshal(data, &recipes)
	if err != nil {
		return nil, nil, nil, err
	}

	// Populate the maps
	for _, recipe := range recipes {
		// Add the tier information
		ingredientsTier[recipe.Result] = recipe.Tier

		// Add the combinations to rawRecipe and reversedRawRecipe
		if recipe.Combinations != nil {
			for _, combo := range recipe.Combinations {
				// Set the mapping from pair to result
				rawRecipe[combo] = recipe.Result

				// Add to reversed map
				reversedRawRecipe[recipe.Result] = append(reversedRawRecipe[recipe.Result], combo)
			}
		}
	}

	return rawRecipe, reversedRawRecipe, ingredientsTier, nil
}