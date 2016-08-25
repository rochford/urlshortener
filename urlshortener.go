// Package urlshortener maps a long URL to a short URL.
package urlshortener

import (
	"errors"
	"math/rand"
	"time"

	"golang.org/x/net/context"
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
// shortened URL on the channel. If the context is Done, then
// this function returns and error.
func GenerateShortURL(ctx context.Context) (string, error) {
	v := ctx.Value("originalUrl")
	if v == nil {
		return "", errors.New("longURL value not found in context")
	}
	longURL := ctx.Value("originalUrl").(string)
	if longURL == "" {
		return "", errors.New("longURL missing")
	}
	req := urlShortenerRequestResponse{longURL, make(chan (string), 1)}
	urlShortenerRequestChannel <- req

	select {
	// success
	case shortURL := <-req.responseChannel:
		return shortURL, nil
	// conext is done
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// ResolveShortURL looks up the shortURL in the map. If shortURL does not exist
// an empty string is returned along with an error. If the context is Done, then
// this function returns and error.
func ResolveShortURL(ctx context.Context) (string, error) {
	v := ctx.Value("id")
	if v == nil {
		return "", errors.New("shortURL value not found in context")
	}
	shortURL := ctx.Value("id").(string)
	if shortURL == "" {
		return "", errors.New("shortURL missing")
	}
	lookup := urlShortenerLookupRequestResponse{shortURL, make(chan (string), 1)}
	urlShortenerLookupChannel <- lookup

	select {
	// success
	case longURL := <-lookup.responseChannel:
		if longURL == "" {
			return "", errors.New("could not find a URL in map")
		}
		return longURL, nil
	// conext is done
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
