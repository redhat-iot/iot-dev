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
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/IoTCLI/cmd/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
)

var (
	enmasseSetupNamespaceFlag string
)

func enmasseSetup() {

	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()
	//download Enmasse v0.30.3
	folderName := utils.DownloadAndUncompress("enmasse.tgz", "https://github.com/EnMasseProject/enmasse/releases/download/0.30.3/enmasse-0.30.3.tgz")
	log.Println("Enmasse Source folder: ", folderName)

	//If you want to deploy in a namespace besides enmasse-infra
	//TODO still Test This
	if enmasseSetupNamespaceFlag != "enmasse-infra" {

		utils.LocalSed((`s/enmasse-infra/` + enmasseSetupNamespaceFlag + `/`), (folderName + "/install/bundles/enmasse/"))
		utils.LocalSed((`s/enmasse-infra/` + enmasseSetupNamespaceFlag + `/`), (folderName + "/install/preview-bundles/iot/"))

	}
	//Fill ain the commands that must be applied to
	//Install Enmasse Core
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/enmasse-infra-namespace.yaml")
	co.Commands = append(co.Commands, folderName+"install/bundles/enmasse")
	co.Commands = append(co.Commands, folderName+"install/components/example-plans")
	co.Commands = append(co.Commands, folderName+"install/components/example-roles")
	co.Commands = append(co.Commands, folderName+"install/components/example-authservices/standard-authservice.yaml")
	co.Commands = append(co.Commands, folderName+"install/components/service-broker")
	//co.Commands = append(co.Commands, folderName+"/install/components/cluster-service-broker")
	//Install Enmasse IoT
	co.Commands = append(co.Commands, folderName+"install/preview-bundles/iot")

	co.Commands = append(co.Commands, folderName+"install/components/iot/examples/infinispan/common")
	co.Commands = append(co.Commands, folderName+"install/components/iot/examples/infinispan/manual")
	co.Commands = append(co.Commands, folderName+"install/components/iot/examples/iot-config.yaml")
	//
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(enmasseSetupNamespaceFlag)

	//Reload config flags after switching context
	log.Println("Provision Enmasse Messaging Service")
	for commandNumber, command := range co.Commands {
		//Once IoT bundles are deployed get host IP to make certs for MQTT adapter
		if commandNumber == 8 {
			//Add ability to put custom IP here
			err := exec.Command("./" + folderName + "/install/components/iot/examples/k8s-tls/create").Run()
			if err != nil {
				log.Fatal(err)
			}
		}
		cmd := apply.NewCmdApply("kubectl", co.CurrentFactory, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}
	os.RemoveAll(folderName)

}

// setupCmd represents the setup command
var enmasseSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Enmasse as a messaging backend",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("setup called")
		enmasseSetup()
	},
}

func init() {
	enmasseCmd.AddCommand(enmasseSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	enmasseSetupCmd.Flags().StringVarP(&enmasseSetupNamespaceFlag, "namespace", "n", "enmasse-infra", "Option to specify namespace for enmasse deployment, defaults to 'enmasse-infra'")
}
