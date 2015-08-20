//go:generate go run main.go -o ../../descriptions.go

package main

import (
	"fmt"
	"path"
	"os"
	"io/ioutil"
	"flag"
	
	log "github.com/Sirupsen/logrus"	
	"gopkg.in/yaml.v2"
)

const (
	header = `package main
	
import (
)

var cmdDescription = map[string]string{
`
)

var (
	flagOut    = flag.String("o", "", "Output file, else stdout.")
	cmdDir     = "../../website/source/docs/commands/"
)

func main() {
	flag.Parse()
	var err error
	w := os.Stdout
	
	if *flagOut != "" {
		if w, err = os.Create(*flagOut); err != nil {
			log.Fatal(err)
		}
		defer w.Close()
	}
	
	fmt.Fprint(w, header)
	
	// check if dir exits
  d, err := os.Open(cmdDir)
  if err != nil {
      log.Fatal(err)
  }
  defer d.Close()
	
	// returns at most n FileInfo structures
  fi, err := d.Readdir(-1)
  if err != nil {
      log.Fatal(err)
  }
	
	// loop over the docs
  for _, fi := range fi {
      if fi.Mode().IsRegular() {
					markdownFile := path.Join(cmdDir, fi.Name())
					
					// data from file
					data, err := ioutil.ReadFile(markdownFile)
					if err != nil {
						log.Fatal(err)
					}
					
					ymlData := make(map[string]string)
	
					err = yaml.Unmarshal([]byte(data), &ymlData)
					if err != nil {
						log.Fatal(err)
					}
					
					fmt.Fprint(w, "  \"", ymlData["sidebar_current"], "\": \"", ymlData["description"], "\",")
					fmt.Fprint(w, "\n")
      }
  }
	
	fmt.Fprint(w, "}\n")
}