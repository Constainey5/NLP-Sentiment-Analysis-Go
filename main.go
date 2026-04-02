package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

var sentimentLexicon = map[string]float64{
	"good": 1.0, "great": 1.2, "excellent": 1.5, "happy": 0.8, "love": 1.3, "positive": 1.0,
	"bad": -1.0, "terrible": -1.2, "poor": -0.8, "sad": -0.9, "hate": -1.3, "negative": -1.0,
	"neutral": 0.0, "ok": 0.1, "fine": 0.2,
}

func analyzeSentiment(text string) (string, float64) {
	words := strings.Fields(strings.ToLower(text))
	var totalScore float64
	var wg sync.WaitGroup
	scoreChan := make(chan float64, len(words))

	for _, word := range words {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			if score, ok := sentimentLexicon[w]; ok {
				scoreChan <- score
			} else {
				scoreChan <- 0.0
			}
		}(word)
	}

	wg.Wait()
	close(scoreChan)

	for score := range scoreChan {
		totalScore += score
	}

	sentiment := "Neutral"
	if totalScore > 0.5 {
		sentiment = "Positive"
	} else if totalScore < -0.5 {
		sentiment = "Negative"
	}

	return sentiment, totalScore
}

type SentimentRequest struct {
	Text string `json:"text"`
}

type SentimentResponse struct {
	Text      string  `json:"text"`
	Sentiment string  `json:"sentiment"`
	Score     float64 `json:"score"`
}

func sentimentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var req SentimentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sentiment, score := analyzeSentiment(req.Text)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SentimentResponse{Text: req.Text, Sentiment: sentiment, Score: score})
}

func main() {
	http.HandleFunc("/analyze", sentimentHandler)
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
