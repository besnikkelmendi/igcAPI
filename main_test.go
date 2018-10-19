package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_postHandler_NotFound(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(Handler2))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != 400 {
		t.Errorf("Expected StatusNotFound %d, received %d. ", 400, resp.StatusCode)
		return
	}

}

func Test_postHandler_Delete_Method(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(postHANDLER1))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != 400 {
		t.Errorf("Expected StatusNotFound %d, received %d. ", 400, resp.StatusCode)
		return
	}

}

func Test_postHandler_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(postHANDLER1))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/randomTextJudtForTesting/",
		ts.URL + "/randommm/1sd5asd/",
		ts.URL + "/test",
		ts.URL + "/test/44a",
		ts.URL + "/45/45",
		ts.URL + "/12/asd",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}

func Test_Handler0_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(Handler))
	defer ts.Close()

	testCases := []string{
		ts.URL + "/randomTextJudtForTesting",
		ts.URL + "/randommm/1sd5asd/",
		ts.URL + "/test",
		ts.URL + "/test44a",
		ts.URL + "/45as45",
		ts.URL + "/12da",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, 200, resp.StatusCode)
			return
		}
	}
}

func Test_Handler2_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(Handler2))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/randomTextJudtForTesting",
		ts.URL + "/randommm/1sd5asd/",
		ts.URL + "/test",
		ts.URL + "/test44a",
		ts.URL + "/45as45",
		ts.URL + "/12da",
		ts.URL + "/12475",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, 400, resp.StatusCode)
			return
		}
	}
}

func Test_Handler3_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(Handler3))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/randomTextJudtForTesting",
		ts.URL + "/randommm/1sd5asd",
		ts.URL + "/test",
		ts.URL + "/test44a",
		ts.URL + "/45as45",
		ts.URL + "/12da",
		ts.URL + "/1221312",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != 400 {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, 400, resp.StatusCode)
			return
		}
	}
}
