package mal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) UpdateAnimeStatus(malId int, status string, episodesWatched int) error {
	url := fmt.Sprintf("https://api.myanimelist.net/v2/anime/%d/my_list_status", malId)
	payload := map[string]any{
		"status":               status,
		"num_episodes_watched": episodesWatched,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MAL API error: %s", resp.Status)
	}

	return nil
}
