package mal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/notfoundy/reamcli/internal/utils"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
}

const tokenFileName = "mal_token.json"

func (c *Client) getTokenPath() string {
	cfgDir := utils.GetConfigDir()
	return filepath.Join(cfgDir, tokenFileName)
}

func (c *Client) SaveToken() error {
	path := c.getTokenPath()
	data, err := json.MarshalIndent(Token{
		AccessToken:  c.AccessToken,
		RefreshToken: c.RefreshToken,
		ExpiresAt:    c.ExpiresAt.Unix(),
	}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (c *Client) LoadToken() error {
	path := c.getTokenPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var t Token
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	c.AccessToken = t.AccessToken
	c.RefreshToken = t.RefreshToken
	c.ExpiresAt = time.Unix(t.ExpiresAt, 0)
	return nil
}
