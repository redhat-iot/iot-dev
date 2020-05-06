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
	"time"

	"github.com/IoTCLI/cmd/utils"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/get"
)

//Made from Instructions @https://opendatahub.io/docs/administration/advanced-installation/object-storage.html for installing
//ceph object storage via the rook operator
func cephSetup() {
	//Make command options for Knative Setup
	co := utils.NewCommandOptions()

	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/scc.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/operator.yaml")
	co.Commands = append(co.Commands, "pods")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/cluster.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/toolbox.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/object.yaml")
	co.Commands = append(co.Commands, "pods")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/object-user.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/route.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	//Switch Context and Reload Config Flags
	co.SwitchContext("rook-ceph-system")

	log.Println("Provision Knative Source")
	for commandNumber, command := range co.Commands {
		//After the system pods are provisioned wait for them to become ready before moving on
		if commandNumber == 2 {
			log.Print("Waiting for Pods to be ready in rook-ceph-system namespace:")
			podStatus := utils.NewpodStatus()
			for podStatus.Running != 7 {
				cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
				cmd.Flags().Set("output", "yaml")
				cmd.Run(cmd, []string{command})
				podStatus.CountPods(out.Bytes())
				log.Print("Waiting...")
				out.Reset()
				time.Sleep(5 * time.Second)
			}
			co.SwitchContext("rook-ceph")

		} else if commandNumber == 6 {
			//After the pods in rook-ceph are provisioned wait for them to become ready before moving on
			log.Print("Waiting for pods to be ready in rook-ceph")
			podStatus := utils.NewpodStatus()
			for podStatus.Running != 9 && podStatus.Succeeded != 3 {
				cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
				cmd.Flags().Set("output", "yaml")
				cmd.Run(cmd, []string{command})
				podStatus.CountPods(out.Bytes())
				log.Print("Waiting...")
				out.Reset()
				time.Sleep(5 * time.Second)
			}
			time.Sleep(5 * time.Second)
		} else {
			cmd := apply.NewCmdApply("kubectl", co.CurrentFactory, IOStreams)
			//Kubectl signals missing field, set validate to false to ignore this
			cmd.Flags().Set("validate", "false")
			err := cmd.Flags().Set("filename", command)
			if err != nil {
				log.Fatal(err)
			}
			cmd.Run(cmd, []string{})
			log.Print(out.String())
			out.Reset()
		}
	}

}

// setupCmd represents the setup command
var cephSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Ceph Object storage",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Ceph setup called")
		cephSetup()
	},
}

func init() {
	cephCmd.AddCommand(cephSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
