package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

type Client struct {
	hc             http.Client
	Token          string
	remainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, hc: c}
}

type Search struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id:"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	ProtographerURL string      `json:"protographer_url`
	src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Potrait   string `json:"potrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

func main() {
	os.Setenv("PexelToken", "mQnfEE1IBvfEqGQXkayvAoN6cCjPs8V7fjajlOFyl0Pdw6pVmFT0y3n9")
	var TOKEN = os.Getenv("PexelToken")

	var c = NewClient(TOKEN)

	result, err := c.SearchPhotos("waves")
	if err != nil {
		fmt.Errorf("search error %v", err)
	}
	if result.Page == 0 {
		fmt.Errorf("search result wrong")
	}
	fmt.Print(result)
}
