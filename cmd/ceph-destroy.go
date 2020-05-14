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

//Made from Instructions @https://opendatahub.io/docs/administration/advanced-installation/object-storage.html for installing
//ceph object storage via the rook operator
func cephDestroy() {
	//Make command options for Knative Setup
	co := utils.NewCommandOptions()

	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/route.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/object-user.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/object-openshift.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/toolbox.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/cluster-on-pvc.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/operator-openshift.yaml")
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/common.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	//Switch Context and Reload Config Flags
	co.SwitchContext("rook-ceph-system")

	log.Println("Removing Ceph Via Rook v1.3.2")
	for _, command := range co.Commands {
		cmd := delete.NewCmdDelete(co.CurrentFactory, IOStreams)
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
	log.Println("For removal to be complete the user must manually remove rook related files from the worker nodes in the following directorys:")
	log.Println("1./var/lib/rook")
	log.Println("2./var/lib/kubelet/plugins")
	log.Println("3./var/lib/kubelet/plugins-registry")
}

// destroyCmd represents the destroy command
var cephDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the Ceph cluster",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("destroy called")
		cephDestroy()
	},
}

func init() {
	cephCmd.AddCommand(cephDestroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
