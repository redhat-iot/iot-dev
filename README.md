# Setting up Clound Native Edge systems with the IoTCLI

The user can choose which tools they want to use in order to either digest, display or process the IoT data. Eventually the user will be allowed to create custom data pipelines to connect the tools. Currently the following tools/commands are supported 

## Getting Started 

Download the latest Release from the [release page](https://github.com/redhat-iot/iot-dev/releases)

Move the Binary `IoTCLI` into your `$PATH`

Run `IoTCLI login` to authenticate with your Openshift cluster 

# Main Components Overview

## Messaging Core 

In any Edge system that incorporates IoT devices, having a scalable cloud native edge system is very important.  With the `IoTCLI` you can easily deploy these messaging services to your Openshift 4.X cluster
    
### Enmasse 

[Enmasse](enmasse.io) is an open source project for managed, self-service messaging on Kubernetes. 

The IoTCLI makes it easy to deploy the system along with its [IoT services](https://enmasse.io/documentation/0.31.0/openshift/#iot-guide-messaging-iot) onto an openshift Cluster 

### Kafka 

[Apache Kafka](https://kafka.apache.org/intro) is an advanced distributed streaming platform. It is deployed to your Openshift cluster using the [Strimzi operator](https://strimzi.io/)

## Application Core

Once IoT devices are streaming or pinging data from the edge developers need a way to process the data either on the Edge or back in the central cloud cloud 

### In-Cloud

#### Knative 

The `IoTCLI` can install and setup Knative on an Openshift cluster, deploy Knative Services, and Container Sources 

##### Currently configured Knative Services 
   
`event-display`

- Simply Drops incoming Cloud-Events into its log 

`iotVideo` 

- Accepts an incoming IoT video Livestream as HLS Playlist File(`.m3u8`), runs image classification using Tensorflow, and serves it back to the user via a simple web appliaction, its repo can be found [here](https://github.com/astoycos/iotKnativeSource) 

`video-analytics` 

- Accepts a video livestream in form of an HLS Playlist file delivered via a JSON structured to work with kafka. Sends video to Tensorflow serving module for inference, and sends inferenced video segments ceph object storage for eventual serving.  Part of [2020 Summit Demo](https://github.com/redhat-iot/2020Summit-IoT-Streaming-Demo)

`video-serving`

- Serves analyzed live video from ceph to user via a simple flask web app. Part of [2020 Summit Demo](https://github.com/redhat-iot/2020Summit-IoT-Streaming-Demo)

#### Currently configured Knative Sources 
    
 `iot` 
    
- AMQP-> CloudEvent broker to ferry messages from the application side of Enmasse to a Knative Service, its repo can be found [here](https://github.com/astoycos/iotContainerSource)

`kafka` 

- KafkaJSON -> Cloud Event broker to ferry video playlist files from Kafka to Knative 

#### Persistent Storage 

Persistent storage is implemented using [Ceph Object Storage](https://ceph.io/ceph-storage/object-storage/) deployed via the [Rook Operator](https://rook.io/)

### On Edge 

#### Gstreamer
    Gstreamer AI plugin to run analytics on low resource devices @Aditya

#### TensorFlow lite 

    TODO

    
## Command Reference 

| Command1 (Tool)   | Command2(Module) | Command3(Module command) | Argument (module input)             | Flags                          | Function                                                                                                                                           |
|-------------------|------------------|--------------------------|-------------------------------------|--------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|
| ceph              | destroy          | None                     | None                                | None                           | Destroy the ceph object store instance on the cluster                                                                                              |
| ceph              | user             | None                     | `<User Name>`                       | None                           | Create an Ceph Object Store User                                                                                                                   |
| ceph              | secrets          | None                     | `<User Name>`                       | None                           | Return user secrets for ceph deployment                                                                                                            |
| ceph              | setup            | None                     | None                                | None                           | Setup Ceph Object Storage via the Rook Operator                                                                                                    |
| enmasse           | destroy          | None                     | None                                | None                           | Remove Enmasse from openshift  cluster                                                                                                             |
| enmasse           | IoT              | addDevice                | `<Messaging Tenant>`, `<DeviceID>`  | None                           | Add a Device with specified ID to the Enmasse device registry for a specified messaging TenantSetup default Credentials                            |
| enmasse           | IoT              | project                  | None                                | --namespace                    | Make a new enmasse IoT project in the specified namespace, defaults to “myapp”                                                                     |
| enmasse           | setup            | None                     | None                                | None                           | Download Enmasse Source, store in current directory. Setup Enmasse Setup IoT services                                                              |
| kafka             | bridge           | None                     | None                                | --namespace                    | Deploy kafka HTTP bridge Deploy nginx ingress to access bridge from outside the cluster (will be transitioned to a route)                          |
| kafka             | destroy          | None                     | None                                | --namespace                    | Destroy the kafka deployment located at the specified namespace                                                                                    |
| kafka             | setup            | None                     | None                                | --namespace                    | Setup a kafka cluster viat the strimzi operator at the specified namespace                                                                         |
| knative           | setup            | None                     | None                                | --status=true/false            | Setup Knative serverless on openshift clusterConfigures both Knative-Eventing and  Knative-ServingSet --status=true to check on Knative deployment |
| knative           | destroy          | None                     | None                                | None                           | Remove Knative deployment from openshift cluster                                                                                                   |
| knative           | service          | None                     | `<Knative service to be deployed>`  | --status=true/false--namespace | Deploy a knative service Set --status=true to check on Knative service deploymentDefault namespace is “knative-eventing”                           |
| knative           | service          | destroy                  | `<Knative service to be destroyed>` | --namespace                    | Remove a specified Knative service from the cluster at specified namespace Default namespace is “knative-eventing”                                 |
| knative           | source           | None                     | `<containersource to be deployed>`  | --namespace                    | Deploy a Knative Source at specified namespace Defaults to namespace “knative-eventing”                                                            |
| knative           | source           | destroy                  | `<containersource to be destroyed>` | --namespace                    | Remove a specified knative source from the cluster from specified namespaceDefault namespace is “knative-eventing”                                 |
| login             | None             | None                     | None                                | None                           | Login to your openshift cluster with the username and password                                                                                     |
| tensorflowServing | setup            | None                     | None                                | --namespace                    | Deploy a tensorflow serving deployment and service Default namespace is “default”                                                                  |
| tensorflowServing | destroy          | None                     | None                                | --namesapce                    | Deploy a tensorflow serving deployment and serviceDefault namespace is “default”                                                                   |