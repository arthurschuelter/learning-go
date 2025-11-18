package info

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arthurschuelter/go-git/models"
	"github.com/arthurschuelter/go-git/utils"
)

func GetRepoData(url string, client *http.Client, token string) []models.Repo {
	req := utils.ConfigRequest(url, token)
	body, err := utils.MakeRequestAndRead(client, req)

	if err != nil {
		panic(err)
	}

	// fmt.Printf("%s\n", body)
	repos, err := ReadRepos(body)

	if err != nil {
		panic(err)
	}

	return repos
}

func ReadRepos(data []byte) ([]models.Repo, error) {
	var repos []models.Repo
	err := json.Unmarshal(data, &repos)
	return repos, err
}

func PrintRepoData(repos []models.Repo) {
	fmt.Printf("%d repositories\n", len(repos))
	stars := 0
	for _, r := range repos {
		stars += r.Stars
	}
	fmt.Printf("Stars: %d ⭐\n", stars)

}
