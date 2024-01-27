package periodicmd

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Name          string        `json:"name"`
	Frequency     TaskFrequency `json:"frequency"`
	StartDate     string        `json:"startDate"`
	OffsetDays    *int          `json:"offsetDays"`
	CreateCommand string        `json:"createCommand"`
	LinkCommand   string        `json:"linkCommand"`
}

type TaskFrequency struct {
	Years  int `json:"years"`
	Months int `json:"months"`
	Weeks  int `json:"weeks"`
	Days   int `json:"days"`
}

func ReadConfig(
	path string,
) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return &config, nil
}
