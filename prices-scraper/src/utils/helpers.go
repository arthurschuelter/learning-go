package utils

import (
	"fmt"
	"items-scraper/src/models"
	"log"
	"net/url"
	"sort"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func LogErr(err error) {
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}
}

func SortList(priceList []models.Item) []models.Item {
	sort.Slice(priceList, func(i, j int) bool {
		return priceList[i].Price < priceList[j].Price
	})
	return priceList
}

func UrlDecode(s string) string {
	decoded, err := url.QueryUnescape(s)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}
