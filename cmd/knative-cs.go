/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"github.com/spf13/cobra"
	"os/exec"
	"os"

)

func iotContainerSource(containerSource string) {
	ocCommands := [][]string{}
	
	messageURI , err := exec.Command("./oc", "-n", "myapp", "get", "addressspace", "iot", "-o", "jsonpath={.status.endpointStatuses[?(@.name=='messaging')].externalHost}").Output()
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("MESSAGE_URI", string(messageURI))
	log.Println(string(messageURI))

	os.Setenv("MESSAGE_PORT","443")

	os.Setenv("MESSAGE_TYPE","telemetry")

	os.Setenv("MESSAGE_TENANT","myapp.iot")

	tlsCert, err := exec.Command("bash", "-c", "oc -n myapp get addressspace iot -o jsonpath={.status.caCert} | base64 --decode").Output()
	if err != nil {
		log.Fatal(err)
	}	

	os.Setenv("TLS_CERT",string(tlsCert))

	os.Setenv("CLIENT_USERNAME","consumer")

	os.Setenv("CLIENT_PASSWORD","foobar")
	
	//ocCommands = append(ocCommands,[]string{"/bin/bash", "-c", ". ./scripts/iotVideoCS-SetupScript.sh"} )
	ocCommands = append(ocCommands,[]string{"/bin/bash","-c","cat yamls/" + containerSource + "ContainerSource.yaml.in | envsubst | oc apply -n knative-eventing -f -"} )
	
	for command := range ocCommands {
		cmd := exec.Command(ocCommands[command][0], ocCommands[command][1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	
}

// csCmd represents the cs command
var csCmd = &cobra.Command{
	Use:   "cs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cs called")
		
		iotContainerSource(args[0])
		
	},
}

func init() {
	knativeCmd.AddCommand(csCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// csCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// csCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
