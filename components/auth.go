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

const (
	// credentialsDir is the name of the directory containing the credentials.
	credentialsDir = ".credentials"

	// tokenFile is the name of the file containing the authentication token.
	tokenFile = "sheets.googleaapis.com-reader.json"

	// ownerPermissions marks read, write and execute permissions for the owner.
	ownerPermissions = 0700
)

// Token uses a OAuth2 Config to retrieve a Token for the current user. The token
// will be cached locally. If the user's permissions change, the cached token must
// be deleted in order to create a new one.
func Token(config *oauth2.Config) (*oauth2.Token, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("Unable to get user info. %v", err)
	}

	tokenDir := filepath.Join(usr.HomeDir, credentialsDir)
	if err := os.MkdirAll(tokenDir, ownerPermissions); err != nil {
		return nil, err
	}

	cacheFile := filepath.Join(tokenDir, url.QueryEscape(tokenFile))

	token, err := getFromFile(cacheFile)
	if err != nil {
		token, err = getFromWeb(config)
		if err != nil {
			return nil, err
		}

		if err := saveToFile(cacheFile, token); err != nil {
			return nil, err
		}
	}

	return token, nil
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
