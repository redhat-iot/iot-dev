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
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/IoTCLI/cmd/utils"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
)

// setupCmd represents the setup command
var gstreamerSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Gstreamer with Openvino toolkit for video streaming and analytics",
	Long: `A longer descriptSion that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gstreamer setup called")
		fstatus, _ := cmd.Flags().GetBool("local")
		if fstatus { // if status is true,
			gstreamerLocalSetup()
		} else {
			gstreamerSetup()
		}
	},
}

func gstreamerLocalSetup() {

	cmd := exec.Command("https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/gstreamer/startup.sh")
	out, error := cmd.Output()
	if error != nil {
		println(error.Error())
		return
	} else {
		log.Println(string(out))
	}

	for {
		fmt.Print("Press 'S' to stop the Gstreamer container: ")
		var key string
		fmt.Scanln(&key)

		if key == "s" {
			cmd2 := exec.Command("/bin/sh", "-c", " docker kill gstreamer_container")
			out2, error2 := cmd2.Output()
			if error2 != nil {
				println(error2.Error())
				return
			} else {
				log.Println(string(out2))
			}
			fmt.Print("Container stopped: ")

			break

		}
	}

}

func gstreamerSetup() {
	//Make command options for Kafka Setup
	co := utils.NewCommandOptions()
	//This section is mimicking the instructions to setup the Strimzi Operator, I.E download the install yaml, and set namespace using sed
	//functionality

	//Fill in the commands that must be applied to
	co.Commands = append(co.Commands, "/home/adkadam/work/golang/iot-dev/yamls/gstreamer/gstreamer-namespace.yaml")
	co.Commands = append(co.Commands, "/home/adkadam/work/golang/iot-dev/yamls/gstreamer/gstreamer-deploy.yaml")
	//
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext("gstreamer")

	log.Println("Provision gstreamer")
	for _, command := range co.Commands {
		cmd := apply.NewCmdApply("kubectl", co.CurrentFactory, IOStreams)
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}
}

func init() {
	gstreamerCmd.AddCommand(gstreamerSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	gstreamerSetupCmd.Flags().BoolP("local", "l", false, "Setup gstreamer locally")
}
