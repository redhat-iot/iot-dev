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
	"github.com/IoTCLI/cmd/utils"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/exec"
)

func login() {

	//download Enmasse v0.30.3
	_ = utils.DownloadAndUncompress("oc.gz", "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/openshift-client-linux.tar.gz")
	log.Println("oc Source folder: ", "oc")
	//Login to OC
	fmt.Print("Enter Openshift Username: ")
	var user string
	fmt.Scanln(&user)

	fmt.Print("Enter Openshift Password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	cmd := exec.Command("bash", "-c", "echo "+user+" "+string(password)+"| ./oc login")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Remove("oc")

}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("login called")
		login()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
