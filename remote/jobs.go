package remote

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/vialtek/whisker/model"
)

func (c *Client) AvailableJobs() []*model.Job {
	jobs := []*model.Job{}
	log.Println("Remote: getting list of jobs")

	resp, err := http.Get(c.BaseURL + "/jobs")
	if err != nil {
		log.Println("Error: AvailableJobs -", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &jobs); err != nil {
		log.Println("Error: AvailableJobs unmarshal -", err)
		return nil
	}

	return jobs
}
