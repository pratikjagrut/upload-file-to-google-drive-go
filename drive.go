package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	credFile   = "credentials.json"
	tokenFile  = "token.json"
	uploadFile = "example.txt"
)

func UploadFileToDrive() error {
	// Load credentials file
	credentialsFile, err := os.ReadFile(credFile)
	if err != nil {
		return fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// Create oauth2 Config with credentials
	config, err := google.ConfigFromJSON(credentialsFile, drive.DriveScope)
	if err != nil {
		return fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	// Create oauth2 token
	token, err := tokenFromFile(tokenFile)
	if err != nil {
		token, err = getTokenFromWeb(config)
		if err != nil {
			return err
		}
	}

	if token.Valid() && (token.Expiry.Before(time.Now()) || token.Expiry.IsZero()) {
		token, err = refreshToken(config, token)
		if err != nil {
			return err
		}
	}

	// Create drive client with oauth2 token
	client := config.Client(context.Background(), token)
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("Unable to retrieve Drive client: %v", err)
	}

	// Open the file to upload
	file, err := os.Open(uploadFile)
	if err != nil {
		return fmt.Errorf("Unable to open file: %v", err)
	}
	defer file.Close()

	// Create the file metadata
	fileMetadata := &drive.File{
		Name: uploadFile,
	}

	// Upload the file to Google Drive
	_, err = srv.Files.Create(fileMetadata).Media(file).Do()
	if err != nil {
		return fmt.Errorf("Unable to upload file: %v", err)
	}

	fmt.Println("File uploaded successfully")
	return nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve token from web: %v", err)
	}
	err = saveToken(tokenFile, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func saveToken(file string, token *oauth2.Token) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	return nil
}

// refreshToken refreshes an OAuth2 token and saves it to a file
func refreshToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	// If it is, refresh the token
	newToken, err := config.TokenSource(context.Background(), token).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %v", err)
	}
	// Save the new token to a file
	err = saveToken(tokenFile, newToken)
	if err != nil {
		return nil, err
	}
	return newToken, nil
}
