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
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/get"
	"time"
)

//Made from Instructions @https://opendatahub.io/docs/administration/advanced-installation/object-storage.html for installing
//ceph object storage via the rook operator
func cephSetup() {
	//Make command options for Knative Setup
	co := utils.NewCommandOptions()

	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/common.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/scc.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/operator-openshift.yaml")
	//co.Commands = append(co.Commands, "pods")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/cluster.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/toolbox.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/object-openshift.yaml")
	co.Commands = append(co.Commands, "pods")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/route.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	//Switch Context and Reload Config Flags
	co.SwitchContext("rook-ceph")

	log.Println("Setup Ceph Object Storage with Rook Operator")
	for commandNumber, command := range co.Commands {

		/*//After the system pods are provisioned wait for them to become ready before moving on
		if commandNumber == 3 {
			log.Info("Waiting for Pods to be ready in rook-ceph-system namespace:")
			podStatus := utils.NewpodStatus()
			for podStatus.Running != 4 {
				cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
				cmd.Flags().Set("output", "yaml")
				cmd.Run(cmd, []string{command})
				podStatus.CountPods(out.Bytes())
				log.Debug(podStatus)
				log.Info("Waiting...")
				out.Reset()
				time.Sleep(5 * time.Second)
			}
		*/
		if commandNumber == 6 {
			//After the pods in rook-ceph are provisioned wait for them to become ready before moving on
			log.Print("Waiting for pods to be ready in rook-ceph")
			podStatus := utils.NewpodStatus()
			for podStatus.Running != 22 {
				cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
				cmd.Flags().Set("output", "yaml")
				cmd.Run(cmd, []string{command})
				podStatus.CountPods(out.Bytes())
				log.Debug(podStatus)
				log.Info("Waiting...")
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
			log.Info(out.String())
			out.Reset()
		}
	}

}

// setupCmd represents the setup command
var cephSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
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
