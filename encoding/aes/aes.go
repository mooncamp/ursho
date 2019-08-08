package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"strconv"

	"github.com/douglasmakey/ursho/encoding"
)

type coder struct {
	nonce []byte
	AEAD  cipher.AEAD
}

func New(secret []byte, nonce []byte) (encoding.Coder, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &coder{
		nonce: nonce,
		AEAD:  aesgcm,
	}, nil
}

func (c *coder) Encode(in int64) string {
	cipher := c.AEAD.Seal(nil, c.nonce, []byte(strconv.FormatInt(in, 16)), nil)
	return hex.EncodeToString(cipher)
}

func (c *coder) Decode(in string) (int64, error) {
	ciphertext, err := hex.DecodeString(in)
	if err != nil {
		return 0, err
	}

	plaintext, err := c.AEAD.Open(nil, c.nonce, ciphertext, nil)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(plaintext), 16, 64)
}
