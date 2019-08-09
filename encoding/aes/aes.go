package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	"github.com/douglasmakey/ursho/encoding"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func nonce() ([]byte, error) {
	buf := make([]byte, 12)
	_, err := rand.Read(buf)
	return buf, err
}

type coder struct {
	AEAD cipher.AEAD
}

func New(secret []byte) (encoding.Coder, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &coder{
		AEAD: aesgcm,
	}, nil
}

func (c *coder) Encode(in int64) (string, error) {
	nonce, err := nonce()
	if err != nil {
		return "", err
	}
	cipher := c.AEAD.Seal(nil, nonce, []byte(strconv.FormatInt(in, 16)), nil)
	return hex.EncodeToString(append(cipher, nonce...)), nil
}

func (c *coder) Decode(in string) (int64, error) {
	code, err := hex.DecodeString(in)
	if err != nil {
		return 0, err
	}

	ciphertext := code[:len(code)-12]
	nonce := code[len(ciphertext):]

	plaintext, err := c.AEAD.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(plaintext), 16, 64)
}
