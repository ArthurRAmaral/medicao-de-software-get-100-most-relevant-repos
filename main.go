package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type QueryResponse struct {
	Data struct {
		Search struct {
			Nodes []struct {
				Name           string `json:"name"`
				StargazerCount int    `json:"stargazerCount"`
				Url            string `json:"url"`
			} `json:"nodes"`
			RepositoryCount int `json:"repositoryCount"`
		} `json:"search"`
	} `json:"data"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	url := "https://api.github.com/graphql"
	bearer := "Bearer " + os.Getenv("getHubToken")

	body := strings.NewReader("{\"query\":\"query MyQuery {\\n  search(type: REPOSITORY, query: \\\"stars:>100 sort:stars-desc\\\", first: 100) {\\n    nodes {\\n      ... on Repository {\\n        name\\n        stargazerCount\\n        url\\n      }\\n    }\\n    repositoryCount\\n  }\\n}\\n\",\"variables\":{},\"operationName\":\"MyQuery\"}")

	req, err := http.NewRequest(http.MethodPost, url, body)

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var dat QueryResponse

	if err := json.Unmarshal(bodyBytes, &dat); err != nil {
		log.Fatal(err)
	}

	repositoryNodes := dat.Data.Search.Nodes

	for item := range repositoryNodes {
		fmt.Println(repositoryNodes[item])
	}
}
