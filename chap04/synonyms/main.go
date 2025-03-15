package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/BhuwanSingh/bookish-potato/chap05/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalln("failes when looking for synonyms")
		}
		if len(syns) == 0 {
			log.Fatalln(" could not find synonyms")
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
