// Package urlshortener maps a long URL to a short URL.
package urlshortener

import (
	"math/rand"
	"time"
)

var stringLength = 6
var urlMap map[string]string

type urlShortenerRequestResponse struct {
	// filled in request
	longURL string
	// filled in request, message on channel indicates short URL
	responseChannel chan (string)
}

type urlShortenerLookupRequestResponse struct {
	// filled in request
	shortURL string
	// filled in request, message on channel indicates long URL
	responseChannel chan (string)
}

var urlShortenerRequestChannel chan (urlShortenerRequestResponse)
var urlShortenerLookupChannel chan (urlShortenerLookupRequestResponse)

func init() {
	rand.Seed(time.Now().Unix())
	urlMap = make(map[string]string)

	urlShortenerRequestChannel = make(chan (urlShortenerRequestResponse))
	urlShortenerLookupChannel = make(chan (urlShortenerLookupRequestResponse))

	go func() {
		for {
			select {
			case create := <-urlShortenerRequestChannel:
				var shortURL string
				for i := 0; i < stringLength; i++ {
					r := rand.Intn(26)
					letter := rune(r + 65)
					shortURL += string(letter)
				}
				urlMap[shortURL] = create.longURL
				// time.Sleep(time.Second * 1)
				// fmt.Println("...sending slow response: ", shortURL)
				create.responseChannel <- shortURL
				close(create.responseChannel)
			case lookup := <-urlShortenerLookupChannel:
				lookup.responseChannel <- urlMap[lookup.shortURL]
				close(lookup.responseChannel)
			}
		}
	}()
}

// GenerateShortURL function returns a shortened URL. It sends a request for a
// shortened URL on the channel
// and expects the response within a time period. XXX if there is a timeout
// on the channel response then an error is returned.
func GenerateShortURL(longURL string) string {
	req := urlShortenerRequestResponse{longURL, make(chan (string), 1)}

	urlShortenerRequestChannel <- req
	select {

	case shortURL := <-req.responseChannel:
		// success
		return shortURL
	case <-time.After(time.Millisecond * 100):
		// timeout on channel
		// TODO: XXX timeout error returned
		return ""
	}
}

// ResolveShortURL looks up the shortURL in the map
func ResolveShortURL(shortURL string) string {
	lookup := urlShortenerLookupRequestResponse{shortURL, make(chan (string), 1)}

	urlShortenerLookupChannel <- lookup
	longURL := <-lookup.responseChannel

	// TODO: XXX timeout error returned
	return longURL
}
