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
	"os"
	"github.com/spf13/cobra"
	//"io"
	"log"
	"os/exec"
	"strconv"
	
)

func enmasseSetup() { 

	//Test you are correctly connected to openshift cluster
	ocCommands := [][]string{}

	
	//install Enmasse
	ocCommands = append(ocCommands,[]string{"bash","-c",". ./scripts/enmasseSetup.sh"} )
	
	//ocCommands = append(ocCommands,[]string{"./oc", "get", "-n", "myapp" ,"iotproject" ,"-o" ,"jsonpath={.items[*].status.isReady}"})
	for command := range ocCommands {
		cmd := exec.Command(ocCommands[command][0], ocCommands[command][1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	//Wait to make iot-user until the IoTproject and IoT addressspace are ready 
	var iotReady=false
	var addrSpaceReady=false
	for(!iotReady && !addrSpaceReady ){
		
		iot := exec.Command("./oc", "get", "-n", "myapp" ,"iotproject" ,"-o" ,"jsonpath={.items[*].status.isReady}")
		iotin,err := iot.Output()
		if err != nil {
			log.Fatal(err)
		}
		
		addrSpace := exec.Command("./oc", "get", "-n", "myapp" ,"addressspace" ,"-o" ,"jsonpath={.items[*].status.isReady}")
		addrSpacein, err := addrSpace.Output()
		if err != nil {
			log.Fatal(err)
		}
		
		iotReady, _ = strconv.ParseBool(string(iotin))
		addrSpaceReady, _ = strconv.ParseBool(string(addrSpacein))
		
	}

	log.Println("iotProject and iotAddressspace ready creating iot-user")

	cmd := exec.Command("./oc", "create" ,"-f" ,"enmasse-0.30.2/install/components/iot/examples/iot-user.yaml")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

}



// setupCmd represents the setup command
var enmasse_setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called")
		enmasseSetup()
	},
}

func init() {
	enmasseCmd.AddCommand(enmasse_setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
