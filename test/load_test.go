package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLoadTransaction(t *testing.T) {
	const (
		numRequests = 100
		targetTPS   = 25
		tolerance   = 0.02 // 2% tolerance
	)

	var wg sync.WaitGroup
	errors := make(chan error, numRequests)
	start := time.Now()

	ticker := time.NewTicker(time.Second / time.Duration(targetTPS))
	defer ticker.Stop()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			makeRequest(id, errors)
		}(i)
		<-ticker.C
	}

	wg.Wait()
	close(errors)

	elapsed := time.Since(start)
	actualTPS := float64(numRequests) / elapsed.Seconds()
	minAcceptableTPS := float64(targetTPS) * (1 - tolerance)

	for err := range errors {
		t.Error(err)
	}

	t.Logf("Completed %d requests in %v (%.2f TPS)", numRequests, elapsed, actualTPS)

	if actualTPS < minAcceptableTPS {
		t.Errorf("Failed to achieve target TPS. Got %.2f, want %.2f", actualTPS, minAcceptableTPS)
	}
}

func makeRequest(id int, errors chan<- error) {
	client := &http.Client{}
	payload := fmt.Sprintf(`{
		"state": "win",
		"amount": "1.00",
		"transactionId": "test-%d"
	}`, id)

	req, err := http.NewRequest("POST", "http://localhost:8080/user/1/transaction", strings.NewReader(payload))
	if err != nil {
		errors <- fmt.Errorf("request %d creation failed: %v", id, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Source-Type", "game")

	resp, err := client.Do(req)
	if err != nil {
		errors <- fmt.Errorf("request %d failed: %v", id, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errors <- fmt.Errorf("request %d: expected status OK, got %v: %s", id, resp.Status, string(body))
	}
}
