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
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

//CommandOptions ...
//Setup options for command
//Eventually move this to its own package
type CommandOptions struct {
	configFlags *genericclioptions.ConfigFlags

	newContext             *api.Context
	newContextName         string
	rawConfig              api.Config
	commands               []string
	userSpecifiedNamespace string

	genericclioptions.IOStreams
}

func newCommandOptions() *CommandOptions {
	return &CommandOptions{}
}

func isContextEqual(ctxA, ctxB *api.Context) bool {
	if ctxA == nil || ctxB == nil {
		return false
	}
	if ctxA.Cluster != ctxB.Cluster {
		return false
	}
	if ctxA.Namespace != ctxB.Namespace {
		return false
	}
	if ctxA.AuthInfo != ctxB.AuthInfo {
		return false
	}

	return true
}

//func switchNamespaceContext()

func kafkaSetup() {

	co := newCommandOptions()

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

	co.commands = ocCommands

	//Load Config for Kubectl Wrapper Function
	co.configFlags = genericclioptions.NewConfigFlags(true)
	co.userSpecifiedNamespace = "kafka"
	//Create a new Credential factory from the kubeconfig file
	f := kcmdutil.NewFactory(co.configFlags)
	co.rawConfig, err = f.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		log.Fatal(err)
	}
	co.IOStreams = genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stdout}

	currentContext, exists := co.rawConfig.Contexts[co.rawConfig.CurrentContext]
	if !exists {
		log.Fatal("Error no Context's avaliable")
	}
	co.newContext = api.NewContext()

	co.newContext.Cluster = currentContext.Cluster
	co.newContext.AuthInfo = currentContext.AuthInfo
	co.newContext.Namespace = co.userSpecifiedNamespace
	contextName := co.userSpecifiedNamespace + "/" + currentContext.Cluster + "/" + strings.Split(currentContext.AuthInfo, "/")[0]
	co.newContextName = contextName

	configAccess := clientcmd.NewDefaultPathOptions()

	if existingContext, exists := co.rawConfig.Contexts[co.newContextName]; !exists || !isContextEqual(co.newContext, existingContext) {
		co.rawConfig.Contexts[co.newContextName] = co.newContext
	}

	co.rawConfig.CurrentContext = co.newContextName
	clientcmd.ModifyConfig(configAccess, co.rawConfig, true)
	//update current factory
	//f.ToRawKubeConfigLoader().RawConfig() = co.rawConfig
	log.Println("Context switched to: ", co.userSpecifiedNamespace)
	//Make a new kubctl command
	//cmd := apply.NewCmdApply("kubectl", f, ioStreams)

	//reload configs after they have been altered
	newconfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionConfig := kcmdutil.NewMatchVersionFlags(newconfigFlags)
	cf := kcmdutil.NewFactory(matchVersionConfig)

	fmt.Println("running")
	for _, command := range ocCommands {
		cmd := apply.NewCmdApply("kubectl", cf, co.IOStreams)

		//commandgroup := templates.CommandGroup{}

		//commandgroup.Commands[0] = cmd

		//templates.UseOptionsTemplates(cmd)
		//cmd.Flags.add("namespace")
		//templates.ActsAsRootCommand(cmd, []string{}, commandgroup)
		//fmt.Println(templater.optionsCmdFor)

		//cmd.Flags().Set("output", "json")
		//cmd.Flags().Set("dry-run", "true")
		cmd.Flags().Set("filename", command)
		/*
			err = cmd.Flags().Set("n", "kafka")
			if err != nil {
				log.Fatal(err)
			}
		*/
		//.fmt.Println(cmd.Flags().Lookup("filename").Value.String())
		//fmt.Println(cmd.Flags().Lookup("namespace").Value.String())
		cmd.Run(cmd, []string{})
		//Allow Resources to stabilize
		time.Sleep(5 * time.Second)
	}

	os.RemoveAll("tmp/")
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
