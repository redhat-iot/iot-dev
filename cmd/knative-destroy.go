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
	"os"
	"os/exec"
	"github.com/spf13/cobra"
)

func destroyKnative(){ 
	ocCommands := [][]string{}
	ocCommands = append(ocCommands,[]string{"./oc" , "delete", "--filename", "https://github.com/knative/eventing/releases/download/v0.13.0/eventing.yaml"})
	ocCommands = append(ocCommands,[]string{"./oc" ,"delete" ,"--selector", "knative.dev/crd-install=true" ,"--filename", "https://github.com/knative/eventing/releases/download/v0.13.0/eventing.yaml"})
	ocCommands = append(ocCommands,[]string{"./oc" , "delete", "namespace", "knative-eventing"} )
	ocCommands = append(ocCommands,[]string{"./oc", "delete", "knativeserving.operator.knative.dev", "knative-serving", "-n" ,"knative-serving"} )
	ocCommands = append(ocCommands,[]string{"./oc" , "delete", "namespace", "knative-serving"} )
	//ocCommands = append(ocCommands,[]string{"./oc","delete","-f","yamls/sub.yaml"} )
	//ocCommands = append(ocCommands,[]string{"./oc","delete","-f","yamls/operatorgroup.yaml"} )
	
	for command := range ocCommands {
		cmd := exec.Command(ocCommands[command][0], ocCommands[command][1:]...)
		cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			//Igonore Resource Not found error and continue, but still notify the user
			log.Println(err)
		}
	}

	//delete Operator

	currentCSV,err := exec.Command("bash","-c","./oc get subscription serverless-operator -n openshift-operators -o jsonpath='{.status.currentCSV}'").Output()
	err = exec.Command("./oc" ,"delete" ,"subscription", "serverless-operator" ,"-n" ,"openshift-operators").Run()
	if err != nil {
		//Igonore Resource Not found error and continue, but still notify the user
		log.Println(err)
	}
	err = exec.Command("./oc" ,"delete" ,"clusterserviceversion",string(currentCSV), "-n" ,"openshift-operators").Run()
	if err != nil {
		//Igonore Resource Not found error and continue, but still notify the user
		log.Println(err)
	}

}

// destroyCmd represents the destroy command
var knative_destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy called")
		destroyKnative()
	},
}

func init() {
	knativeCmd.AddCommand(knative_destroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
