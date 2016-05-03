curl -X "POST" "http://localhost:8080/log" \
	-H "X-Omnilogger-Stream: test-header-value" \
	-H "Authorization: Bearer test-password-hash" \
	-H "Content-Type: text/plain" \
	-d $'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam id turpis sit amet nibh tempus fringilla. Vivamus lacinia metus et neque dignissim egestas eu non sem. Phasellus pretium augue ultrices, tristique dui vel, euismod est. Maecenas egestas mauris quis diam maximus laoreet. Curabitur mattis, diam sed mollis posuere, felis ipsum rhoncus nulla, non gravida metus ipsum lobortis orci. Mauris quis tellus et enim elementum fermentum.
'


