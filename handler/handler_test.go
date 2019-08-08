package handler

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/douglasmakey/ursho/aes"
	"github.com/douglasmakey/ursho/storage/dgraph"
)

func Test_store_and_retrive_link_with_dgraph(t *testing.T) {
	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	coder, err := aes.New(key, nonce)
	if err != nil {
		t.Fatalf("new coder: %v", err)
	}

	s, err := dgraph.New("localhost", "9080", coder)
	if err != nil {
		t.Fatalf("new dgraph storage service: %v", err)
	}

	testServer := httptest.NewServer(New("", s))

	body := map[string]interface{}{
		"url": "www.google.de",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("encode body: %v", err)
	}

	client := http.Client{}
	resp, err := client.Post(fmt.Sprintf("%s/encode/", testServer.URL), "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		t.Fatalf("client do: %v", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("response status not ok: %d %s", resp.StatusCode, string(respBody))
	}

	var responseData response
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	url := fmt.Sprintf("%s/%s", testServer.URL, responseData.Data.(string))

	getResponse, err := client.Get(url)
	if err != nil {
		t.Fatalf("get response: %v", err)
	}

	getRespBody, err := ioutil.ReadAll(getResponse.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("response status not ok: %d %s", getResponse.StatusCode, string(getRespBody))
	}
}
