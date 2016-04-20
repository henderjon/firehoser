curl -X "POST" "http://localhost:8080/" \
	-H "X-Omnilogger-Stream: test-header-value" \
	-H "Authorization: Bearer test-password-hash" \
	-H "Content-Type: text/plain" \
	-d $'Lorem	ipsum	dolor	sit	amet	consectetur	adipiscing	elit	Sed
felis	ligula	laoreet	at	sapien	a	sodales	facilisis	massa
Nulla	eleifend	ac	purus	auctor	consectetur	Morbi	imperdiet	dictum
ex	in	imperdiet	Quisque	et	mauris	neque	Praesent	at
nibh	venenatis	egestas	ipsum	ac	convallis	tortor	Sed	cursus
lectus	odio	et	tempor	risus	malesuada	eu	Praesent	nulla
turpis	hendrerit	nec	orci	quis	gravida	pulvinar	est	Vestibulum
congue	tellus	et	congue	pretium	Nunc	posuere	consequat	molestie'


