package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func sendOmnilogger() {
	// omnilogger (POST http://localhost:8080/)

	body := strings.NewReader(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam id turpis sit amet nibh tempus fringilla. Vivamus lacinia metus et neque dignissim egestas eu non sem. Phasellus pretium augue ultrices, tristique dui vel, euismod est. Maecenas egestas mauris quis diam maximus laoreet. Curabitur mattis, diam sed mollis posuere, felis ipsum rhoncus nulla, non gravida metus ipsum lobortis orci. Mauris quis tellus et enim elementum fermentum.
`)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://localhost:8080/log", body)

	// Headers
	req.Header.Add("X-Omnilogger-Stream", "test-header-value")
	req.Header.Add("Authorization", "Bearer test-password-hash")
	req.Header.Add("Content-Type", "text/plain")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}


