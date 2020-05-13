package utils

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
	"os"
	"strings"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

//CommandOptions ...
//Setup options for command
//Eventually move this to its own package
type CommandOptions struct {
	configFlags            *kcmdutil.MatchVersionFlags
	CurrentFactory         kcmdutil.Factory
	newContext             *api.Context
	newContextName         string
	clientConfig           *rest.Config
	rawConfig              api.Config
	Commands               []string
	userSpecifiedNamespace string

	IOStreams genericclioptions.IOStreams
}

//NewCommandOptions ...
//Use this to make a new CommandOptions Struct
func NewCommandOptions() *CommandOptions {

	newconfigFlags := genericclioptions.NewConfigFlags(true)
	matchVersionConfig := kcmdutil.NewMatchVersionFlags(newconfigFlags)
	return &CommandOptions{
		configFlags: matchVersionConfig,

		CurrentFactory: kcmdutil.NewFactory(matchVersionConfig),
	}
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

//SwitchContext ...
//Use this function to Switch kubeconfig Contexts
func (co *CommandOptions) SwitchContext(nameSpace string) {
	//Load Config for Kubectl Wrapper Function
	co.userSpecifiedNamespace = nameSpace
	//Create a new Credential factory from the kubeconfig file
	co.rawConfig, _ = co.CurrentFactory.ToRawKubeConfigLoader().RawConfig()

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
	newconfigFlags := genericclioptions.NewConfigFlags(true)
	co.configFlags = kcmdutil.NewMatchVersionFlags(newconfigFlags)
	co.CurrentFactory = kcmdutil.NewFactory(co.configFlags)

	log.Info("Context switched to namespace:", co.userSpecifiedNamespace)

}

//GetUserToken ...
//Returns Token for current user
func (co *CommandOptions) GetUserToken() string {
	newconfigFlags := genericclioptions.NewConfigFlags(true)
	co.configFlags = kcmdutil.NewMatchVersionFlags(newconfigFlags)
	f := kcmdutil.NewFactory(co.configFlags)
	co.clientConfig, _ = f.ToRESTConfig()
	return (co.clientConfig.BearerToken)
}
