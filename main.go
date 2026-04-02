
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SentimentResult struct to hold the sentiment analysis result
type SentimentResult struct {
	Text      string `json:"text"`
	Sentiment string `json:"sentiment"`
	Score     float64 `json:"score"`
}

// analyzeSentiment performs a very basic sentiment analysis
// In a real application, this would involve more sophisticated NLP techniques
func analyzeSentiment(text string) (string, float64) {
	text = strings.ToLower(text)
	positiveWords := []string{"good", "great", "excellent", "happy", "love", "positive"}
	negativeWords := []string{"bad", "terrible", "awful", "sad", "hate", "negative"}

	positiveScore := 0
	negativeScore := 0

	words := strings.Fields(text)
	for _, word := range words {
		for _, pWord := range positiveWords {
			if word == pWord {
				positiveScore++
			}
		}
		for _, nWord := range negativeWords {
			if word == nWord {
				negativeScore++
			}
		}
	}

	if positiveScore > negativeScore {
		return "positive", float64(positiveScore - negativeScore)
	} else if negativeScore > positiveScore {
		return "negative", float64(negativeScore - positiveScore)
	} else {
		return "neutral", 0.0
	}
}

// sentimentHandler handles HTTP requests for sentiment analysis
func sentimentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body to get the text
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}
	text := r.FormValue("text")

	if text == "" {
		http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
		return
	}

	// Perform sentiment analysis
	sentiment, score := analyzeSentiment(text)

	// Create a SentimentResult object
	result := SentimentResult{
		Text:      text,
		Sentiment: sentiment,
		Score:     score,
	}

	// Convert the result to JSON and send it as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	// Register the sentiment handler for the /analyze endpoint
	http.HandleFunc("/analyze", sentimentHandler)

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Starting sentiment analysis server on port %s
", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
