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
	"github.com/spf13/cobra"
	"log"

	//in package import
	"github.com/IoTCLI/cmd/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	//"k8s.io/kubectl/pkg/cmd/"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/get"
)

var (
	status                      = false
	logView                     = false
	knativeServiceNamespaceFlag string
)

func service(service string) {

	//Make command options for Knative Setup
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving
	co.Commands = append(co.Commands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/knative/services/"+service+".yaml")
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(knativeServiceNamespaceFlag)

	log.Println("Provision Knative Service: ", service)
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

func serviceStatus() {

	//Make command options for Knative Setup
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(knativeServiceNamespaceFlag)

	//Reload config flags after switching contex

	log.Println("Get Knative Service Status")

	cmd := get.NewCmdGet("kubectl", co.CurrentFactory, IOStreams)
	cmd.Run(cmd, []string{"ksvc"})
	log.Print(out.String())
	out.Reset()
}

/*
func logs(name string) {


	//Make command options for Knative Setup
	co := utils.NewCommandOptions()

	//Install Openshift Serveless and  Knative Serving

	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()
	co.Commands = append(co.Commands, "pods")

	co.SwitchContext(knativeServiceNamespaceFlag)

	//Reload config flags after switching context
	newconfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionConfig := kcmdutil.NewMatchVersionFlags(newconfigFlags)
	cf := kcmdutil.NewFactory(matchVersionConfig)


	log.Println("Get Knative Service Status")

	cmd := get.NewCmdGet("kubectl", cf, IOStreams)
	cmd.Flags().Set("selector", "serving.knative.dev/service")
	cmd.Run(cmd, []string{})

	log.Print(out.String())
	out.Reset()


	podName, err := exec.Command("./oc", "get", "pods", "--selector='serving.knative.dev/service'").Output()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("./oc", "logs", string(podName), "-c", "user-container", "--since=10m")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
*/

// serviceCmd represents the service command
var knativeServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Knative Service called")

		if status {
			serviceStatus()
		} else if logView {
			//logs(args[0])
		} else {
			service(args[0])
		}
	},
}

func init() {
	knativeCmd.AddCommand(knativeServiceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	knativeServiceCmd.Flags().BoolVarP(&status, "status", "S", false, "Show Status of the Service")
	knativeServiceCmd.Flags().BoolVarP(&logView, "logView", "l", false, "Show logs of the Service")
	knativeServiceCmd.Flags().StringVarP(&knativeServiceNamespaceFlag, "namespace", "n", "knative-eventing", "Option to specify namespace for knative service deployment, defaults to 'knative-eventing'")
}
