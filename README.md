# mbgo
Mountebank preprocessor to generate data for a single JSON imposter post to mountebank server.

Mountebank http://www.mbtest.org/ is a real cool HTTP mock which supports Java Script as a mighty simulation language. Mountebank supports configuration files which do allow more human readable format. The configuration files are limited for those who deploy Mountebank and the test scripts local. But when you want to control Mountebank remotely through some test clients you need to stick to some really crazy JSON style content.

To make my life easier deploying Mountebank with the application server and running the test suites on Jenkins I came up with a small GOLANG tool that adds some preprocessing capabilities.

Let us start with a Mountebank imposter POST request JSON body. Instead of adding all required information we use the macro INCLUDE.

```
{
	port": 9005,
        "protocol": "http",
        "stubs": [{
                "responses": [{
                        "is": { "statusCode": 200, "headers": { "Content-Type": "application/xml" }, "body": "INCLUDE=./ip_verification.xml" },
                        "_behaviors": { "decorate": "INCLUDE=./risk.js" }
                }],
                "predicates": [{
                        "and": [
                                { "equals": { "path": "/risk/verify", "method": "POST", "headers": { "Content-Type": "application/xml" } } },
                                { "contains": { "body": "<ip-address>" } }
                        ]
                }]
        },
```


The sample sets up a risk engine web service that evaluates the risk for a certain customer transaction.

The web service uses XML content and needs some java script logic to take over some tags from the request to the response and generate a random response id.

The xml file can be written in a readable format:

```
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<risk-response>
        <application>
                $RESPONSE_ID
                $REQUEST_ID
        </application>
        <decision>
                <return-code>0</return-code>
                <return-message>Successful system entry.</return-message>
        </decision>
        <verifications>
                <ip-verification>
                        <return-code>0</return-code>
                        <return-message>Successful system entry.</return-message>
                        <country>DE</country>
                        <country-confidence-factor>99</country-confidence-factor>
                        $IPADDRESS
                </ip-verification>
        </verifications>
</risk-response>
```

$RESPONSE_ID, $REQUEST_ID and $IPADDRESS are placeholders for the javascript code that inserts values there.

The javascript part can be written in a separate file:

```
function(request, response) {
        var responseId=("000000000000000000000000" + Math.round(Math.random()*10000000000000)).slice(-24);
        var requestId=request.body.match(/<request-id>.*request-id>/);
        var ipaddress=request.body.match(/<ip-address>.*ip-address>/);

        response.body=response.body.replace("$RESPONSE_ID", "<responseId>"+responseId+"</responseId>");
        response.body=response.body.replace("$REQUEST_ID", requestId);
        response.body=response.body.replace("$IPADDRESS", ipaddress);
}
```

Now when you run mbgo in the test folder the generates an mb.json that contains all the content json encoded.

you just need to use e.g. curl to POST the data with @mb.json.

The upload.sh in test shows the registration of the mock and the run of three different requests.


To build mbgo you just need to checkout my repository and run "go install"

To run the test you need to isntall mountebank on your machine and run in with "mb --allowInjection".

For the use in some secure environments I would use iptable to restrict mountebank to certain ports and also create an own mountebank user with no rights.



