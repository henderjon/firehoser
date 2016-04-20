package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func sendOmnilogger() {
	// omnilogger (POST http://localhost:8080/)

	body := strings.NewReader(`Lorem	ipsum	dolor	sit	amet	consectetur	adipiscing	elit	Sed
felis	ligula	laoreet	at	sapien	a	sodales	facilisis	massa
Nulla	eleifend	ac	purus	auctor	consectetur	Morbi	imperdiet	dictum
ex	in	imperdiet	Quisque	et	mauris	neque	Praesent	at
nibh	venenatis	egestas	ipsum	ac	convallis	tortor	Sed	cursus
lectus	odio	et	tempor	risus	malesuada	eu	Praesent	nulla
turpis	hendrerit	nec	orci	quis	gravida	pulvinar	est	Vestibulum
congue	tellus	et	congue	pretium	Nunc	posuere	consequat	molestie`)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://localhost:8080/", body)

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


