package util

import "log"

const   (
	apiKey = "f8720fd"
)
func GetSubtitle(movieName, year string) (subtitileMessage *SearchResponse) {
	client := New(apiKey)
	res, err := client.Search(movieName, year)
	if err != nil { /* ... */ }
	log.Printf(">>> results: %+v", res)
	return res
}
