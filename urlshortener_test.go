package urlshortener

import (
	"sync"
	"testing"
	"unicode"

	"golang.org/x/net/context"
)

func TestGenerateShortURL(t *testing.T) {
	expectedLongURL := "http://www.abc.org"
	for i := 0; i < 1e3; i++ {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "originalURL", expectedLongURL)
		str, _ := GenerateShortURL(ctx)
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
			ctx := context.Background()
			ctx = context.WithValue(ctx, "originalURL", "http://www.something.org")
			GenerateShortURL(ctx)
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
			ctx := context.Background()
			ctx = context.WithValue(ctx, "id", expectedMapKey)

			actualLongURL, err := ResolveShortURL(ctx)
			if actualLongURL != expectedMapValue {
				t.Errorf("TestResolveShortURL incorrect longURL")
			}
			if err != nil {
				t.Errorf("TestResolveShortURL error value must be nil")
			}
		}()
	}

	wg.Wait()
}

func TestResolveShortURLUnknown(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "id", "UNKNOWN")

	actualLongURL, err := ResolveShortURL(ctx)
	if actualLongURL != "" {
		t.Errorf("TestResolveShortURL incorrect longURL")
	}
	if err == nil {
		t.Errorf("TestResolveShortURL error value must not be nil")
	}
}

func TestGenerateShortURLMissingLongURL(t *testing.T) {
	ctx := context.Background()
	shortURL, err := GenerateShortURL(ctx)

	if err == nil {
		t.Errorf("TestGenerateShortURLContextDone error value must not be nil")
	}
	if shortURL != "" {
		t.Errorf("TestGenerateShortURLContextDone incorrect shortURL")
	}
}

func TestResolveShortURLMissing(t *testing.T) {
	ctx := context.Background()

	actualLongURL, err := ResolveShortURL(ctx)
	if err == nil {
		t.Errorf("TestResolveShortURLContextDone error value must not be nil")
	}
	if actualLongURL != "" {
		t.Errorf("TestResolveShortURLContextDone incorrect longURL")
	}
}

func TestGenerateShortURLEmptyLongURL(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "originalURL", "")
	shortURL, err := GenerateShortURL(ctx)

	if err == nil {
		t.Errorf("TestGenerateShortURLContextDone error value must not be nil")
	}
	if shortURL != "" {
		t.Errorf("TestGenerateShortURLContextDone incorrect shortURL")
	}
}

func TestResolveShortURLEmpty(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "id", "")

	actualLongURL, err := ResolveShortURL(ctx)
	if err == nil {
		t.Errorf("TestResolveShortURLContextDone error value must not be nil")
	}
	if actualLongURL != "" {
		t.Errorf("TestResolveShortURLContextDone incorrect longURL")
	}
}

func TestGenerateShortURLContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "originalUrl", "http://www.something.org")
	cancel()
	<-ctx.Done()

	shortURL, err := GenerateShortURL(ctx)

	if shortURL != "" {
		t.Errorf("TestGenerateShortURLContextDone incorrect shortURL")
	}
	if err == nil {
		t.Errorf("TestGenerateShortURLContextDone error value must not be nil")
	}
}

func TestResolveShortURLContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "id", "UNKNOWN")
	cancel()
	<-ctx.Done()

	actualLongURL, err := ResolveShortURL(ctx)
	if actualLongURL != "" {
		t.Errorf("TestResolveShortURLContextDone incorrect longURL")
	}
	if err == nil {
		t.Errorf("TestResolveShortURLContextDone error value must not be nil")
	}
}
