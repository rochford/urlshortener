# urlshortener
Package urlshortener maps a long URL to a short URL.

First Go package.
GenerateShortURL function sends a request for a shortened URL on the channel and
expects the response within a time period. XXX if there is a timeout on the
channel response then an error is returned.

ResolveShortURL looks up the shortURL in the map. XXX if there is a timeout on
the channel response then an error is returned.
