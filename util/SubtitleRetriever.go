package util

import (
	"log"
	"os"
)

func GetSubtitle(movieName, year string) (subtitileMessage *SearchResponse) {
	client := New(os.Getenv("API_KEY"))
	res, err := client.Search(movieName, year)
	if err != nil { /* ... */ }
	log.Printf(">>> results: %+v", res)
	return res
}
