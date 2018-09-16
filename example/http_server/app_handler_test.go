package main_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/valerykalashnikov/moments"
	main "github.com/valerykalashnikov/moments/example/http_server"
)

func TestAppHandler_ServeHTTP(t *testing.T) {
	counter := moments.NewMomentsCounter(time.Minute)

	handler := main.NewAppHandler(counter)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "1"
	actual := string(body)

	if string(body) != expected {
		t.Error("Expected: ", expected, " actual: ", actual)
	}

}
