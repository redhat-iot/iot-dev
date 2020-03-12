#!/bin/bash

./oc new-project enmasse-infra

./oc apply -f enmasse-0.30.2/install/bundles/enmasse

./oc apply -f enmasse-0.30.2/addressspaces.crd.yaml

./oc apply -f enmasse-0.30.2/install/components/example-plans

./oc apply -f enmasse-0.30.2/install/components/example-roles

./oc apply -f enmasse-0.30.2/install/components/example-authservices/standard-authservice.yaml

./oc apply -f enmasse-0.30.2/install/preview-bundles/iot

export CLUSTERIP = $(oc get svc -n openshift-ingress -o jsonpath='{.items[*].spec.clusterIP'})

CLUSTER=$CLUSTERIP.nip.io install/components/iot/examples/k8s-tls/create

./oc create secret tls iot-mqtt-adapter-tls --key=install/components/iot/examples/k8s-tls/build/iot-mqtt-adapter-key.pem --cert=install/components/iot/examples/k8s-tls/build/iot-mqtt-adapter-fullchain.pem

./oc apply -f install/components/iot/examples/infinispan/common

./oc apply -f install/components/iot/examples/infinispan/manual

./oc apply -f install/components/iot/examples/iot-config.yaml
