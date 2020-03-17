#!/bin/bash

./oc new-project enmasse-infra

./oc apply -n enmasse-infra -f enmasse-0.30.2/install/bundles/enmasse

./oc apply -n enmasse-infra -f enmasse-0.30.2/install/components/example-plans

./oc apply -n enmasse-infra -f enmasse-0.30.2/install/components/example-roles

./oc apply -n enmasse-infra -f enmasse-0.30.2/install/components/example-authservices/standard-authservice.yaml

./oc apply -n enmasse-infra -f enmasse-0.30.2/install/preview-bundles/iot

export CLUSTERIP=$(oc get svc -n openshift-ingress -o jsonpath='{.items[*].spec.clusterIP}')

CLUSTER=$CLUSTERIP.nip.io enmasse-0.30.2/install/components/iot/examples/k8s-tls/create

./oc create -n enmasse-infra secret tls iot-mqtt-adapter-tls --key=enmasse-0.30.2/install/components/iot/examples/k8s-tls/build/iot-mqtt-adapter-key.pem --cert=enmasse-0.30.2/install/components/iot/examples/k8s-tls/build/iot-mqtt-adapter-fullchain.pem

./oc apply -n enmasse-infra -f enmasse-0.30.2/install/components/iot/examples/infinispan/common

./oc apply -n enmasse-infra -f enmasse-0.30.2/install/components/iot/examples/infinispan/manual

./oc apply -n enmasse-infra -f enmasse-0.30.2/install/components/iot/examples/iot-config.yaml

./oc new-project myapp

./oc create -n my-app -f enmasse-0.30.2/install/components/iot/examples/iot-project-managed.yaml
