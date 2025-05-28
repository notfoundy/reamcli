package ani

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Request struct {
	Query     string `json:"query"`
	Variables any    `json:"variables"`
}

type Response struct {
	Data   json.RawMessage  `json:"data"`
	Errors []map[string]any `json:"errors,omitempty"`
}

func doRequest(url, referer, userAgent string, gqlReq Request) (*Response, error) {
	payload, err := json.Marshal(gqlReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gqlResp Response
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, err
	}

	if len(gqlResp.Errors) > 0 {
		return &gqlResp, err
	}

	return &gqlResp, nil
}
