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
	// Encode image to base64
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	// Prepare request payload
	requestBody := map[string]interface{}{
		"model": "llava-llama3",
		"prompt": `Extract receipt information as JSON with these fields: 
                  "shop" (string), "amount" (float), 
                  "date" (YYYY-MM-DD format). 
                  Return ONLY valid JSON without any additional text.`,
		"stream": false,
		"format": "json",
		"images": []string{encodedImage},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error encoding request: %v", err)
	}

	// Create API request
	req, err := http.NewRequestWithContext(context.Background(), "POST", g.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := g.client.Do(req)
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, time.Time{}, "", fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	// Parse response
	var apiResponse struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error decoding response: %v", err)
	}

	// Extract JSON from response
	fmt.Printf("\n\n%s\n\n", apiResponse.Response)
	jsonStr, err := extractJSON(apiResponse.Response)
	if err != nil {
		return "", 0, time.Time{}, "", err
	}

	// Parse extracted JSON
	var result struct {
		Shop    string  `json:"shop"`
		Amount  float64 `json:"amount"`
		Date    string  `json:"date"`
		Product string  `json:"product"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("error parsing JSON: %v", err)
	}

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", result.Date)
	if err != nil {
		// return "", 0, time.Time{}, "", fmt.Errorf("error parsing date: %v", err)
		fmt.Println(result.Date)
		parsedDate = time.Now()
	}

	return result.Shop, result.Amount, parsedDate, result.Product, nil
}

// Helper function to extract JSON from response text
func extractJSON(input string) (string, error) {
	re := regexp.MustCompile(`(?s)\{.*\}`)
	match := re.FindString(input)
	if match == "" {
		return "", fmt.Errorf("no JSON found in response")
	}
	return match, nil
}
