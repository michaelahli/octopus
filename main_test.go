package main

import (
	"net/http"
	"testing"

	"github.com/michaelahli/octopus/src/services"
	"github.com/spf13/viper"
)

type MockService struct{}

func (s *MockService) HandleBooks(w http.ResponseWriter, r *http.Request) {}

func (s *MockService) CommonHandler(w http.ResponseWriter, r *http.Request) {}

func TestStartServer(t *testing.T) {
	config := viper.New()
	config.Set("server.port", "8080")

	mockService := &MockService{}

	go serveTest(config, mockService)

	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatalf("Error making request to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func serveTest(cfg *viper.Viper, svc services.Services) {
	port := cfg.GetString("server.port")

	http.HandleFunc("/book", svc.HandleBooks)
	http.HandleFunc("/", svc.CommonHandler)

	port = ":" + port
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			panic(err)
		}
	}()
}
