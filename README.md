# IoT CLI 

This CLI is intended to make deploying Cloud native IoT and Edge applications easier, More information can be found on the following [document](https://docs.google.com/document/d/1lS5YWPVCF4OhbVfB3reJDtxAojpp69W_ZBk5gwAvd6M/edit?usp=sharing)

## Prerequsites 

An functioning openshift 4.X cluster is required for this CLI's usage 

## Getting Started 

[Download a release](https://github.com/redhat-iot/iot-dev/releases), move the executable to your `$PATH`, and simply run the following command to get started  

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
   
`IoTCLI` Service

- Simply Drops incoming Cloud-Events into its log 

`iotVideo` Service 

- Accepts an incoming IoT video Livestream, runs image classification using Tensorflow, and serves it back to the user via a simple web appliaction, its repo can be found [here](https://github.com/astoycos/iotKnativeSource) 

#### Currently configured Knative Container Sources 
    
 `iot` ContainerSource
    
- AMQP-> CloudEvent broker to ferry messages from the application side of Enmasse to a Knative Service, its repo can be found [here](https://github.com/astoycos/iotContainerSource)

### Kafka
    TODO 

### OpenDataHub 
    TODO


### Persistent Storage (Most likely Ceph)
    TODO 
    
## Command Reference 

|          |                  |                          |                                   |                     |                                                                                                                                              |
|----------|------------------|--------------------------|-----------------------------------|---------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| **Command1** | **Command2(Module)** | **Command3(Module command)** | **Argument (module input)**           | **Flags**               | **Function**                                                                                                                                     |
| setup    | None             | None                     | None                              | None                | Download required files and binaries for all available tools                                                                                 |
| enmasse  | setup            | None                     | None                              | None                | Setup Enmasse Setup IoT services                                                                                                             |
| enmasse  | destroy          | None                     | None                              | None                | Remove Enmasse from openshift cluster                                                                                                        |
| enmasse  | addDevice        | None                     | <Messaging Tenant> <DeviceID>     | None                | Add a Device with specified ID to the Enmasse device registry for a specified messaging TenantSetup default Credentials                      |
| knative  | setup            | None                     | None                              | --status=true/false | Setup Knative serverless on openshift clusterConfigures Knative-Eventing and Knative-ServingSet --status=true to check on Knative deployment |
| knative  | destroy          | None                     | None                              | None                | Remove Knative deployment from openshift cluster                                                                                             |
| knative  | service          | None                     | <Knative service to be deployed>  | --status=true/false | Deploy a knative service Set --status=true to check on Knative service deployment                                                            |
| knative  | service          | destroy                  | <Knative service to be destroyed> | None                | Remove a specified Knative service from the cluster                                                                                          |
| knative  | cs               |                          | <containersource to be deployed>  | None                | Deploy a Knative ContainerSource                                                                                                             |
| Knative  | cs               | destroy                  | <containersource to be destroyed> | None                | Remove a specified containersource from the cluster                                                                                          |
| kafka    | setup            |                          |                                   |                     |                                                                                                                                              |
| Kafka    | destroy          |                          |                                   |                     |                                                                                                                                              |
