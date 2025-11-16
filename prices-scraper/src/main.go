package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"items-scraper/src/models"
	"items-scraper/src/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	host     string
	port     int
	user     string
	password string
	dbname   string
)

func main() {
	setupDatabase()
	db := connectDB()

	items := []models.Item{
		{ID: "1", Title: "Switch 2", MinPrice: 2200, MaxPrice: 9999},
		{ID: "2", Title: "Steam Deck", MinPrice: 2200, MaxPrice: 9999},
	}

	scapeAll(items, db)

	defer func() {
		err := db.Close()
		utils.LogErr(err)
	}()
}

func setupDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[ERROR]: .env file not found, using environment variables or defaults")
		log.Println(err)
		panic(err)
	}

	host = getEnv("DB_HOST")
	port = getEnvAsInt("DB_PORT", 5432)
	user = getEnv("DB_USER")
	password = getEnv("DB_PASSWORD")
	dbname = getEnv("DB_NAME")
}

func connectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func scapeAll(items []models.Item, db *sql.DB) {
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

func findProductByIdProduct(db *sql.DB, id_product string) (int, error) {
	sqlQuery := "SELECT * FROM products WHERE id_product = $1"
	rows, err := db.Query(sqlQuery, id_product)
	utils.CheckErr(err)

	var id int
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.IDProduct,
			&product.ProductName,
			&product.URL,
			&product.Retailer,
			&product.CreatedAt,
		)
		if err != nil {
			return -1, err
		}
		id = product.ID
	}
	return id, nil
}

func insertProduct(db *sql.DB, product models.Product) (int, error) {
	sqlStatement := `
		INSERT INTO products (id_product, product_name, url, retailer)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id int
	err := db.QueryRow(
		sqlStatement,
		product.IDProduct,
		product.ProductName,
		product.URL,
		product.Retailer,
	).Scan(&id)

	if err != nil {
		fmt.Printf("[ERROR] inserting product: %v\n", err)
		return 0, err
	}

	fmt.Printf("Successfully inserted product with ID: %d\n", id)

	return id, err
}

func insertPriceHistory(db *sql.DB, ph models.PriceHistory) (int, error) {
	if !canInsertPriceHistory(db, ph) {
		err := errors.New("duplicate entry: this item has already been inserted today")
		fmt.Printf("[ERROR] %v\n", err)
		return 0, err
	}

	sqlQuery := `
		INSERT INTO price_history (product_id, price, currency, price_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id int
	err := db.QueryRow(
		sqlQuery,
		ph.ProductID,
		ph.Price,
		ph.Currency,
		time.Now(),
	).Scan(&id)

	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return 0, err
	}

	fmt.Printf("Successfully inserted price_history with ID: %d\n", id)
	return id, nil

}

func canInsertPriceHistory(db *sql.DB, ph models.PriceHistory) bool {
	sqlQuery := `
		select ph.id
		from price_history ph 
		where 1=1
			and ph.product_id = $1
			and ph.price_date = CURRENT_DATE
		`

	var id int
	err := db.QueryRow(
		sqlQuery,
		ph.ProductID,
	).Scan(&id)

	if err != nil {
		id = 0
	}

	return id == 0
}

func validateItem(item models.Item, compareList []models.Item) bool {
	title := strings.ToUpper(item.Title)
	title = strings.ReplaceAll(title, "™", "")

	for _, compare := range compareList {
		if strings.Contains(title, strings.ToUpper(compare.Title)) && item.Price >= compare.MinPrice && item.Price <= compare.MaxPrice {
			return true
		}
	}

	return false
}

func createPriceHistory(p models.Product, item models.Item, db *sql.DB) models.PriceHistory {
	// id := find id in db
	id, err := findProductByIdProduct(db, p.IDProduct)
	utils.CheckErr(err)

	if id == -1 {
		id, err = insertProduct(db, p)
		utils.CheckErr(err)
	}

	priceHistory := models.PriceHistory{
		ProductID: id,
		Price:     item.Price,
		Currency:  item.Currency,
	}
	return priceHistory
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default %d", key, defaultValue)
		return defaultValue
	}
	return value
}

// func CheckErr(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func LogErr(err error) {
// 	if err != nil {
// 		fmt.Printf("[ERROR] %s\n", err)
// 	}
// }
