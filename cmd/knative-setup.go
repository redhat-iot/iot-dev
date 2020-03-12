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
	"log"
	"os/exec"
	"github.com/spf13/cobra"
	"os"
	"time"

)

var status = false

func knativeServing() {
	
	ocCommands := [][]string{}

	ocCommands = append(ocCommands,[]string{"./oc","apply","-f","yamls/operatorgroup.yaml"} )
	ocCommands = append(ocCommands,[]string{"./oc","apply","-f","yamls/sub.yaml"} )
	ocCommands = append(ocCommands,[]string{"./oc","apply","-f","yamls/knative-serving.yaml"} )
	
	for command := range ocCommands {
		cmd := exec.Command(ocCommands[command][0], ocCommands[command][1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		
		//After installing the operator give cluster the time to distribute to all namespaces
		log.Println(command)
		if command == 1 {
			time.Sleep(60.0 * time.Second)
		}
	}
}

func knativeEventing() {

	ocCommands := [][]string{}

	ocCommands = append(ocCommands,[]string{"./oc","apply","--selector","knative.dev/crd-install=true","--filename", "https://github.com/knative/eventing/releases/download/v0.13.0/eventing.yaml"} )
	ocCommands = append(ocCommands,[]string{"./oc","apply", "--filename","https://github.com/knative/eventing/releases/download/v0.13.0/eventing.yaml"} )
	
	for command := range ocCommands {
		cmd := exec.Command(ocCommands[command][0], ocCommands[command][1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func knativeStatus(){
	ocCommands := [][]string{}

	ocCommands = append(ocCommands,[]string{"./oc","project"} )
	ocCommands = append(ocCommands,[]string{"./oc", "get", "knativeserving/knative-serving" ,"-n", "knative-serving", "--template='{{range .status.conditions}}{{printf '%s=%s' .type .status}}{{end}}'"} )
	ocCommands = append(ocCommands,[]string{"./oc","get", "pods","--namespace","knative-eventing"} )
	
	for command := range ocCommands {
		cmd := exec.Command(ocCommands[command][0], ocCommands[command][1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// setupCmd represents the setup command
var knative_setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if(status){
			log.Println("Checking on knative Eventing and Serving download status")
			knativeStatus()
		}else{
			fmt.Println("Installing Knative Serving")
			knativeServing()
			fmt.Println("Installing Knative Eventing")
			knativeEventing()
			log.Println("Checking on knative Eventing and Serving download status")
			knativeStatus()
		}
		
	},
}

func init() {
	knativeCmd.AddCommand(knative_setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	knative_setupCmd.Flags().BoolVarP(&status, "status", "S", false, "Check on status of knative install")
}
