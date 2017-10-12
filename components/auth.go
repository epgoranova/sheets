package components

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/browser"

	"golang.org/x/oauth2"
)

// Token uses a OAuth2 Config to retrieve a Token for the current user. The token
// is cached locally. If the caching fails, an error is not returned. If the
// user's permissions change, the cached token must be deleted in order to create
// a new one.
func Token(config *oauth2.Config) (*oauth2.Token, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return getFromWeb(config)
	}

	token, err := getFromFile(cachePath)
	if err != nil {
		token, err = getFromWeb(config)
		if err != nil {
			return nil, err
		}

		// Even if saving fails, return the generated token.
		saveToFile(cachePath, token)
	}

	return token, nil
}

// CacheToken generates a new OAuth2 Token and caches it locally. It will return
// an error if the caching fails.
func CacheToken(config *oauth2.Config) error {
	token, err := getFromWeb(config)
	if err != nil {
		return err
	}

	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	return saveToFile(cachePath, token)
}

const (
	// credentialsDir is the name of the directory containing the credentials.
	credentialsDir = ".credentials"

	// tokenFile is the name of the file containing the authentication token.
	tokenFile = "sheets.googleaapis.com-reader.json"

	// ownerPermissions marks read, write and execute permissions for the owner.
	ownerPermissions = 0700
)

// getCachePath constructs the path to the local credentials cache.
func getCachePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Unable to get user info. %v", err)
	}

	tokenDir := filepath.Join(usr.HomeDir, credentialsDir)
	if err := os.MkdirAll(tokenDir, ownerPermissions); err != nil {
		return "", err
	}

	return filepath.Join(tokenDir, url.QueryEscape(tokenFile)), nil
}

// getFromWeb uses the Config to request a Token. It requires additional user
// input. It returns the retrieved Token.
func getFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Authentication link: \n%v\n\n", authURL)

	browser.OpenURL(authURL)

	fmt.Print("Type the authorization code: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}

	return config.Exchange(oauth2.NoContext, code)
}

// getFromFile retrieves a Token from a given file path.
func getFromFile(filepath string) (*oauth2.Token, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(file).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

// saveToFile uses a file path to create a file and store the token in it.
func saveToFile(filepath string, token *oauth2.Token) error {
	file, err := os.OpenFile(
		filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, ownerPermissions)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(token)
}
