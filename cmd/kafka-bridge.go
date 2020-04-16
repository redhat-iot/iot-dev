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
	"github.com/IoTCLI/cmd/utils"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
	"log"
)

var (
	kafkaBridgeNamespaceFlag string
)

func kafkaBridge() {

	co := utils.NewCommandOptions()

	//Setup kafka bridge
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/kafka-bridge.yaml")
	//Setup Nginix Ingress **CONVERT TO OPENSHIFT ROUTE AT SOME POINT** to connect to bridge from outside the cluster
	//Get Nginix controller and apply to cluster
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/cloud-generic.yaml")
	//Seutp the K8s ingress resource

	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ingress.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(kafkaBridgeNamespaceFlag)

	//Reload config flags after switching context
	newconfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionConfig := kcmdutil.NewMatchVersionFlags(newconfigFlags)
	cf := kcmdutil.NewFactory(matchVersionConfig)

	log.Println("Provision Kafka Http Bridge")
	for _, command := range co.Commands {
		cmd := apply.NewCmdApply("kubectl", cf, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}
	log.Println("To check status of Kafka HTTP bridge run 'curl -v GET http://my-bridge.io/healthy'")
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
		log.Println("Kafka Http Bridge called")
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
	kafkaBridgeCmd.Flags().StringVarP(&kafkaBridgeNamespaceFlag, "namespace", "n", "kafka", "Option to specify namespace for kafka deletion, defaults to 'kafka'")
}
