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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/rwtodd/Go.Sed/sed"
)

func kafkaSetup() {

	ocCommands := []string{}

	//This section is mimicking the instructions to setup the Strimzi Operator, I.E download the install yaml, and set namespace
	os.Mkdir("tmp/", 0755)

	resp, err := http.Get("https://github.com/strimzi/strimzi-kafka-operator/releases/download/0.17.0/strimzi-cluster-operator-0.17.0.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	engine, err := sed.New(strings.NewReader(`s/namespace: .*/namespace: kafka/`))
	myOutput := new(bytes.Buffer)

	myOutput.ReadFrom(engine.Wrap(resp.Body))

	ioutil.WriteFile("tmp/strim.yaml", myOutput.Bytes(), 0755)

	//End of Strimzi install section

	//List of commands to install strimzi and then provision kafka
	ocCommands = append(ocCommands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/kafka-namespace.yaml")
	ocCommands = append(ocCommands, "tmp/strim.yaml")
	//ocCommands = append(ocCommands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/kafka-rolebindings.yaml")
	ocCommands = append(ocCommands, "https://raw.githubusercontent.com/redhat-iot/iot-dev/master/yamls/kafka.yaml")

	//Load Config for Kubectl Wrapper Function
	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionKubeConfigFlags := kcmdutil.NewMatchVersionFlags(kubeConfigFlags)

	//Create a new Credential factory
	f := kcmdutil.NewFactory(matchVersionKubeConfigFlags)

	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stdout}

	//Make a new kubctl command
	//cmd := apply.NewCmdApply("kubectl", f, ioStreams)
	fmt.Println("running")
	for _, command := range ocCommands {

		cmd := apply.NewCmdApply("kubectl", f, ioStreams)

		//cmd.Flags().Set("output", "json")
		//cmd.Flags().Set("dry-run", "true")
		cmd.Flags().Set("filename", command)
		cmd.Flags().Set("namespace", "kafka")
		cmd.Run(cmd, []string{})
		//Allow Resources to stabilize
		time.Sleep(10 * time.Second)
	}
}

// setupCmd represents the setup command
var kafkaSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Kafka with Strimzi Operator on a single Openshift namespace",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kafka setup called")
		kafkaSetup()

	},
}

func init() {
	kafkaCmd.AddCommand(kafkaSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
