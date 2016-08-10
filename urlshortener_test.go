package urlshortener

import (
	"sync"
	"testing"
	"unicode"
)

func TestGenerateShortURL(t *testing.T) {
	expectedLongURL := "http://www.abc.org"
	for i := 0; i < 1e3; i++ {
		str := GenerateShortURL(expectedLongURL)
		if len(str) != stringLength {
			t.Errorf("GenerateShortURL returned incorrect length %d, expected %d",
				len(str), stringLength)
		}
		for j := 0; j < len(str); j++ {
			if unicode.IsLetter(rune(str[j])) == false {
				t.Errorf("GenerateShortURL not all letters: %s", str)
				break
			}
		}
		acutalLongURL := urlMap[str]
		if acutalLongURL != expectedLongURL {
			t.Errorf("GenerateShortURL did not update the urlMap")
		}
	}
}

func TestConcurrentAccessUrlMap(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		// Don't put the Add() call in the goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			GenerateShortURL("http://www.something.org")
		}()
	}

	wg.Wait()
}

func TestConcurrentResolveShortURL(t *testing.T) {

	expectedMapKey := "ABCEDEF"
	expectedMapValue := "http://www.wikipedia.org"
	urlMap[expectedMapKey] = expectedMapValue
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		// Don't put the Add() call in the goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			actualLongURL := ResolveShortURL(expectedMapKey)
			if actualLongURL != expectedMapValue {
				t.Errorf("TestResolveShortURL incorrect longURL")
			}
		}()
	}

	wg.Wait()
}
