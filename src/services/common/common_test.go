package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

func TestCommonHandler(t *testing.T) {
	config := viper.New()
	config.Set("runtime.environment", "test")

	commonSvc := New(config)

	req, err := http.NewRequest("GET", "/test-path", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()

	commonSvc.CommonHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := "Hello from test /test-path "
	if rr.Body.String() != expectedResponse {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}
