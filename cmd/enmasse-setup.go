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
	"os"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"archive/tar"
	"compress/gzip"
	"path/filepath"
	"log"
	"os/exec"
	
)

var (
	ocUrl = "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/openshift-client-linux-4.3.5.tar.gz"
	enmasseUrl = "https://github.com/EnMasseProject/enmasse/releases/download/0.30.2/enmasse-0.30.2.tgz"
)

func downloadPkg(filepath string, url string) error{
	resp, err := http.Get(url)
	if err != nil{
		return err
	}

	defer resp.Body.Close()

	out,err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_,err = io.Copy(out, resp.Body)
	return err
}

func Untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func setup() { 
	downloadPkg("oc.tar.gz", ocUrl)
	downloadPkg("enmasse.tgz",enmasseUrl)
	
	ocContent, err := os.Open("oc.tar.gz")
	if err != nil {
		log.Fatal(err)
	}

	enmasseContent,err := os.Open("enmasse.tgz")
	if err != nil {
		log.Fatal(err)
	}

	path, err := os.Getwd()
	if err != nil {
    	log.Fatal(err)
	}
	
	Untar(path,ocContent)
	Untar(path,enmasseContent)

	cmd := exec.Command("./oc","project")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

}



// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called")
		setup()
	},
}

func init() {
	enmasseCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
