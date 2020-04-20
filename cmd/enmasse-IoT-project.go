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
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	//in package import
	"github.com/IoTCLI/cmd/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/get"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
)

var (
	enmasseIoTProjectNamespaceFlag string
)

func createProject() {
	//Wait to make iot-user until the IoTproject and IoT addressspace are ready
	var iotReady = false
	var addrSpaceReady = false
	var enmasseFolderName = ""
	//Find correct enmasse path name
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			if strings.Split(f.Name(), "-")[0] == "enmasse" {
				enmasseFolderName = f.Name()
			}
		}
	}

	if enmasseFolderName == "" {
		log.Fatal("Enmasse Bundle isn't Downloaded")
	}

	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()

	//Fill in the commands that must be applied to
	//Install Enmasse Core
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/myapp-namespace.yaml")
	co.Commands = append(co.Commands, enmasseFolderName+"/install/components/iot/examples/iot-project-managed.yaml")
	co.Commands = append(co.Commands, enmasseFolderName+"/install/components/iot/examples/iot-user.yaml")
	//
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(enmasseIoTProjectNamespaceFlag)

	//Reload config flags after switching context
	newconfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionConfig := kcmdutil.NewMatchVersionFlags(newconfigFlags)
	cf := kcmdutil.NewFactory(matchVersionConfig)

	log.Println("Provision Enmasse Messaging Service")
	for commandNumber, command := range co.Commands {
		//Once IoT bundles are deployed get host IP to make certs for MQTT adapterz
		if commandNumber == 2 {

			for !iotReady && !addrSpaceReady {
				cmd := get.NewCmdGet("kubectl", cf, IOStreams)
				cmd.Flags().Set("output", "jsonpath={.items[*].status.isReady}")
				if err != nil {
					log.Fatal(err)
				}
				cmd.Run(cmd, []string{"iotproject"})
				iotReady, _ = strconv.ParseBool(out.String())
				log.Print("IoTProject is Ready: ", out.String())
				out.Reset()
				cmd.Run(cmd, []string{"addressspace"})
				addrSpaceReady, _ = strconv.ParseBool(out.String())
				log.Print("Addressspace is Ready: ", out.String())
				out.Reset()
				time.Sleep(2 * time.Second)
			}

		}
		cmd := apply.NewCmdApply("kubectl", cf, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}

}

// projectCmd represents the project command
var enmasseIoTProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("project called")
		createProject()
	},
}

func init() {
	enmasseIoTCmd.AddCommand(enmasseIoTProjectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	enmasseIoTProjectCmd.Flags().StringVarP(&enmasseIoTProjectNamespaceFlag, "namespace", "n", "myapp", "Option to specify namespace for enmasse deployment, defaults to 'myapp'")
}
