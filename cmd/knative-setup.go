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
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	//in package import
	"github.com/IoTCLI/cmd/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	//"k8s.io/kubectl/pkg/cmd/"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/get"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"

	"time"
)

//var setupStatus = false

func knativeServing() {
	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/operatorgroup.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/sub.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/knative-serving.yaml")
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext("knative-serving")

	//Reload config flags after switching context
	newconfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionConfig := kcmdutil.NewMatchVersionFlags(newconfigFlags)
	cf := kcmdutil.NewFactory(matchVersionConfig)

	log.Println("Provision Openshift Serverless Operator and Knative Serving")
	for commandNumber, command := range co.Commands {
		cmd := apply.NewCmdApply("kubectl", cf, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
		//Allow time for Operator to distribute to all namespaces
		if commandNumber == 1 {
			time.Sleep(2.0 * time.Second)
		}
	}

	var dependencies = false
	var deployments = false
	var install = false
	var ready = false
	//m = make(map[string]int)

	for !deployments && !install && !ready && !dependencies {

		cmd := get.NewCmdGet("kubectl", cf, IOStreams)
		err := cmd.Flags().Set("output", "jsonpath={.status.conditions}")
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{"knativeserving.operator.knative.dev/knative-serving"})
		knativeStatus := out.String()
		knativeStatusSlice := strings.Split(knativeStatus, " ")

		dependencies, _ = strconv.ParseBool(knativeStatusSlice[1][7:11])
		deployments, _ = strconv.ParseBool(knativeStatusSlice[4][7:11])
		install, _ = strconv.ParseBool(knativeStatusSlice[7][7:11])
		ready, _ = strconv.ParseBool(knativeStatusSlice[10][7:11])

		log.Print("knative Serving Install Status:\nDependenciesInstalled=" + knativeStatusSlice[1][7:11] + "\n" +
			"DeploymentsAvaliable=" + knativeStatusSlice[4][7:11] + "\n" + "InstallSucceeded=" + knativeStatusSlice[7][7:11] +
			"\n" + "Ready=" + knativeStatusSlice[10][7:11] + "\n")
		out.Reset()

		time.Sleep(5 * time.Second)

	}

}

//kubectl delete --selector knative.dev/crd-install=true --filename https://github.com/knative/eventing/releases/download/v0.13.0/eventing.yaml

func knativeEventing() {

	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving
	co.Commands = append(co.Commands, "https://github.com/knative/eventing/releases/download/v0.13.0/eventing.yaml")
	co.Commands = append(co.Commands, "https://github.com/knative/eventing/releases/download/v0.13.0/eventing.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext("knative-eventing")

	//Reload config flags after switching context
	newconfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionConfig := kcmdutil.NewMatchVersionFlags(newconfigFlags)
	cf := kcmdutil.NewFactory(matchVersionConfig)

	log.Println("Provision Knative Eventing")
	for commandNumber, command := range co.Commands {

		cmd := apply.NewCmdApply("kubectl", cf, IOStreams)
		if commandNumber == 0 {
			cmd.Flags().Set("selector", "knative.dev/crd-install")
		}
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}
	time.Sleep(5 * time.Second)
}

func knativeStatus() {
	//Make command options for knative Status
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving
	co.Commands = append(co.Commands, "knativeserving.operator.knative.dev/knative-serving")
	co.Commands = append(co.Commands, "pods")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext("knative-serving")

	//Reload config flags after switching context
	newconfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionConfig := kcmdutil.NewMatchVersionFlags(newconfigFlags)
	cf := kcmdutil.NewFactory(matchVersionConfig)

	cmd := get.NewCmdGet("kubectl", cf, IOStreams)
	cmd.Flags().Set("template", "'{{range .status.conditions}}{{printf \"%s=%s\" .type .status}}{{end}}'")
	cmd.Run(cmd, []string{co.Commands[0]})
	log.Print(out.String())
	out.Reset()

	co.SwitchContext("knative-eventing")

	//Reload config flags after switching context
	newconfigFlags = genericclioptions.NewConfigFlags(true)
	matchVersionConfig = kcmdutil.NewMatchVersionFlags(newconfigFlags)
	cf = kcmdutil.NewFactory(matchVersionConfig)

	cmd = get.NewCmdGet("kubectl", cf, IOStreams)
	cmd.Run(cmd, []string{co.Commands[1]})
	log.Print(out.String())
	out.Reset()

}

// setupCmd represents the setup command
var knativeSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if status {
			log.Println("Checking on knative Eventing and Serving download status")
			knativeStatus()
		} else {
			log.Println("Installing Knative Serving")
			knativeServing()
			log.Println("Installing Knative Eventing")
			knativeEventing()
			log.Println("Checking on Knative Eventing and Serving download status")
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
