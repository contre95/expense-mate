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

type ExpenseGuess struct {
	Shop    string
	Amount  float64
	Date    time.Time
	Product string
}

func NewGuesser() (*Guesser, error) {
	return &Guesser{
		apiURL: "http://localhost:11434/api/generate",
		client: &http.Client{Timeout: 300 * time.Second},
	}, nil
}

func (g *Guesser) GuessExpense(imageData []byte) ([]ExpenseGuess, error) {
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	requestBody := map[string]interface{}{
		"model": "llama3.2-vision:11b-instruct-q4_K_M",
		"prompt": `STRICT INSTRUCTIONS:
1. Return ONLY a JSON array of transaction objects.
2. Each object MUST have EXACTLY these fields:
   - shop (string): The name of the shop or store.
   - amount (positive float): The total amount spent, as a positive number (e.g., 4.50, not +€4.50 or -€4.50).
   - date (string in YYYY-MM-DD format): The date of the transaction, in the exact format YYYY-MM-DD.
3. DO NOT include any nested groups/categories.
4. DO NOT add other fields or explanations.
5. DO NOT include currency symbols (e.g., €, $) or signs (e.g., +, -) in the amount.
6. DO NOT include time in the date field.
7. If the date or amount cannot be parsed, omit the transaction entirely.
GOOD EXAMPLE: 
[
  {
    "shop": "Coffee Shop",
    "amount": 4.50,
    "date": "2023-10-15"
  },
  {
    "shop": "Supermarket",
    "amount": 12.30,
    "date": "2023-10-16"
  }
]

Current receipt to parse:`,
		"stream": false,
		"format": "json",
		"images": []string{encodedImage},
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error encoding request: %v", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", g.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

var apiResponse struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("\n\n%s\n\n", apiResponse.Response)
	jsonStr, err := extractJSONArray(apiResponse.Response)
	if err != nil {
		return nil, err
	}

	var results []struct {
		Shop   string  `json:"shop"`
		Amount float64 `json:"amount"`
		Date   string  `json:"date"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &results); err != nil {
		return nil, fmt.Errorf("error parsing JSON array: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no items found in response")
	}

	guesses := make([]ExpenseGuess, 0, len(results))
	for _, result := range results {
		parsedDate, err := time.Parse("2006-01-02", result.Date)
		if err != nil {
			return nil, fmt.Errorf("error parsing date %q: %v", result.Date, err)
		}

		guesses = append(guesses, ExpenseGuess{
			Shop:    result.Shop,
			Amount:  result.Amount,
			Date:    parsedDate,
			Product: "AI Generated",
		})
	}

	return guesses, nil
}

func extractJSONArray(input string) (string, error) {
	re := regexp.MustCompile(`(?s)\[.*\]`)
	match := re.FindString(input)
	if match == "" {
		return "", fmt.Errorf("no JSON array found in response")
	}
	return match, nil
}
