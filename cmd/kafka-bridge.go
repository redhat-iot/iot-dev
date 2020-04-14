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
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func kafkaBridge() {

	ocCommands := []string{}

	//Setup kafka bridge
	ocCommands = append(ocCommands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/kafka-bridge.yaml")
	//Setup Nginix Ingress **CONVERT TO OPENSHIFT ROUTE AT SOME POINT** to connect to bridge from outside the cluster
	//Get Nginix controller and apply to cluster
	ocCommands = append(ocCommands, "https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml")
	ocCommands = append(ocCommands, "https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/cloud-generic.yaml")
	//Seutp the K8s ingress resource

	ocCommands = append(ocCommands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ingress.yaml")

	//Load Config for Kubectl Wrapper Function
	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionKubeConfigFlags := kcmdutil.NewMatchVersionFlags(kubeConfigFlags)

	//Create a new Credential factory
	f := kcmdutil.NewFactory(matchVersionKubeConfigFlags)

	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stdout}

	//Make a new kubctl command
	//cmd := apply.NewCmdApply("kubectl", f, ioStreams)

	//Setup the bridge
	for _, command := range ocCommands {
		cmd := apply.NewCmdApply("kubectl", f, ioStreams)
		cmd.Flags().Set("filename", command)
		cmd.Flags().Set("namespace", "kafka")
		//cmd.Flags().Set("output", "json")
		//cmd.Flags().Set("dry-run", "true")
		cmd.Run(cmd, []string{})

	}

	//wait until pods are up

}

// bridgeCmd represents the bridge command
var kafkaBridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bridge called")
		kafkaBridge()
	},
}

func init() {
	kafkaCmd.AddCommand(kafkaBridgeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bridgeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bridgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
