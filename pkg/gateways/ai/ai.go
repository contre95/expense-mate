package ai

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/otiai10/gosseract/v2"
)

type Guesser struct {
	ocrClient *gosseract.Client
}

func NewGuesser() (*Guesser, error) {
	client := gosseract.NewClient()
	err := client.SetLanguage("eng")
	if err != nil {
		return nil, fmt.Errorf("failed to set OCR language: %w", err)
	}
	return &Guesser{ocrClient: client}, nil
}

func (g *Guesser) GuessExpense(imageData []byte) (shop string, amount float64, date time.Time, product string, err error) {
	g.ocrClient.SetImageFromBytes(imageData)
	text, err := g.ocrClient.Text()
	if err != nil {
		return "", 0, time.Time{}, "", fmt.Errorf("OCR failed: %w", err)
	}

	// Preprocess text for better parsing
	lines := strings.Split(text, "\n")
	cleanedLines := make([]string, 0)
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			cleanedLines = append(cleanedLines, trimmed)
		}
	}

	// Debug: Print OCR results
	fmt.Println("OCR Text:\n", strings.Join(cleanedLines, "\n"))

	// Parse information
	amount, err = extractAmount(cleanedLines)
	if err != nil {
		return "", 0, time.Time{}, "", err
	}

	date = extractDate(cleanedLines)
	shop = extractShop(cleanedLines)
	product = extractProduct(cleanedLines)

	return shop, amount, date, product, nil
}

func extractAmount(lines []string) (float64, error) {
	// Look for the main total amount (last € value)
	var amountStr string
	re := regexp.MustCompile(`€\s*([+-]?\d+[\.,]\d{2})`)

	for i := len(lines) - 1; i >= 0; i-- {
		if matches := re.FindStringSubmatch(lines[i]); len(matches) > 1 {
			amountStr = strings.Replace(matches[1], ",", ".", 1)
			break
		}
	}

	if amountStr == "" {
		return 0, fmt.Errorf("no amount found")
	}

	var amount float64
	_, err := fmt.Sscanf(amountStr, "%f", &amount)
	return amount, err
}

func extractDate(lines []string) time.Time {
	// Look for date keywords or relative dates
	now := time.Now()
	dateFormats := []string{
		"2/1/2006",
		"2006-01-02",
		"January 2, 2006",
	}

	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "today") {
			return now
		}
		if strings.Contains(lower, "yesterday") {
			return now.AddDate(0, 0, -1)
		}

		for _, format := range dateFormats {
			if _, err := time.Parse(format, line); err == nil {
				if t, err := time.Parse(format, line); err == nil {
					return t
				}
			}
		}
	}
	return now
}

func extractShop(lines []string) string {
	// Look for lines that appear before time stamps
	timePattern := regexp.MustCompile(`\d{1,2}:\d{2}\s*[AP]M`)

	for i, line := range lines {
		if timePattern.MatchString(line) && i > 0 {
			return lines[i-1]
		}
	}
	return "Unknown Shop"
}

func extractProduct(lines []string) string {
	// Look for lines starting with "For" or after shop name
	for i, line := range lines {
		if strings.HasPrefix(line, "For ") && i+1 < len(lines) {
			return lines[i+1]
		}
		if strings.Contains(line, "€") && i > 0 {
			return lines[i-1]
		}
	}
	return "General Purchase"
}

