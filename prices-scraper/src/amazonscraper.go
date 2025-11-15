package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var amazonDomains = []string{
	"www.amazon.com.br",
}

func ScrapeAmazon(itemList []Item, db *sql.DB) {
	baseURL := "https://www.amazon.com.br/s?k="
	retailer := "Amazon"

	links := []string{}
	for _, item := range itemList {
		url := strings.ReplaceAll(strings.ToLower(item.Title), " ", "+")
		links = append(links, baseURL+url)
	}

	fmt.Printf("Scanning %s for:\n", retailer)
	for _, item := range itemList {
		fmt.Printf("- %s\n", item.Title)
	}

	c := colly.NewCollector(
		colly.AllowedDomains(amazonDomains...),
	)

	priceList := []Item{}
	productList := []Product{}
	priceHistoryList := []PriceHistory{}

	c.OnHTML("div[role='listitem']", func(e *colly.HTMLElement) {
		title := getTitleAmazon(e)
		price := getPriceAmazon(e)
		l := getLinkAmazon(e)
		id := getProductIDAmazon(l)

		item := Item{
			Title:    title,
			Currency: "R$",
			Price:    price,
			ID:       id,
			Link:     l,
		}

		product := Product{
			IDProduct:   id,
			ProductName: title,
			URL:         l,
			Retailer:    retailer,
		}

		if validateItem(item, itemList) {
			priceList = append(priceList, item)
			productList = append(productList, product)
			priceHistory := createPriceHistory(product, item, db)
			if priceHistory.ProductID != -1 {
				insertPriceHistory(db, priceHistory)
			}
			priceHistoryList = append(priceHistoryList, priceHistory)
		}
	})

	for i, link := range links {
		fmt.Printf("Scanning %s\n%s\n", itemList[i].Title, link)
		c.Visit(link)
		priceList = sortList(priceList)
		priceList = []Item{}
	}
}

func getTitleAmazon(e *colly.HTMLElement) string {
	return e.ChildText("div[data-cy='title-recipe'] h2 span")
}

func getPriceAmazon(e *colly.HTMLElement) float64 {
	var price float64

	var priceWhole = e.ChildText("div[data-cy='price-recipe'] span.a-price-whole")
	var priceFraction = e.ChildText("div[data-cy='price-recipe'] span.a-price-fraction")
	priceStr := priceWhole + priceFraction
	priceStr = strings.ReplaceAll(priceStr, ".", "")
	priceStr = strings.ReplaceAll(priceStr, ",", ".")

	var err error
	price, err = strconv.ParseFloat(priceStr, 64)
	if err != nil {
		price = 0
	}

	return price
}

func getLinkAmazon(e *colly.HTMLElement) string {
	baseUrl := "https://www.amazon.com.br"

	link := e.ChildAttr("span[data-component-type='s-product-image'] a.a-link-normal", "href")
	link = urlDecode(link)

	reURL := regexp.MustCompile(`^([^?]+?)(?:/ref=|$|\?)`)
	if match := reURL.FindStringSubmatch(link); len(match) > 1 {
		cleanURL := match[1]
		link = cleanURL
	}

	return baseUrl + link
}

func getProductIDAmazon(link string) string {
	reID := regexp.MustCompile(`/dp/([A-Z0-9]+)`)
	productID := ""
	if match := reID.FindStringSubmatch(link); len(match) > 1 {
		productID = match[1]
	}
	return productID
}
