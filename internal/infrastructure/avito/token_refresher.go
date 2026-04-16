package avito

import (
	"context"
	"encoding/json"
	"fmt"
	"hr-bot-ddd-example/internal/domain/entity"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// TokenRefresher реализует port.TokenRefresher для Avito.
type TokenRefresher struct{}

func NewTokenRefresher() *TokenRefresher {
	return &TokenRefresher{}
}

func (r *TokenRefresher) Provider() entity.Provider {
	return entity.ProviderAvito
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (r *TokenRefresher) Refresh(ctx context.Context, cred entity.Credential, decryptionKey string) (string, int, error) {
	decryptedSecret := fmt.Sprintf("%s %s", cred.ClientSecretEncrypted, decryptionKey)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", cred.ClientID)
	data.Set("client_secret", decryptedSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.avito.ru/token/", strings.NewReader(data.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("avito token API returned %d: %s", resp.StatusCode, body)
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", 0, fmt.Errorf("parse response: %w", err)
	}

	return tokenResp.AccessToken, tokenResp.ExpiresIn, nil
}
