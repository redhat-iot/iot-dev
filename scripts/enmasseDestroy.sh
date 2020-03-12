#!/bin/bash
./oc delete -f enmasse-0.30.2/install/components/iot/examples/iot-config.yaml

./oc delete -f enmasse-0.30.2/install/components/iot/examples/infinispan/common
./oc delete -f enmasse-0.30.2/install/components/iot/examples/infinispan/manual

./oc delete -f enmasse-0.30.2/install/preview-bundles/iot

./oc delete -f enmasse-0.30.2/install/components/example-authservices/standard-authservice.yaml

./oc delete -f enmasse-0.30.2/install/components/example-roles

./oc delete -f enmasse-0.30.2/install/components/example-plans

./oc delete -f enmasse-0.30.2/install/bundles/enmasse

./oc delete project enmasse-infra