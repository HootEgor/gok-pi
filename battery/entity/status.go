package entity

import (
	"encoding/json"
	"fmt"
)

type Status struct {
	Level int `json:"level"`
}

func Parse(body []byte) (*Status, error) {
	var status Status
	err := json.Unmarshal(body, &status)
	if err != nil {
		return nil, fmt.Errorf("unmarshal status body: %s", err)
	}
	return &status, nil
}
