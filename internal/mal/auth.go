package mal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/notfoundy/reamcli/internal/utils"
)

func (c *Client) StartOAuth() error {
	authURL := fmt.Sprintf(
		"https://myanimelist.net/v1/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s&code_challenge_method=plain",
		c.ClientID,
		url.QueryEscape(c.RedirectURI),
		c.State,
		c.CodeChallenge,
	)

	fmt.Println("Opening browser for authorization...")
	openBrowser(authURL)

	codeCh := make(chan string)

	go func() {
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "No code in URL", http.StatusBadRequest)
				return
			}
			fmt.Fprintln(w, "Authorization successful. You can close this tab.")
			codeCh <- code
		})

		authPort := utils.GetCallbackPort()
		log.Printf("Waiting for callback on http://localhost:%d/callback ...\n", authPort)
		_ = http.ListenAndServe(fmt.Sprintf(":%d", authPort), nil)
	}()

	code := <-codeCh
	return c.exchangeCodeForToken(code)
}

func (c *Client) exchangeCodeForToken(code string) error {
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("client_id", c.ClientID)
	values.Set("redirect_uri", c.RedirectURI)
	values.Set("code_verifier", c.CodeVerifier)

	resp, err := c.HTTPClient.PostForm("https://myanimelist.net/v1/oauth2/token", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to exchange token, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var t Token
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.AccessToken = t.AccessToken
	c.RefreshToken = t.RefreshToken
	c.ExpiresAt = time.Now().Add(time.Duration(t.ExpiresAt) * time.Second)
	err = c.SaveToken()
	return err
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		cmd = "xdg-open"
	}
	if cmd != "" {
		args = append([]string{url}, args...)
		exec.Command(cmd, args...).Start()
	}
}
