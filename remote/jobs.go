package remote

import (
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

	resp, err := http.Get(url)
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
