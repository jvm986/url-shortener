package shortener_test

import (
	"testing"

	"github.com/jvm986/url-shortener/internal/pkg/shortener"
	"github.com/stretchr/testify/assert"
)

type ShortenCase struct {
	input         string
	length        int
	expectedKey   string
	expectedUrl   string
	expectedError string
}

func TestShorten(t *testing.T) {
	tests := map[string]ShortenCase{
		"happy path": {
			input:       "https://example.com",
			length:      10,
			expectedKey: "c984d06aaf",
			expectedUrl: "https://example.com",
		},
		"happy path, long length": {
			input:       "https://example.com",
			length:      999,
			expectedKey: "c984d06aafbecf6bc55569f964148ea3",
			expectedUrl: "https://example.com",
		},
		"happy path with trailing slash": {
			input:       "https://www.example.com/",
			expectedKey: "e149be135a", // same value as without slash
			expectedUrl: "https://www.example.com",
			length:      10,
		},
		"happy path with query params": {
			input:       "https://www.example.com?param=value&otherparam=othervalue",
			expectedKey: "67a5634dee",
			expectedUrl: "https://www.example.com?param=value&otherparam=othervalue",
			length:      10,
		},
		"happy path with query params and trailing slash": {
			input:       "https://www.example.com/?param=value&otherparam=othervalue",
			expectedKey: "67a5634dee", // same value as without slash
			expectedUrl: "https://www.example.com?param=value&otherparam=othervalue",
			length:      10,
		},
		"happy path without protocol": {
			input:       "www.example.com",
			expectedKey: "e149be135a",
			expectedUrl: "https://www.example.com",
			length:      10,
		},
		"invalid url": {
			input:         "i am invalid",
			expectedKey:   "",
			expectedUrl:   "",
			length:        10,
			expectedError: "unable to sanitize url",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := shortener.NewMD5Shortener(shortener.MD5ShortenerConfig{
				PathLength: tc.length,
			})
			actualKey, actualUrl, err := s.Shorten(tc.input)
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			}
			assert.Equal(t, tc.expectedKey, actualKey)
			assert.Equal(t, tc.expectedUrl, actualUrl)
		})
	}
}
