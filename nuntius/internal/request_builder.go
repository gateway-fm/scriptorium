package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gateway-fm/scriptorium/nuntius/config"
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

// TODO: Remove TempRequest !!!

func TempRequest() error {
	conf := &config.Config{}
	payload, err := conf.ParsePayload()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	emptybody := EmptyBody(payload.Jsonrpc, payload.Method, payload.Params, payload.Id)
	fmt.Println(payload.Jsonrpc)

	payloadie := strings.NewReader(emptybody)

	urlAdd := payload.Url

	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, urlAdd, payloadie)

	if err != nil {
		return fmt.Errorf("%w", err)

	}
	req.Header.Add("Authorization", "Bearer "+payload.BearerKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Println(string(body))
	return nil
}
