package components

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"golang.org/x/oauth2/google"
)

const clientCredentialsPath = "client_credentials.json"

// Client creates a new OAuth2 HTTP client that has the given scope of access.
func Client(ctx context.Context, scopeURL string) (*http.Client, error) {
	config, err := ClientConfig(scopeURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file: %v", err)
	}

	token, err := Token(config)
	if err != nil {
		return nil, fmt.Errorf("Unable to authorize requests: %v", err)
	}

	return config.Client(ctx, token), nil
}

// ClientConfig returns an OAuth2 client config for the given scope.
func ClientConfig(scopeURL string) (*oauth2.Config, error) {
	credentials, err := getClientCredentials()
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	return google.ConfigFromJSON(credentials, scopeURL)
}

// getClientCredentials gets OAuth2 client credentials.
func getClientCredentials() ([]byte, error) {
	return ioutil.ReadFile(clientCredentialsPath)
}
