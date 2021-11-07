package importing

import (
	"context"
	"encoding/json"
	"expenses/pkg/app"
	"expenses/pkg/domain/expense"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type SheetsImporertResp struct {
	ID  expense.CategoryID
	Msg string
}

type SheetsImporertReq struct {
}

// The createCategory use case creates a category for a expense
type ImportFromSheetsUseCase struct {
	logger   app.Logger
	expenses expense.Expenses
}

func NewSheetsImporter(l app.Logger, e expense.Expenses) *ImportFromSheetsUseCase {
	return &ImportFromSheetsUseCase{l, e}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Create use cases function creates a new category
func (u *ImportFromSheetsUseCase) Import(req SheetsImporertReq) (*SheetsImporertResp, error) {
	panic("Implement me!")
}
