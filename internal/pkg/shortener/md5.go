package shortener

import (
	"crypto/md5"
	"encoding/hex"
)

type MD5Shortener struct{}

var _ = Shortener(&MD5Shortener{})

func (s *MD5Shortener) Shorten(input string) (string, string, error) {
	hash := md5.Sum([]byte(input))
	str := hex.EncodeToString(hash[:])
	return str, input, nil
}
