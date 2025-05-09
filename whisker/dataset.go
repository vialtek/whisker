package whisker

import (
	"log"
	"os"
	"strings"
)

type Dataset struct {
	Name string
}

func loadDatasets() []*Dataset {
	var datasets []*Dataset

	entries, err := os.ReadDir(GetConfig().DatasetDirPath)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				datasets = append(datasets, &Dataset{Name: entry.Name()})
			}
		}
	}

	log.Println("Datasets loaded:", len(datasets))
	if len(datasets) == 0 {
		log.Println("Warning: no dataset loaded.")
	}

	return datasets
}
