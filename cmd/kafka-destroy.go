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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/delete"
)

var (
	kafkaDestroyNamespaceFlag string
)

func kafkaDestroy() {

	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()

	_ = utils.DownloadAndUncompress("oc.gz", "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/openshift-client-linux.tar.gz")
	log.Println("oc Source folder: ", "oc")

	//Fill in the commands that must be applied to
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/kafka/setup/kafka.yaml")
	co.Commands = append(co.Commands, "https://github.com/strimzi/strimzi-kafka-operator/releases/download/0.17.0/strimzi-cluster-operator-0.17.0.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/kafka/setup/kafka-namespace.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml")
	//co.Commands = append(co.Commands, "https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/cloud-generic.yaml")
	//
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(kafkaDestroyNamespaceFlag)

	log.Println("Destroy Kafka from cluster")
	for _, command := range co.Commands {
		cmd := delete.NewCmdDelete(co.CurrentFactory, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}
	//Remove tempfile when done
}

// destroyCmd represents the destroy command
var kafkaDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Destroy called")
		kafkaDestroy()
	},
}

func init() {
	kafkaCmd.AddCommand(kafkaDestroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	kafkaDestroyCmd.Flags().StringVarP(&kafkaDestroyNamespaceFlag, "namespace", "n", "kafka", "Option to specify namespace for kafka deletion, defaults to 'kafka'")
}
