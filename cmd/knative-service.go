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
)

var (
	status = false
	logView = false
)

func service(service string) {
	
	cmd := exec.Command("./oc","apply", "-n", "knative-eventing","-f","yamls/" + service + "Service.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func serviceStatus() { 
	
	cmd := exec.Command("./oc","get","ksvc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func logs(name string) {
	
	podName, err := exec.Command("./oc" ,"get", "pods", "--selector='serving.knative.dev/service'").Output()
	if err != nil {
		log.Fatal(err)
	}	
	
	cmd := exec.Command("./oc", "logs", string(podName), "-c" ,"user-container", "--since=10m")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// serviceCmd represents the service command
var knative_serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service called")
		
		if (status) {
			serviceStatus()
		}else if(logView){
			logs(args[0])
		}else{
			service(args[0])
		}
	},
}

func init() {
	knativeCmd.AddCommand(knative_serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	knative_serviceCmd.Flags().BoolVarP(&status, "status", "S", false, "Show Status of the Service")
	knative_serviceCmd.Flags().BoolVarP(&logView, "logView", "l", false, "Show logs of the Service")

}
