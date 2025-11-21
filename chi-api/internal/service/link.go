package service

import (
	"chi-api/internal/models"
	"chi-api/internal/repository"
	"fmt"
)

type LinkServicer interface {
	ListAllLinks() ([]models.Link, error)
	// Save(link *models.Link) error
	// FindByShortCode(code string) (*models.Link, error)
}

type LinkService struct {
	Repo repository.LinksRepo
}

func NewLinkService(repo repository.LinksRepo) LinkServicer {
	return &LinkService{
		Repo: repo,
	}
}

func (ls *LinkService) ListAllLinks() ([]models.Link, error) {
	links, err := ls.Repo.FindAll()

	if err != nil {
		return nil, fmt.Errorf("failed to get all links: %w", err)
	}

	return links, nil
}
