package nuntius

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func EmptyBody(json, method string, params, id interface{}) string {
	body := fmt.Sprintf(`{
    "jsonrpc": "%s",
    "method": "%s",
    "params": [%s],
    "id": %s
    }`, json, method, params, id)
	return body
}

func (p *Payload) NewPayload(json, method string, params, id interface{}) io.Reader {
	//var builder strings.Builder
	p.Jsonrpc, p.Method, p.Params, p.Id = json, method, params, id

	body := EmptyBody(json, method, params, id)
	aa := strings.NewReader(body)

	return aa
}

type Payload struct {
	Jsonrpc string
	Id      interface{}
	Method  string
	Params  interface{}
}

// TODO: Remove TempRequest !!!

func TempRequest() {

	p := &Payload{}
	payloadie := p.NewPayload("2.0", "eth_blockNumber", "", "1")

	fmt.Println(payloadie)
	url := "http://127.0.0.1:8080/v1/ethereum/mainnet"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadie)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer c8_oaXv7CUybRGueZPs1HegXIrNJIMUT.FDm9FQgWvn5K6CLv")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
