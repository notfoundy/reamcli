package mal

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	malApi            = "https://api.myanimelist.net/v2"
	StatusWatching    = "watching"
	StatusCompleted   = "completed"
	StatusOnHold      = "on_hold"
	StatusDropped     = "dropped"
	StatusPlanToWatch = "plan_to_watch"
)

var validStatuses = map[string]bool{
	StatusWatching:    true,
	StatusCompleted:   true,
	StatusOnHold:      true,
	StatusDropped:     true,
	StatusPlanToWatch: true,
}

func (c *Client) UpdateAnimeStatus(malId string, status string, episodesWatched int) error {
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %q", status)
	}

	reqUrl := fmt.Sprintf("%s/anime/%s/my_list_status", malApi, malId)
	payload := url.Values{}
	payload.Set("status", status)
	payload.Set("num_watched_episodes", strconv.Itoa(episodesWatched))

	req, err := http.NewRequest("PATCH", reqUrl, bytes.NewReader([]byte(payload.Encode())))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
