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

	"github.com/IoTCLI/cmd/utils"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/delete"
)

func gstreamerDestroy() {

	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()

	_ = utils.DownloadAndUncompress("oc.gz", "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/openshift-client-linux.tar.gz")
	log.Println("oc Source folder: ", "oc")

	//Fill in the commands that must be applied to
	co.Commands = append(co.Commands, "/home/adkadam/work/golang/iot-dev/yamls/gstreamer/gstreamer-deploy.yaml")
	co.Commands = append(co.Commands, "/home/adkadam/work/golang/iot-dev/yamls/gstreamer/gstreamer-namespace.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext("gstreamer")

	log.Println("Destroy Gstreamer deployment from cluster")
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
var gstreamerDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy gstreamer deployment",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Gstreamer destroy called")
		gstreamerDestroy()
	},
}

func init() {
	gstreamerCmd.AddCommand(gstreamerDestroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
