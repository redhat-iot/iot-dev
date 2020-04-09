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
	"k8s.io/kubectl/pkg/cmd/delete"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func kafkaDestroy() {

	ocCommands := []string{}
	//Resources to delete
	ocCommands = append(ocCommands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/kafka.yaml")
	ocCommands = append(ocCommands, "https://github.com/strimzi/strimzi-kafka-operator/releases/download/0.17.0/strimzi-cluster-operator-0.17.0.yaml")
	ocCommands = append(ocCommands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/kafka-namespace.yaml")

	//Load Config for Kubectl Wrapper Function
	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionKubeConfigFlags := kcmdutil.NewMatchVersionFlags(kubeConfigFlags)

	//Create a new Credential factory
	f := kcmdutil.NewFactory(matchVersionKubeConfigFlags)

	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stdout}

	//Make a new kubctl command
	cmd := delete.NewCmdDelete(f, ioStreams)

	for _, command := range ocCommands {

		cmd.Flags().Set("filename", command)
		cmd.Flags().Set("namespace", "kafka")
		//cmd.Flags().Set("output", "json")
		//cmd.Flags().Set("dry-run", "true")
		cmd.Run(cmd, []string{})

	}

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
		fmt.Println("destroy called")
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
}
