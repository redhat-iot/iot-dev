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
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	//in package import
	"github.com/IoTCLI/cmd/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	//"k8s.io/kubectl/pkg/cmd/"
	"k8s.io/kubectl/pkg/cmd/delete"
)

///STILL NOT UPDATED TO Wrappers
func destroyKnative() {

	//Make command options for Knative Setup
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/knative/setup/knative-eventing.yaml")

	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/knative/setup/knative-serving.yaml")
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext("knative-eventing")

	log.Info("Remove Openshift Serverless Operator and Knative Serving")
	for commandNumber, command := range co.Commands {
		if commandNumber == 1 {
			co.SwitchContext("knative-serving")
		}
		cmd := delete.NewCmdDelete(co.CurrentFactory, IOStreams)
		cmd.Flags().Set("filename", command)
		cmd.Run(cmd, []string{})
		log.Info(out.String())
		out.Reset()
		//Allow time for Operator to distribute to all namespaces
	}

	/*

		currentCSV, err := exec.Command("bash", "-c", "./oc get subscription serverless-operator -n openshift-operators -o jsonpath='{.status.currentCSV}'").Output()
		err = exec.Command("./oc", "delete", "subscription", "serverless-operator", "-n", "openshift-operators").Run()
		if err != nil {
			//Ignore Resource Not found error and continue, but still notify the user
			log.Println(err)
		}
		err = exec.Command("./oc", "delete", "clusterserviceversion", string(currentCSV), "-n", "openshift-operators").Run()
		if err != nil {
			//Igonore Resource Not found error and continue, but still notify the user
			log.Println(err)
		}
		os.Remove("oc")
	*/

}

// destroyCmd represents the destroy command
var knativeDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Knative destroy called")
		destroyKnative()
	},
}

func init() {
	knativeCmd.AddCommand(knativeDestroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
