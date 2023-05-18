package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	src := []byte("hello world")

	encoded := Base64Encode(src)
	decoded, err := Base64Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(decoded))

	assert.Equal(t, src, decoded)
}
