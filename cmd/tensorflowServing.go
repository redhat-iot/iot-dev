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

	"github.com/spf13/cobra"
)

var (
	tensorflowServingNamespaceFlag string
)

// tensorflowServingCmd represents the tensorflowServing command
var tensorflowServingCmd = &cobra.Command{
	Use:   "tensorflowServing",
	Short: "Setup a tensorflow service for analytics",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tensorflowServing called")
	},
}

func init() {
	rootCmd.AddCommand(tensorflowServingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tensorflowServingCmd.PersistentFlags().String("foo", "", "A help for foo")
	tensorflowServingCmd.PersistentFlags().StringVarP(&tensorflowServingNamespaceFlag, "namespace", "n", "default", "Option to specify namespace for Tensorflow Serving Deployment, defaults to 'default'")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tensorflowServingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
