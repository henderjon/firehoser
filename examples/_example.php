<?php

// Get cURL resource
$ch = curl_init();

// Set url
curl_setopt($ch, CURLOPT_URL, 'http://localhost:8080/');

// Set method
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'POST');

// Set options
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

// Set headers
curl_setopt($ch, CURLOPT_HTTPHEADER, [
  "X-Omnilogger-Stream: test-header-value",
  "Authorization: Bearer test-password-hash",
  "Content-Type: text/plain",
 ]
);
// Create body
$body = 'Lorem	ipsum	dolor	sit	amet	consectetur	adipiscing	elit	Sed
felis	ligula	laoreet	at	sapien	a	sodales	facilisis	massa
Nulla	eleifend	ac	purus	auctor	consectetur	Morbi	imperdiet	dictum
ex	in	imperdiet	Quisque	et	mauris	neque	Praesent	at
nibh	venenatis	egestas	ipsum	ac	convallis	tortor	Sed	cursus
lectus	odio	et	tempor	risus	malesuada	eu	Praesent	nulla
turpis	hendrerit	nec	orci	quis	gravida	pulvinar	est	Vestibulum
congue	tellus	et	congue	pretium	Nunc	posuere	consequat	molestie';

// Set body
curl_setopt($ch, CURLOPT_POST, 1);
curl_setopt($ch, CURLOPT_POSTFIELDS, $body);

// Send the request & save response to $resp
$resp = curl_exec($ch);

if(!$resp) {
  die('Error: "' . curl_error($ch) . '" - Code: ' . curl_errno($ch));
} else {
  echo "Response HTTP Status Code : " . curl_getinfo($ch, CURLINFO_HTTP_CODE);
  echo "\nResponse HTTP Body : " . $resp;
}

// Close request to clear up some resources
curl_close($ch);


