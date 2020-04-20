package utils

import (
	"bytes"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rwtodd/Go.Sed/sed"
	"path/filepath"
)

//RemoteSed ...
//Use this function to emulate the bash Sed command in golang
func RemoteSed(command string, url string) *bytes.Buffer {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	engine, err := sed.New(strings.NewReader(command))
	myOutput := new(bytes.Buffer)
	myOutput.ReadFrom(engine.Wrap(resp.Body))
	return myOutput

}

func LocalSed(command string, dir string) {

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			file, err := os.Open(path)
			if err != nil {
				log.Panic(err)
			}

			engine, err := sed.New(strings.NewReader(command))

			if err != nil {
				log.Panic(err)
			}

			myOutput := new(bytes.Buffer)
			myOutput.ReadFrom(engine.Wrap(file))
			_, err = file.Write(myOutput.Bytes())

		}
		return nil

	})
	if err != nil {
		log.Panic(err)
	}
}
