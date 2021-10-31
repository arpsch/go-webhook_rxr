package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create a request %v", err)
	}

	apph := appHandlers{}
	rec := httptest.NewRecorder()
	apph.IndexHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("accepted status 200, got %v", rec.Code)
	}

	fmt.Printf("### %v ######\n", rec.Result().Cookies())
	fmt.Printf("### %v ######\n", rec.Body)
	/*
	   if !strings.Contains(rec.Body.String(), "hello world") {
	       t.Errorf("unexpected body in response %q", rec.Body.String())
	   }
	*/

}
