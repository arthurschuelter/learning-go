package main

type Item struct {
	ID       string
	Currency string
	Price    float64
	MinPrice float64
	Title    string
}

func main() {
	items := []Item{
		{ID: "1", Title: "Switch 2", MinPrice: 2200},
		{ID: "2", Title: "Steam Deck", MinPrice: 2200},
	}

	ScrapeMercadoLivre(items)
}
