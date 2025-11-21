package repository

import (
	"chi-api/internal/models"
	"database/sql"
)

type LinksRepo interface {
	FindAll() ([]models.Link, error)
}

type SQLLinkRepo struct {
	DB *sql.DB
}

type MockLinkRepo struct {
}

// [TODO] Add DB here
func NewLinkRepo() LinksRepo {
	return MockLinkRepo{}
}

func (m MockLinkRepo) FindAll() ([]models.Link, error) {
	return []models.Link{
		{
			ShortCode: "abcdef",
			URL:       "google.com",
		},
		{
			ShortCode: "fedcba",
			URL:       "youtube.com",
		},
	}, nil
}
