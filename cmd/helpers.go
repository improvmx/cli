package cmd

import (
	"encoding/json"
	"fmt"
)

func parseResponse(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	return nil
}
