package utils

import (
	"log"
	"os"
	"strings"

	//"time"

	//"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"

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

func (co *CommandOptions) switchContext(nameSpace string) {
	//Load Config for Kubectl Wrapper Function
	co.configFlags = genericclioptions.NewConfigFlags(true)
	co.userSpecifiedNamespace = "kafka"
	//Create a new Credential factory from the kubeconfig file
	f := kcmdutil.NewFactory(co.configFlags)
	co.rawConfig, _ = f.ToRawKubeConfigLoader().RawConfig()

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

}
