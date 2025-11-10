package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var domains = []string{
	"www.mercadolivre.com.br",
	"lista.mercadolivre.com.br",
}

func ScrapeMercadoLivre(itemList []Item) {
	baseURL := "https://lista.mercadolivre.com.br/"
	links := []string{}
	for _, item := range itemList {
		url := strings.ReplaceAll(strings.ToLower(item.Title), " ", "-")
		links = append(links, baseURL+url)
	}

	fmt.Printf("Scanning Mercado Livre for:\n")
	for _, item := range itemList {
		fmt.Printf("- %s\n", item.Title)
	}

	c := colly.NewCollector(
		colly.AllowedDomains(domains...),
	)

	priceList := []Item{}

	c.OnHTML("div.ui-search-result__wrapper", func(e *colly.HTMLElement) {
		title := e.ChildText("h3.poly-component__title-wrapper")
		var price float64
		e.ForEach("div.poly-price__current span.andes-money-amount__fraction", func(i int, el *colly.HTMLElement) {
			if i == 0 {
				priceStr := strings.ReplaceAll(el.Text, ".", "")
				var err error
				price, err = strconv.ParseFloat(priceStr, 64)
				if err != nil {
					price = 0
				}

			}
		})

		item := Item{
			Title:    title,
			Currency: "R$",
			Price:    price,
			ID:       "",
		}
		if validateItem(item, itemList) {
			// printItem(item)
			priceList = append(priceList, item)
			// fmt.Printf("Link: %s\n", link)
		}
	})

	for i, link := range links {
		fmt.Printf("Scanning %s\n", itemList[i].Title)
		c.Visit(link)
		priceList = sortList(priceList)
		for _, item := range priceList {
			printItem(item)
		}
		priceList = []Item{}
	}
}

func printItem(item Item) {
	fmt.Printf("  * %s -- %s %.2f\n", item.Title, item.Currency, item.Price)
}

func validateItem(item Item, compareList []Item) bool {
	title := strings.ToUpper(item.Title)

	for _, compare := range compareList {
		if strings.Contains(title, strings.ToUpper(compare.Title)) && item.Price >= compare.MinPrice {
			return true
		}
	}

	return false
}

func sortList(priceList []Item) []Item {
	sort.Slice(priceList, func(i, j int) bool {
		return priceList[i].Price < priceList[j].Price
	})
	return priceList
}
