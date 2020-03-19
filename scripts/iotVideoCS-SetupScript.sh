#!/bin/bash
export IP=$(oc get pods --selector='serving.knative.dev/service' -o jsonpath='{.items[*].status.podIP'})
echo $IP

export MESSAGE_URI=$(oc -n myapp get addressspace iot -o jsonpath={.status.endpointStatuses[?\(@.name==\'messaging\'\)].externalHost})
echo $MESSAGE_URI

export MESSAGE_PORT=443
echo $MESSAGE_PORT

export MESSAGE_TYPE=telemetry
echo $MESSAGE_TYPE
 
export MESSAGE_TENANT=myapp.iot
echo $MESSAGE_TENANT

export TLS_CONFIG="1"
echo $TLS_CONFIG

export TLS_CERT=$(oc -n myapp get addressspace iot -o jsonpath={.status.caCert} | base64 --decode)
echo $TLS_CERT

export CLIENT_USERNAME=consumer
echo $CLIENT_USERNAME

export CLIENT_PASSWORD=foobar
echo $CLIENT_PASSWORD