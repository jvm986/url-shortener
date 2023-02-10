package shortener

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type MD5Shortener struct {
	length int
}

type MD5ShortenerConfig struct {
	PathLength int
}

func NewMD5Shortener(cfg MD5ShortenerConfig) Shortener {
	return &MD5Shortener{
		length: cfg.PathLength,
	}
}

func (s *MD5Shortener) Shorten(input string) (string, string, error) {
	u, err := s.SantizeURL(input)
	if err != nil {
		return "", "", errors.Wrapf(err, "unable to sanitize url")
	}

	hash := md5.Sum([]byte(u))
	str := hex.EncodeToString(hash[:])
	if s.length > len(str) {
		// We could log a warning here
		return str, u, nil
	}
	return str[0:s.length], u, nil
}

// SanitizeURL validates the input url, removes trailing slashes from path and add https where protocol is missing
func (s *MD5Shortener) SantizeURL(input string) (string, error) {
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "https://" + input
	}
	u, err := url.ParseRequestURI(input)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse input url: %s", input)
	}

	if u.Path == "/" {
		u.Path = ""
	}

	return u.String(), nil
}
