package aes

import (
	"encoding/hex"
	"testing"
)

func Test_it_encodes_and_decodes(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	coder, err := New(key)
	if err != nil {
		t.Fatalf("new coder: %v", err)
	}

	expected := int64(392048493)
	c, err := coder.Encode(expected)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	plain, err := coder.Decode(c)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if plain != expected {
		t.Errorf("%d != %d", expected, plain)
	}
}
