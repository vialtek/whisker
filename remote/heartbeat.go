package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) SendHeartbeat(status map[string]string) {
	url := fmt.Sprintf("%s/heartbeat", c.BaseURL)

	log.Println("Remote: sending heartbeat to " + url)

	json_data, err := json.Marshal(status)
	if err != nil {
		log.Println("Error: SendHeartbeat marshal -", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Println("Error: SendHeartbeat -", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: SendHeartbeat reading response -", err)
		return
	}

	log.Println("Heartbeat response: " + string(body))
}
