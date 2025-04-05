package utils

import (
	"fmt"
)

func ConvertToJSON(data []map[string]interface{}, headerKey string) (map[string][]map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	result := make(map[string][]map[string]interface{})

	for _, row := range data {
		// Check if headerKey exists
		if headerValue, ok := row[headerKey]; ok {
			headerStr, valid := headerValue.(string)
			if !valid {
				return nil, fmt.Errorf("header key '%s' is not a string", headerKey)
			}
			// Append the row to the appropriate group
			result[headerStr] = append(result[headerStr], row)
		} else {
			return nil, fmt.Errorf("header key '%s' not found in row", headerKey)
		}
	}

	return result, nil
}
