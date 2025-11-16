package main

import (
	"items-scraper/src/config"
	"items-scraper/src/models"
	"items-scraper/src/scrapers"
	"items-scraper/src/utils"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	db, err := cfg.Connect()
	utils.CheckErr(err)

	items := []models.Item{
		{ID: "1", Title: "Switch 2", MinPrice: 2200, MaxPrice: 9999},
		{ID: "2", Title: "Steam Deck", MinPrice: 2200, MaxPrice: 9999},
	}

	scrapers.ScapeAll(items, db)

	defer func() {
		err := db.Close()
		utils.LogErr(err)
	}()
}
