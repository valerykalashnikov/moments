package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestApp(t *testing.T) {

	for i := 0; i < 60; i++ {
		time.Sleep(1 * time.Second)

		_, err := http.Get("http://localhost:8080")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

	}
	time.Sleep(time.Second)
	_, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	time.Sleep(time.Second)
	res, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// after making requests every 1 second returned value should be constantly 60
	expected := "60"
	actual := string(body)

	if string(body) != expected {
		t.Error("Expected: ", expected, " actual: ", actual)
	}

	fmt.Println(string(body))

}
