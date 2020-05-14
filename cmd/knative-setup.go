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
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	//in package import
	"github.com/IoTCLI/cmd/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	//"k8s.io/kubectl/pkg/cmd/"
	"strings"
	"time"

	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/get"
)

//var setupStatus = false

func knativeServing() {
	//Make command options for Knative Setup
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/knative/setup/operatorgroup.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/knative/setup/sub.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/knative/setup/knative-serving.yaml")
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext("knative-serving")

	log.Info("Provision Openshift Serverless Operator and Knative Serving")
	for commandNumber, command := range co.Commands {
		cmd := apply.NewCmdApply("kubectl", co.CurrentFactory, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Info(out.String())
		out.Reset()

		//Allow time for Operator to distribute to all namespaces
		if commandNumber == 1 {
			log.Info("Pausing to let Openshift Serverless operator distribute to all namespaces")
			time.Sleep(10.0 * time.Second)
		}
	}

	var dependencies = false
	var deployments = false
	var install = false
	var ready = false
	var err error

	for !deployments || !install || !ready || !dependencies {

		cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
		//err := cmd.Flags().Set("output", "jsonpath={.status.conditions}")
		cmd.Flags().Set("template", "{{range .status.conditions}}{{printf \"%s \" .status}}{{end}}")
		cmd.Run(cmd, []string{"knativeserving.operator.knative.dev/knative-serving"})

		knativeStatus := strings.Split(out.String(), " ")
		out.Reset()

		log.Debug("The current response from Knative: ", knativeStatus)
		log.Debug("Length of response: ", len(knativeStatus))

		if len(knativeStatus) != 5 {
			log.Info("Still Starting Knative Resources")
			continue
		}
		dependencies, err = strconv.ParseBool(knativeStatus[0])
		if err != nil {
			dependencies = false
		}
		deployments, err = strconv.ParseBool(knativeStatus[1])
		if err != nil {
			deployments = false
		}
		install, err = strconv.ParseBool(knativeStatus[2])
		if err != nil {
			install = false
		}
		ready, err = strconv.ParseBool(knativeStatus[3])
		if err != nil {
			ready = false
		}

		log.Info("knative Serving Install Status:\nDependenciesInstalled=" + knativeStatus[0] + "\n" +
			"DeploymentsAvaliable=" + knativeStatus[1] + "\n" + "InstallSucceeded=" + knativeStatus[2] +
			"\n" + "Ready=" + knativeStatus[3] + "\n")

		time.Sleep(5 * time.Second)

	}

}

func knativeEventing() {

	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/knative/setup/knative-eventing.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext("knative-eventing")

	log.Info("Provision Knative Eventing")
	for _, command := range co.Commands {

		cmd := apply.NewCmdApply("kubectl", co.CurrentFactory, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Info(out.String())
		out.Reset()
	}
	time.Sleep(5 * time.Second)

	var install = false
	var ready = false
	var err error

	for !install || !ready {

		cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
		cmd.Flags().Set("template", "{{range .status.conditions}}{{printf \"%s= %s \" .type .status}}{{end}}")
		cmd.Run(cmd, []string{"knativeeventing.operator.knative.dev/knative-eventing"})

		knativeStatus := strings.Split(out.String(), " ")

		install, err = strconv.ParseBool(knativeStatus[1])
		if err != nil {
			install = false
		}
		ready, err = strconv.ParseBool(knativeStatus[3])
		if err != nil {
			ready = false
		}
		log.Info("Knative Eventing Install Status: ")
		log.Info(out.String())
		out.Reset()
		time.Sleep(5 * time.Second)

	}

}

func knativeStatus() {
	//Make command options for knative Status
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving
	co.Commands = append(co.Commands, "knativeserving.operator.knative.dev/knative-serving")
	co.Commands = append(co.Commands, "pods")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext("knative-serving")

	cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
	cmd.Flags().Set("template", "{{range .status.conditions}}{{printf \"%s \" .status}}{{end}}")
	cmd.Run(cmd, []string{co.Commands[0]})

	knativeStatus := strings.Split(out.String(), " ")
	out.Reset()

	log.Info("knative Serving Install Status:\nDependenciesInstalled=" + knativeStatus[0] + "\n" +
		"DeploymentsAvaliable=" + knativeStatus[1] + "\n" + "InstallSucceeded=" + knativeStatus[2] +
		"\n" + "Ready=" + knativeStatus[3] + "\n")

	co.SwitchContext("knative-eventing")

	cmd = get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
	cmd.Flags().Set("template", "{{range .status.conditions}}{{printf \"%s=%s \" .type .status}}{{end}}")
	cmd.Run(cmd, []string{"knativeeventing.operator.knative.dev/knative-eventing"})
	log.Info(out.String())
	out.Reset()

	cmd = get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
	cmd.Run(cmd, []string{co.Commands[1]})
	log.Info(out.String())
	out.Reset()

}

// setupCmd represents the setup command
var knativeSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Knative serving and eventing for serverless application",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if status == true {
			log.Info("Checking on knative Eventing and Serving download status")
			knativeStatus()
		} else {
			log.Info("Installing Knative Serving")
			knativeServing()
			log.Info("Installing Knative Eventing")
			knativeEventing()
			log.Info("Checking on Knative Eventing and Serving download status")
			knativeStatus()
		}

	},
}

func init() {
	knativeCmd.AddCommand(knativeSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	knativeSetupCmd.Flags().BoolVarP(&status, "status", "S", false, "Check on status of knative install")
}
