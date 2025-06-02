package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/vialtek/whisker/model"
)

func (c *Client) AvailableJobs() []*model.Job {
	url := c.BaseURL + "/jobs"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error: AvailableJobs -", err)
		return nil
	}
	defer resp.Body.Close()

	jobs := []*model.Job{}
	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &jobs); err != nil {
		log.Println("Error: AvailableJobs unmarshal -", err)
		return nil
	}

	return jobs
}

func (c *Client) ChangeJobState(guid string, state string) {
	url := fmt.Sprintf("%s/jobs/%s/%s", c.BaseURL, guid, state)

	log.Println("Remote: changing job state on " + url)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		log.Println("ERROR: ChangeJobState -", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: ChangeJobState reading response -", err)
		return
	}

	log.Println("ChangeJobState response:", string(body))
}

func (c *Client) SendJobOutput(guid string, output []string) {
	url := fmt.Sprintf("%s/jobs/%s/output_log", c.BaseURL, guid)

	payload := map[string][]string{
		"output_log": output,
	}

	log.Println("Remote: sending output log to " + url)
	json_data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error: SendJobOutput marshal -", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Println("Error: SendJobOutput -", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: SendJobOutput reading response -", err)
		return
	}

	log.Println("Heartbeat SendJobOutput: " + string(body))
}
