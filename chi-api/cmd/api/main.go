package main

import (
	"chi-api/internal/handler"
	"chi-api/internal/repository"
	"chi-api/internal/service"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	PrintBanner()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// --- Links Module Setup ---
	linksRepo := repository.NewLinkRepo()
	linksService := service.NewLinkService(linksRepo)
	linksHandler := handler.NewLinksHandler(linksService)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World from Chi v0"))
	})

	// --- Links Routes ---
	r.Route("/links", func(r chi.Router) {
		r.Get("/", linksHandler.GetLinks)
	})

	http.ListenAndServe(":3000", r)
}

func PrintBanner() {
	banner := `
(Chi @ v0)
 ________  ________          ________  ________  ___     
|\   ____\|\   __  \        |\   __  \|\   __  \|\  \    
\ \  \___|\ \  \|\  \       \ \  \|\  \ \  \|\  \ \  \   
 \ \  \  __\ \  \\\  \       \ \   __  \ \   ____\ \  \  
  \ \  \|\  \ \  \\\  \       \ \  \ \  \ \  \___|\ \  \ 
   \ \_______\ \_______\       \ \__\ \__\ \__\    \ \__\
    \|_______|\|_______|        \|__|\|__|\|__|     \|__|
	`

	fmt.Println(banner)
}
