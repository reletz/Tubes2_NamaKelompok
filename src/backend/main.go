package main

import (
	"backend/scraper"
	"backend/util"
	"fmt"
)


func main(){
	fmt.Println("Halo!");
	var recipes map[util.Pair]string;
	scraper.Scraper(&recipes, false);
	// scraper.GraphScraper();
}