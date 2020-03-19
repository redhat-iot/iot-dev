#!/bin/bash
./oc delete -n myapp iotproject iot

./oc delete -f enmasse-0.30.2/install/components/iot/examples/iot-config.yaml

./oc get iotproject -o yaml

./oc delete -n enmasse-infra -f enmasse-0.30.2/install/components/iot/examples/infinispan/common
./oc delete -n enmasse-infra -f enmasse-0.30.2/install/components/iot/examples/infinispan/manual

./oc delete -n enmasse-infra -f enmasse-0.30.2/install/preview-bundles/iot

./oc delete -n enmasse-infra -f enmasse-0.30.2/install/components/example-authservices/standard-authservice.yaml

./oc delete -n enmasse-infra -f enmasse-0.30.2/install/components/example-roles

./oc delete -n enmasse-infra -f enmasse-0.30.2/install/components/example-plans

./oc delete -n enmasse-infra -f enmasse-0.30.2/install/bundles/enmasse

./oc delete project myapp
./oc delete project enmasse-infra
