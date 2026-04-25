package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const skillHubBaseURL = "https://www.skillhub.club/api/v1"

type SkillHubResult struct {
	Items      []map[string]interface{} `json:"items"`
	Raw        interface{}              `json:"raw"`
	Configured bool                     `json:"configured"`
	Source     string                   `json:"source"`
}

func fetchSkillHubCatalog(limit int) (SkillHubResult, error) {
	apiKey := strings.TrimSpace(os.Getenv("SKILLHUB_API_KEY"))
	if apiKey == "" {
		return SkillHubResult{
			Items:      nil,
			Raw:        nil,
			Configured: false,
			Source:     "skillhub",
		}, nil
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/skills/catalog?limit=%d&sort=score", skillHubBaseURL, limit), nil)
	if err != nil {
		return SkillHubResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return doSkillHubRequest(req)
}

func fetchSkillHubSearch(query string, limit int) (SkillHubResult, error) {
	apiKey := strings.TrimSpace(os.Getenv("SKILLHUB_API_KEY"))
	if apiKey == "" {
		return SkillHubResult{
			Items:      nil,
			Raw:        nil,
			Configured: false,
			Source:     "skillhub",
		}, nil
	}

	body, err := json.Marshal(map[string]interface{}{
		"query":  query,
		"limit":  limit,
		"method": "hybrid",
	})
	if err != nil {
		return SkillHubResult{}, err
	}

	req, err := http.NewRequest(http.MethodPost, skillHubBaseURL+"/skills/search", bytes.NewReader(body))
	if err != nil {
		return SkillHubResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	return doSkillHubRequest(req)
}

func doSkillHubRequest(req *http.Request) (SkillHubResult, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return SkillHubResult{}, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return SkillHubResult{}, err
	}
	if resp.StatusCode >= 400 {
		return SkillHubResult{}, fmt.Errorf("SkillHub API request failed: %s", resp.Status)
	}

	var raw interface{}
	if err := json.Unmarshal(content, &raw); err != nil {
		return SkillHubResult{}, err
	}

	return SkillHubResult{
		Items:      extractSkillHubItems(raw),
		Raw:        raw,
		Configured: true,
		Source:     "skillhub",
	}, nil
}

func extractSkillHubItems(raw interface{}) []map[string]interface{} {
	switch data := raw.(type) {
	case []interface{}:
		return interfaceSliceToMaps(data)
	case map[string]interface{}:
		for _, key := range []string{"items", "skills", "results", "data"} {
			if candidate, ok := data[key]; ok {
				switch typed := candidate.(type) {
				case []interface{}:
					return interfaceSliceToMaps(typed)
				case map[string]interface{}:
					for _, nestedKey := range []string{"items", "skills", "results"} {
						if nested, ok := typed[nestedKey]; ok {
							if arr, ok := nested.([]interface{}); ok {
								return interfaceSliceToMaps(arr)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func interfaceSliceToMaps(items []interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if entry, ok := item.(map[string]interface{}); ok {
			results = append(results, entry)
		}
	}
	return results
}
