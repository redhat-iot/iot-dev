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
	"net/http"
	"bytes"
	"crypto/tls"
	"strings"

)

func device(tenant string, deviceID string){

	//Get device Registy
	registryHost , err := exec.Command("./oc" ,"-n", "enmasse-infra", "get" ,"routes", "device-registry", "--template={{ .spec.host }}").Output()
	if err != nil {
		log.Fatal("Error with getting registry host:", err)
	}

	//Get token
	token , err := exec.Command("./oc" ,"whoami", "--show-token").Output()
	if err != nil {
		log.Fatal(err)
	}

	strtoken := strings.TrimSuffix(string(token), "\n")
	//POST device
	urlDevice := "https://" + string(registryHost) + "/v1/devices/" + string(tenant) + "/" + string(deviceID)

	urlCredentials := "https://" + string(registryHost) + "/v1/credentials/" + string(tenant) + "/" + string(deviceID)
	
	credentialsJSON := []byte(`[{
			"type": "hashed-password",
			"auth-id": "sensor1",
			"secrets": [{
				"pwd-plain":"hono-secret"
			}]
		}]`)
	
	
	log.Println("Device Url:",urlDevice)
	log.Println("Credentials Url:",urlCredentials)
	log.Println("Payload:", string(bytes.TrimSpace(credentialsJSON)))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	addDevice, err := http.NewRequest("POST", urlDevice,nil)
	addCredentials, err := http.NewRequest("PUT", urlCredentials, bytes.NewBuffer(bytes.TrimSpace(credentialsJSON)))
	if err != nil {
		// handle err
	}

	addDevice.Header.Set("Content-Type", "application/json")
	addDevice.Header.Set("Authorization", "Bearer " + strtoken)
	addCredentials.Header.Set("Content-Type", "application/json")
	addCredentials.Header.Set("Authorization", "Bearer " + strtoken)

	devResp, err := client.Do(addDevice)
	if err != nil {
		// handle err
		log.Fatal("Http POST error: ", err)
	}
	log.Println("Device Post Response:", devResp)
	defer devResp.Body.Close()

	creResp, err := client.Do(addCredentials)
	if err != nil {
		log.Fatal("Http POST error: ", err)
	}

	log.Println("Credential Post Response:", creResp)
	defer creResp.Body.Close()
}

// addDeviceCmd represents the addDevice command
var enmasse_addDeviceCmd = &cobra.Command{
	Use:   "addDevice",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enmasse addDevice called")
		device(args[0] , args[1])
	},
}

func init() {
	enmasseCmd.AddCommand(enmasse_addDeviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addDeviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addDeviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
