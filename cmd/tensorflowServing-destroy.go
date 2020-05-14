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

func tensorflowServingDestroy() {

	co := utils.NewCommandOptions()

	co.Commands = append(co.Commands, "tensorflow-deployment")
	co.Commands = append(co.Commands, "coco-service")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(tensorflowServingNamespaceFlag)

	//Reload config flags after switching context

	log.Println("Provision Tensorflow Serving Pod")
	for commandNumber, command := range co.Commands {
		cmd := delete.NewCmdDelete(co.CurrentFactory, IOStreams)
		if commandNumber == 0 {
			cmd.Run(cmd, []string{"deployment", command})
		} else {
			cmd.Run(cmd, []string{"service", command})
		}
		log.Print(out.String())
		out.Reset()
	}
}

// tensorflowServingDestroyCmd represents the tensorflowServingDestroy command
var tensorflowServingDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the deployed tensorflow service",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("tensorflowServingDestroy called")
		tensorflowServingDestroy()
	},
}

func init() {
	tensorflowServingCmd.AddCommand(tensorflowServingDestroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tensorflowServingDestroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tensorflowServingDestroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
