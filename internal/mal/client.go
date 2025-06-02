package mal

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/notfoundy/reamcli/internal/utils"
)

type Client struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string

	// PKCE
	State         string
	CodeVerifier  string
	CodeChallenge string

	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time

	HTTPClient *http.Client
	mu         sync.Mutex
}

// CodeVerifier and CodeChallenge are the same because method is plain
// MAL only support plain method
func NewClient(clientID, redirectURI string) *Client {
	verifier := utils.RandomString(64)
	c := &Client{
		ClientID:      clientID,
		RedirectURI:   redirectURI,
		State:         uuid.NewString(),
		CodeVerifier:  verifier,
		CodeChallenge: verifier,
		HTTPClient:    &http.Client{},
	}
	_ = c.LoadToken()

	return c
}
