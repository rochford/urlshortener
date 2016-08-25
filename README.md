# urlshortener
Package urlshortener maps a long URL to a short URL.

First Go package.

GenerateShortURL function sends a request for a shortened URL. Returns an error
if the context is done.

ResolveShortURL looks up the shortURL in the map. Returns an error if the
context is done.
