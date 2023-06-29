package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// Define the task to be performed on a sentence
func processSentence(sentence string) map[string]int {
	words := strings.Fields(sentence)
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}
	return counts
}

// Define the worker that processes sentences
func worker(sentences <-chan string, results chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	for sentence := range sentences {
		result := processSentence(sentence)
		results <- result
	}
}

func main() {
	// Open the CSV file
	file, err := os.Open("misc/sentences.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Create channels for communication between goroutines
	sentenceChannel := make(chan string)
	resultChannel := make(chan map[string]int)
	done := make(chan bool)

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Launch the goroutine to read sentences from the CSV file
	go func() {
		defer close(sentenceChannel)
		// Read all records from the CSV file
		records, err := reader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		// Extract the sentences from the records (adjust column index if needed)
		for _, record := range records {
			sentence := record[0]
			sentenceChannel <- sentence
		}
	}()

	// Launch worker goroutines
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(sentenceChannel, resultChannel, &wg)
	}

	// Create a map to store the word counts
	wordCounts := make(map[string]int)

	// Process results
	go func() {
		wg.Wait()
		close(resultChannel) // Close the resultChannel after all workers finish
		done <- true
	}()

	// Distribute sentences among workers using round-robin assignment
	go func() {
		// Create a slice of channels to hold the worker channels
		workerChannels := make([]chan string, numWorkers)
		for i := 0; i < numWorkers; i++ {
			workerChannels[i] = make(chan string)
		}

		// Round-robin assignment of sentences to worker channels
		i := 0
		for sentence := range sentenceChannel {
			workerChannels[i] <- sentence
			i = (i + 1) % numWorkers
		}

		// Close worker channels
		for _, ch := range workerChannels {
			close(ch)
		}
	}()

	// Collect and process results
	for result := range resultChannel {
		// Update the word counts based on each result
		for word, count := range result {
			wordCounts[word] += count
		}
	}

	// Wait for all processing to complete
	<-done

	// Print the word counts
	for word, count := range wordCounts {
		fmt.Printf("%s: %d\n", word, count)
	}

	fmt.Println("Processing complete")
}
