package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/usace-nsi/consequences-server/handlers"
)

func TestComputeEndpoint(t *testing.T) {
	// 1. Setup the real Echo instance exactly like your main.go
	e := echo.New()
	handler := handlers.Handler{}
	e.POST("/compute", handler.Compute)

	// 2. Start the real local server
	ts := httptest.NewServer(e)
	defer ts.Close()

	body, err := os.ReadFile("/workspaces/consequences-server/data.json")
	if err != nil {
		t.Fatalf("Could not read data.json: %v", err)
	}
	// 4. Execute a real POST request
	resp, err := http.Post(ts.URL+"/compute", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to send POST: %v", err)
	}
	defer resp.Body.Close()

	// 5. Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200; got %d", resp.StatusCode)
	}

	// Since you use a GeoJsonResultsWriter, let's verify we got some data back
	respBody, _ := io.ReadAll(resp.Body)
	if len(respBody) == 0 {
		t.Error("Expected GeoJSON response body, but got empty response")
	}

	fmt.Printf("Response received: %s", string(respBody[:200])) // Log first 50 chars
}

func BenchmarkComputeEndpoint_Stress(b *testing.B) {
	// 1. Setup server once for the entire benchmark
	e := echo.New()
	handler := handlers.Handler{}
	e.POST("/compute", handler.Compute)
	ts := httptest.NewServer(e)
	defer ts.Close()

	// 2. Pre-load the real data into memory so file I/O doesn't bottleneck the test
	body, err := os.ReadFile("/workspaces/consequences-server/data.json")
	if err != nil {
		b.Fatalf("Could not read data.json: %v", err)
	}

	b.ResetTimer() // Don't count setup time

	// 3. Run in parallel to simulate multiple concurrent users
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Post(ts.URL+"/compute", "application/json", bytes.NewBuffer(body))
			if err != nil {
				b.Errorf("Request failed: %v", err)
				continue
			}

			// Always drain and close the body to avoid leaking file descriptors
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected 200, got %d", resp.StatusCode)
			}
			respBody, _ := io.ReadAll(resp.Body)
			fmt.Printf("Response received: %s", string(respBody[:200]))
		}
	})
}

func BenchmarkRemoteComputeStress(b *testing.B) {
	// 1. Point to the real remote URL
	targetURL := "https://consequences.hecdev.net/compute"

	// 2. Pre-load your data.json once
	body, err := os.ReadFile("/workspaces/consequences-server/data.json")
	if err != nil {
		b.Fatalf("Could not read data.json: %v", err)
	}

	b.ResetTimer()

	// 3. Parallel execution simulates high concurrent traffic
	b.RunParallel(func(pb *testing.PB) {
		client := &http.Client{} // Shared client is thread-safe
		for pb.Next() {
			resp, err := client.Post(targetURL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				b.Errorf("Request to remote failed: %v", err)
				continue
			}

			// Crucial: Drain and close the body to avoid hitting open file limits
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Remote returned non-200: %d", resp.StatusCode)
			}
		}
	})
}
