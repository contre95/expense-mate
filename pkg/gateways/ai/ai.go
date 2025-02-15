package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type Guesser struct {
	apiURL string
	client *http.Client
}

func NewGuesser() (*Guesser, error) {
	return &Guesser{
		apiURL: "http://localhost:11434/api/generate",
		client: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (g *Guesser) GuessExpense(imageData []byte) (shop string, amount float64, date time.Time, product string, err error) {
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	requestBody := map[string]interface{}{
		"model": "llava-llama3",
		"prompt": `STRICT INSTRUCTIONS:
  1. Return ONLY a JSON array of transaction objects
  2. Each object MUST have EXACTLY these fields:
     - shop (string)
     - amount (positive float)
     - date (YYYY-MM-DD, no negative dates)
     - product (string)
  3. DO NOT include any nested groups/categories
  4. DO NOT add other fields or explanations

  BAD EXAMPLE: {"Transactions": [...]}
  GOOD EXAMPLE: [
    {
      "shop": "Coffee Shop",
      "amount": 4.50,
      "date": "2023-10-15",
      "product": "Cappuccino"
    }
  ]
  Current receipt to parse:`,
		"stream": false,
		"format": "json",
		"images": []string{encodedImage},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error encoding request: %v", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", g.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, time.Time{}, "", fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var apiResponse struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("\n\n%s\n\n", apiResponse.Response)
	jsonStr, err := extractJSONArray(apiResponse.Response)
	if err != nil {
		return "", 0, time.Time{}, "", err
	}

	var results []struct {
		Shop    string  `json:"shop"`
		Amount  float64 `json:"amount"`
		Date    string  `json:"date"`
		Product string  `json:"product"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &results); err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error parsing JSON array: %v", err)
	}

	if len(results) == 0 {
		return "", 0, time.Time{}, "", fmt.Errorf("no items found in response")
	}

	// Take the first item from the array
	result := results[0]
	parsedDate, err := time.Parse("2006-01-02", result.Date)
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error parsing date: %v", err)
	}

	return result.Shop, result.Amount, parsedDate, result.Product, nil
}

// Updated to extract JSON arrays
func extractJSONArray(input string) (string, error) {
	re := regexp.MustCompile(`(?s)\[.*\]`)
	match := re.FindString(input)
	if match == "" {
		return "", fmt.Errorf("no JSON array found in response")
	}
	return match, nil
}
package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type Guesser struct {
	apiURL string
	client *http.Client
}

func NewGuesser() (*Guesser, error) {
	return &Guesser{
		apiURL: "http://localhost:11434/api/generate",
		client: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (g *Guesser) GuessExpense(imageData []byte) (shop string, amount float64, date time.Time, product string, err error) {
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	requestBody := map[string]interface{}{
		"model": "llava-llama3",
		"prompt": `STRICT INSTRUCTIONS:
  1. Return ONLY a JSON array of transaction objects
  2. Each object MUST have EXACTLY these fields:
     - shop (string)
     - amount (positive float)
     - date (YYYY-MM-DD, no negative dates)
     - product (string)
  3. DO NOT include any nested groups/categories
  4. DO NOT add other fields or explanations

  BAD EXAMPLE: {"Transactions": [...]}
  GOOD EXAMPLE: [
    {
      "shop": "Coffee Shop",
      "amount": 4.50,
      "date": "2023-10-15",
      "product": "Cappuccino"
    }
  ]
  Current receipt to parse:`,
		"stream": false,
		"format": "json",
		"images": []string{encodedImage},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error encoding request: %v", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", g.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, time.Time{}, "", fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var apiResponse struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("\n\n%s\n\n", apiResponse.Response)
	jsonStr, err := extractJSONArray(apiResponse.Response)
	if err != nil {
		return "", 0, time.Time{}, "", err
	}

	var results []struct {
		Shop    string  `json:"shop"`
		Amount  float64 `json:"amount"`
		Date    string  `json:"date"`
		Product string  `json:"product"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &results); err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error parsing JSON array: %v", err)
	}

	if len(results) == 0 {
		return "", 0, time.Time{}, "", fmt.Errorf("no items found in response")
	}

	// Take the first item from the array
	result := results[0]
	parsedDate, err := time.Parse("2006-01-02", result.Date)
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error parsing date: %v", err)
	}

	return result.Shop, result.Amount, parsedDate, result.Product, nil
}

// Updated to extract JSON arrays
func extractJSONArray(input string) (string, error) {
	re := regexp.MustCompile(`(?s)\[.*\]`)
	match := re.FindString(input)
	if match == "" {
		return "", fmt.Errorf("no JSON array found in response")
	}
	return match, nil
}
