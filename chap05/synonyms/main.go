package main

import (
	"os"

	"github.com/BhuwanSingh/bookish-potato/chap05/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{apiKey}

}
