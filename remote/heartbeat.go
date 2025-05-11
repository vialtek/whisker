package remote

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// TODO: load dynamically
func webShellEndpoint() string {
	return "http://localhost:4567"
}

func SendHeartbeat(status map[string]string) {
	log.Println("Remote: sending heartbeat to " + webShellEndpoint())

	json_data, err := json.Marshal(status)
	if err != nil {
		log.Println("Error: SendHeartbeat marshal -", err)
		return
	}

	resp, err := http.Post(webShellEndpoint()+"/heartbeat", "application/json", bytes.NewBuffer(json_data))
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
