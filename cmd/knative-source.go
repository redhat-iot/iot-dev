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
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/IoTCLI/cmd/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
)

var (
	knativeSourceNamespaceFlag string
)

func iotContainerSource(containerSource string) {
	ocCommands := [][]string{}

	messageURI, err := exec.Command("./oc", "-n", "myapp", "get", "addressspace", "iot", "-o", "jsonpath={.status.endpointStatuses[?(@.name=='messaging')].externalHost}").Output()
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("MESSAGE_URI", string(messageURI))
	log.Println(string(messageURI))

	os.Setenv("MESSAGE_PORT", "443")

	os.Setenv("MESSAGE_TYPE", "telemetry")

	os.Setenv("MESSAGE_TENANT", "myapp.iot")

	tlsCert, err := exec.Command("bash", "-c", "oc -n myapp get addressspace iot -o jsonpath={.status.caCert} | base64 --decode").Output()
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("TLS_CERT", string(tlsCert))

	os.Setenv("CLIENT_USERNAME", "consumer")

	os.Setenv("CLIENT_PASSWORD", "foobar")

	//ocCommands = append(ocCommands,[]string{"/bin/bash", "-c", ". ./scripts/iotVideoCS-SetupScript.sh"} )
	ocCommands = append(ocCommands, []string{"/bin/bash", "-c", "cat yamls/" + containerSource + "ContainerSource.yaml.in | envsubst | oc apply -n knative-eventing -f -"})

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

func source(source string, sink string) {
	//Make command options for Knative Setup
	co := utils.NewCommandOptions()

	//Custom source configs for various Knative Sources
	switch source {
	case "kafka":
		co.Commands = append(co.Commands, "https://storage.googleapis.com/knative-releases/eventing-contrib/latest/kafka-source.yaml")

	}

	//Make A Temporary file set Sink value in Source Yaml
	tmpFile, err := ioutil.TempFile(os.TempDir(), "source-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(tmpFile.Name())

	sedCommands := []string{`/^ *sink:/,/^ name:/s/name: .*/name: ` + sink + `/`}

	myOutput := utils.RemoteSed(sedCommands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/knative/sources/kafka.yaml")

	tmpFile.Write(myOutput.Bytes())
	log.Println("the Source file: ", myOutput.String())
	//Close Tempfile after writing
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	co.Commands = append(co.Commands, tmpFile.Name())
	IOStreams, _, out, _ := genericclioptions.NewTestIOStreams()

	co.SwitchContext(knativeSourceNamespaceFlag)

	log.Println("Provision Knative Source")
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

// knativeSourceCmd represents the cs command
var knativeSourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Deploy a knative source",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Knative Source called")
		if args[0] == "iotContainer" {
			iotContainerSource(args[0])
		} else {
			source(args[0], args[1])
		}
	},
}

func init() {
	knativeCmd.AddCommand(knativeSourceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// knativeSourceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// knativeSourceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	knativeSourceCmd.Flags().StringVarP(&knativeSourceNamespaceFlag, "namespace", "n", "knative-eventing", "Option to specify namespace for knative service deployment, defaults to 'knative-eventing'")
}
