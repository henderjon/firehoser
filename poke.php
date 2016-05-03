<?php

for($i = 0; $i < 150; $i += 1){

// Get cURL resource
$ch = curl_init();

// Set url
curl_setopt($ch, CURLOPT_URL, 'http://127.0.0.1:8080/log');

// Set method
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'POST');

// Set options
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

// Set headers
curl_setopt($ch, CURLOPT_HTTPHEADER, [
  "Content-Type: text/plain",
 ]
);
// Create body
$body = 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam id turpis sit amet nibh tempus fringilla. Vivamus lacinia metus et neque dignissim egestas eu non sem. Phasellus pretium augue ultrices, tristique dui vel, euismod est. Maecenas egestas mauris quis diam maximus laoreet. Curabitur mattis, diam sed mollis posuere, felis ipsum rhoncus nulla, non gravida metus ipsum lobortis orci. Mauris quis tellus et enim elementum fermentum.
';

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

}
