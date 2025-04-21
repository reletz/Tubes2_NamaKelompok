package main

import (
	"backend/scraper"
	"fmt"
)


func main(){
	fmt.Println("Halo!");
	scraper.Scraper();
	scraper.MapScraper();
}