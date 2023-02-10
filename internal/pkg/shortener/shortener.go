package shortener

//go:generate mockgen -destination=mocks/mock_shortener.go -package=shortenermocks . Shortener
type Shortener interface {
	// Shorten converts a *valid* input url to a unique key, returning the key and a sanitized input url
	Shorten(string) (string, string, error)
}
