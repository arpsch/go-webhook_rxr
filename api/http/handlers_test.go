package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzHandler(t *testing.T) {

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:9999/",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create a request %v", err)
	}

	rh := receiverHandlers{}
	rec := httptest.NewRecorder()
	rh.HealthzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("accepted status 200, received status %v", rec.Code)
	}

	body, err := ioutil.ReadAll(rec.Body)

	if err != nil {
		log.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("accepted response OK, received status %s", body)
	}
}
