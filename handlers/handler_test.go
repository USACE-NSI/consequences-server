package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestCompute(t *testing.T) {
	// 1. Setup Echo and Handler
	e := echo.New()
	h := &Handler{}

	// 2. Read the data.json file you created earlier
	payload, err := os.ReadFile("/workspaces/consequences-server/data.json")
	if err != nil {
		t.Fatalf("Could not read data.json: %v", err)
	}

	// 3. Create a POST request with the JSON payload
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// 4. Create a ResponseRecorder to capture the output
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 5. Call the handler
	if err := h.Compute(c); err != nil {
		t.Errorf("Handler returned an error: %v", err)
	}

	// 6. Assertions using standard library logic
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rec.Code)
	}

	// Check if we actually got a response body back
	resBody, _ := io.ReadAll(rec.Body)
	if len(resBody) == 0 {
		t.Error("Expected a response body, but got an empty string")
	}

	fmt.Printf("Response received: %s", string(resBody))
}
