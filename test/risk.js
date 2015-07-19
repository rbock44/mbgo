function(request, response) {
	var responseId=("000000000000000000000000" + Math.round(Math.random()*10000000000000)).slice(-24);
	var requestId=request.body.match(/<request-id>.*request-id>/);
	var ipaddress=request.body.match(/<ip-address>.*ip-address>/);

	response.body=response.body.replace("$RESPONSE_ID", "<responseId>"+responseId+"</responseId>");
	response.body=response.body.replace("$REQUEST_ID", requestId);
	response.body=response.body.replace("$IPADDRESS", ipaddress); 
}
