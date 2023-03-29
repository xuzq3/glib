package rand

import (
	randlib "crypto/rand"
	"encoding/binary"
	"math/rand"
	"time"
)

const (
	CharsetAlphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	CharsetNumeral      = "1234567890"
	CharsetHex          = "1234567890abcdef"
	CharsetAlphanumeric = CharsetAlphabet + CharsetNumeral
)

func init() {
	seed, err := NewCryptoSeed()
	if err != nil {
		seed = NewTimeSeededSource()
	}
	rand.Seed(seed)
}

func NewCryptoSeed() (int64, error) {
	var seed int64
	err := binary.Read(randlib.Reader, binary.BigEndian, &seed)
	if err != nil {
		return 0, err
	}
	return seed, nil
}

func NewTimeSeededSource() int64 {
	return time.Now().UnixNano()
}

func Intn(n int) int {
	return rand.Intn(n)
}

func IntRange(min, max int) int {
	if min > max {
		return 0
	}
	if min == max {
		return min
	}
	r := rand.Intn(max - min)
	return min + r
}

func Charset(n int, charset string) string {
	randStr := make([]byte, n)
	charLen := len(charset)
	for i := 0; i < n; i++ {
		j := rand.Intn(charLen)
		randStr[i] = charset[j]
	}
	return string(randStr)
}

func String(n int) string {
	return Charset(n, CharsetAlphanumeric)
}

func Alphabet(n int) string {
	return Charset(n, CharsetAlphabet)
}

func NumeralStr(n int) string {
	return Charset(n, CharsetNumeral)
}

func HexStr(n int) string {
	return Charset(n, CharsetHex)
}

func ChoiceString(choices []string) string {
	var res string
	length := len(choices)
	i := rand.Intn(length)
	res = choices[i]
	return res
}

func ChoiceInt(choices []int) int {
	var res int
	length := len(choices)
	i := rand.Intn(length)
	res = choices[i]
	return res
}
