# IoT CLI 

This CLI is intended to make deploying Cloud native IoT and Edge applications easier, More information can be found on the following [document](https://docs.google.com/document/d/1lS5YWPVCF4OhbVfB3reJDtxAojpp69W_ZBk5gwAvd6M/edit?usp=sharing)

## Prerequsites 

An functioning openshift 4.X cluster is required for this CLI's usage 

## Getting Started 

Clone the Repo, move the executable to your `$PATH`, and simply run the following command to get started  

```
IoTCLI setup --user=<Openshift admin Username> --password=<Openshift admin Password>
```
## Setting up IoT Cloud Native Messaging 

The user can chose to utilize Enmasse or Kafka for the middleware messaging layer


### Setup Enmasse 

To enable [Enmasse](enmasse.io) along with its [IoT services](https://enmasse.io/documentation/0.30.2/openshift/#'iot-guide-messaging-iot) simply run 

```
IoTCLI enmasse setup 
```

## Setting up Clound Native application tools

The user can choose which tools they want to use in order to either digest, display or process the IoT data. Eventually the user will be allowed to create custom data pipelines to connect the tools. Currently the following tools will be supported 

### Knative 

The `IoTCLI` can install and setup Knative on an Openshift cluster, deploy Knative Services, and Container Sources 

#### Currently configured Knative Services 
   
    `IoTCLI` 

- Simply Drops incoming Cloud-Events into its log 

    `iotVideo` 

- Accepts an incoming IoT video Livestream, runs image classification using Tensorflow, and serves it back to the user via a simple web appliaction, its repo can be found [here](https://github.com/astoycos/iotKnativeSource) 

#### Currently configured Knative Container Sources 
    
    `iot` 
    
- AMQP-> CloudEvent broker to ferry messages from the application side of Enmasse to a Knative Service, its repo can be found [here](https://github.com/astoycos/iotContainerSource)

### Kafka
    TODO 

### OpenDataHub 
    TODO


### Persistent Storage (Most likely Ceph)
    TODO 
