package main

import (
	"os"

	"bookish-potato/chap05/thesaurus/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{apiKey}

}
