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
	log.Println("Remote: get available jobs at " + url)

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

	log.Println("Jobs in queue:", len(jobs))
	return jobs
}

func (c *Client) AcceptJob(guid string) {
	url := fmt.Sprintf("%s/jobs/%s/accept", c.BaseURL, guid)

	req, _ := http.NewRequest(http.MethodPatch, url, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
  resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: AcceptJob -", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: AcceptJob reading response -", err)
		return
	}

	log.Println("AcceptJob response: " + string(body))
}

// TODO: upload Result into the job
func (c *Client) FinishedJob(guid string) {
	url := fmt.Sprintf("%s/jobs/%s/finished", c.BaseURL, guid)

	req, _ := http.NewRequest(http.MethodPatch, url, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
  resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: FinishedJob -", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: FinishedJob reading response -", err)
		return
	}

	log.Println("FinishedJob response: " + string(body))
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
