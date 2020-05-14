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
	"io/ioutil"
	"os"

	"github.com/IoTCLI/cmd/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
)

func cephUser(user string) {
	co := utils.NewCommandOptions()

	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/object-user.yaml")

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	//Switch Context and Reload Config Flags
	co.SwitchContext("rook-ceph")

	tmpFile, err := ioutil.TempFile(os.TempDir(), "service-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(tmpFile.Name())

	sedCommands := []string{`s/name: .*/name: ` + user + `/`, `s/displayName: .*/displayName: "` + user + `"/`}

	myOutput := utils.RemoteSed(sedCommands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/ceph/setup/object-user.yaml")

	//Write updated yaml to tempfile
	tmpFile.Write(myOutput.Bytes())
	log.Println("the Source file: ", myOutput.String())
	//Close Tempfile after writing
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
	co.Commands = append(co.Commands, tmpFile.Name())

	log.Println("Provision Ceph User")
	for _, command := range co.Commands {
		cmd := apply.NewCmdApply("kubectl", co.CurrentFactory, IOStreams)
		//Kubectl signals missing field, set validate to false to ignore this
		//cmd.Flags().Set("validate", "false")
		err := cmd.Flags().Set("filename", command)
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run(cmd, []string{})
		log.Print(out.String())
		out.Reset()
	}

}

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user called")
		cephUser(args[0])
	},
}

func init() {
	cephCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
