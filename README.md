# mbgo
Mountebank preprocessor to generate data for a single JSON imposter post to the mountebank server.

Mountebank http://www.mbtest.org/ is a real cool HTTP mock which supports Java Script as a mighty simulation language. Mountebank supports configuration files which do allow more human readable format, but the configuration files are limited for those who deploy Mountebank and the test scripts locally. When you want to control Mountebank remotely through some test clients or simple CURL you need to stick to some really crazy JSON style post data.

To make my life easier deploying Mountebank locally to an application server host, but running the test suites on a separate Jenkins host, I came up with a small GOLANG tool that adds some preprocessing capabilities for the Mountebank POST data.

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

The web service uses XML content and uses some java script logic to take over some tags from the request to the response and generate a random response id.

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

Now when you run mbgo in the test folder it generates an mb.json file that contains all the content in escaped json format.

It is enough to just use CURL to post the data @mb.json.

The sample test script upload.sh in the test folder shows the registration of the mock and the run of three different requests.


To build mbgo you just need to checkout my repository and run "go install"

To run the test you need to install mountebank on your machine and run in with "mb --allowInjection".

For the use in some secure environments I would use iptable to restrict mountebank to certain ports and also create an own mountebank user with no rights.
