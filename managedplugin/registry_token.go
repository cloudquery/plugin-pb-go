package managedplugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/distribution/reference"
)

const (
	dockerRegistryUsername = "managedplugin"
	pluginVersionHeader    = "X-Meta-Plugin-Version"
	teamNameHeader         = "X-Meta-User-Team-Name"
)

var authChallengeParamRegexp = regexp.MustCompile(`([a-zA-Z0-9_-]+)="([^"]*)"`)

type registryAuthChallenge struct {
	realm   string
	service string
}

func fetchRegistryBearerToken(ctx context.Context, httpClient *http.Client, registryBaseURL string, ref reference.NamedTagged, authToken, teamName string) (string, error) {
	challenge, err := fetchRegistryAuthChallenge(ctx, httpClient, registryBaseURL)
	if err != nil {
		return "", err
	}

	tokenURL, err := url.Parse(challenge.realm)
	if err != nil {
		return "", fmt.Errorf("failed to parse registry auth realm %q: %w", challenge.realm, err)
	}
	query := tokenURL.Query()
	query.Set("scope", fmt.Sprintf("repository:%s:pull", reference.Path(ref)))
	if challenge.service != "" {
		query.Set("service", challenge.service)
	}
	tokenURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create registry auth request: %w", err)
	}
	req.SetBasicAuth(dockerRegistryUsername, authToken)
	req.Header.Set(teamNameHeader, teamName)
	req.Header.Set(pluginVersionHeader, ref.Tag())

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request registry auth token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("registry auth request returned %s", resp.Status)
	}

	var payload struct {
		Token       string `json:"token"`
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("failed to decode registry auth response: %w", err)
	}
	if payload.Token != "" {
		return payload.Token, nil
	}
	if payload.AccessToken != "" {
		return payload.AccessToken, nil
	}
	return "", errors.New("registry auth response contained no token")
}

func fetchRegistryAuthChallenge(ctx context.Context, httpClient *http.Client, registryBaseURL string) (registryAuthChallenge, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimSuffix(registryBaseURL, "/")+"/v2/", nil)
	if err != nil {
		return registryAuthChallenge{}, fmt.Errorf("failed to create registry challenge request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return registryAuthChallenge{}, fmt.Errorf("failed to request registry challenge: %w", err)
	}
	defer resp.Body.Close()

	header := resp.Header.Get("WWW-Authenticate")
	if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return registryAuthChallenge{}, fmt.Errorf("registry %s did not return a bearer auth challenge", registryBaseURL)
	}

	challenge := registryAuthChallenge{}
	for _, match := range authChallengeParamRegexp.FindAllStringSubmatch(header[len("bearer "):], -1) {
		switch strings.ToLower(match[1]) {
		case "realm":
			challenge.realm = match[2]
		case "service":
			challenge.service = match[2]
		}
	}
	if challenge.realm == "" {
		return registryAuthChallenge{}, fmt.Errorf("registry %s auth challenge contained no realm", registryBaseURL)
	}
	return challenge, nil
}
