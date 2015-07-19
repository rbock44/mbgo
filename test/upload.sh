HOST=localhost

echo "deregister the imposter"
curl -i -X DELETE "http://${HOST}:2525/imposters"
mbgo
curl -i -X POST -H 'Content-Type: application/xml' -d@mb.json "http://${HOST}:2525/imposters"
echo "send GET request"
curl -i -X POST -H 'Content-Type: application/xml' -d "<request-id>1</request-id><ip-address>192.168.56.1</ip-address>" "http://${HOST}:9005/risk/verify"
curl -i -X POST -H 'Content-Type: application/xml' -d "<request-id>2</request-id><paymentMethod>"  "http://${HOST}:9005/risk/verify"
curl -i -X POST -H 'Content-Type: application/xml' -d "<request-id>2</request-id><credit-limit>"  "http://${HOST}:9005/risk/verify"
