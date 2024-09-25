package network

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func PostReq(url string,data []byte) map[string]interface{} {
    response, err := http.Post(url, "application/json", bytes.NewBuffer(data))
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		log.Print(err)
	}

	return res
}