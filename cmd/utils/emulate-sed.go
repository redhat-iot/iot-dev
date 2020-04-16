package utils

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/rwtodd/Go.Sed/sed"
)

func emulateSed(command string, url string) *bytes.Buffer {

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
