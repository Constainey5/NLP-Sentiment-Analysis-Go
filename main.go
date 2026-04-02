
package main

import (
	"fmt"
	"log"
	"net/http"
)

func sentimentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
		return
	}

	// Placeholder for actual sentiment analysis logic
	sentiment := analyzeSentiment(text)

	fmt.Fprintf(w, "{"text": "%s", "sentiment": "%s"}", text, sentiment)
}

func analyzeSentiment(text string) string {
	// Simple placeholder logic
	if len(text) > 20 && text[0] == 'I' {
		return "positive"
	} else if len(text) < 10 {
		return "neutral"
	} else {
		return "negative"
	}
}

func main() {
	http.HandleFunc("/analyze", sentimentHandler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
