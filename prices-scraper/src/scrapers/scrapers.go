package scrapers

import (
	"database/sql"
	"fmt"
	"items-scraper/src/models"
	"items-scraper/src/repository"
	"items-scraper/src/utils"
	"sync"
)

func ScapeAll(items []models.Item, db *sql.DB) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ScrapeMercadoLivre(items, db)
	}()

	go func() {
		defer wg.Done()
		ScrapeAmazon(items, db)
	}()

	wg.Wait()
	fmt.Println("✓ All scrapers completed!")

}

// TODO: Refactor this
func InsertProductPrice(item models.Item, product models.Product, productRepo *repository.ProductRepository) {
	priceHistory := CreatePriceHistory(product, item, productRepo)
	if priceHistory.ProductID != -1 {
		_, err := productRepo.InsertPriceHistory(priceHistory)
		utils.LogErr(err)
	}
}

// TODO: Refactor this
func CreatePriceHistory(p models.Product, item models.Item, repository *repository.ProductRepository) models.PriceHistory {
	id, err := repository.FindByIDProduct(p.IDProduct)
	utils.CheckErr(err)

	if repository.CanInsertProduct(p, id) {
		id, err = repository.InsertProduct(p)
		utils.CheckErr(err)
	}

	priceHistory := models.PriceHistory{
		ProductID: id,
		Price:     item.Price,
		Currency:  item.Currency,
	}
	return priceHistory
}
