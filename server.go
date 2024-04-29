package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	dataChan := make(chan string)

	go func() {
		messages := []string{
			"Hello, how can I assist you today?",
			"I'm an AI language model trained to help with various tasks.",
			"Feel free to ask me anything, and I'll do my best to provide a helpful response.",
			"I can help with writing, analysis, coding, math, and more.",
			"Let me know if you have any specific questions or topics you'd like to discuss.",
		  "I have a vast knowledge base covering various subjects, so feel free to explore any topic you're interested in.",
			"Whether you need help with problem-solving, creative writing, or research, I'm here to assist you.",
			"I can provide explanations, generate ideas, and offer insights on a wide range of topics.",
			"My goal is to provide accurate and informative responses to help you learn and achieve your objectives.",
			"If you have any doubts or need clarification on anything, please don't hesitate to ask.",
		}

		for _, message := range messages {
			for _, char := range message {
				dataChan <- string(char)
				time.Sleep(time.Duration(rand.Intn(300)+50) * time.Millisecond)
			}
			dataChan <- "\n"
			time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
		}
	}()

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		for {
			select {
			case char := <-dataChan:
				if char == "\n" {
					fmt.Fprintf(w, "data: %s\n\n", char)
				} else {
					fmt.Fprintf(w, "data: %s\n\n", char)
				}
				flusher.Flush()
			case <-r.Context().Done():
				
        return
			}
		}
	})

	log.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
