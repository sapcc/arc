//go:generate go run main.go -o ../../descriptions.go
//go:generate goimports -w ../../descriptions.go

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	header = `package main

import (
)

`
)

var (
	flagOut = flag.String("o", "", "Output file, else stdout.")
	cmdDir  = "../../website/source/docs/commands/"
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

	generateUsage(w)

	fmt.Fprint(w, "\n")

	generateDescription(w)
}

func generateUsage(w io.Writer) {
	fmt.Fprint(w, `var cmdUsage = map[string]string{`)
	fmt.Fprint(w, "\n")

	mapGenerator(w,
		func(data string) string {
			ymlData := make(map[string]string)
			err := yaml.Unmarshal([]byte(data), &ymlData)
			if err != nil {
				log.Fatal(err)
			}
			parsedData := strings.Replace(ymlData["description"], "`", `"`, -1)
			return parsedData
		})

	fmt.Fprint(w, "}\n")
}

func generateDescription(w io.Writer) {
	fmt.Fprint(w, `var cmdDescription = map[string]string{`)
	fmt.Fprint(w, "\n")

	mapGenerator(w,
		func(data string) string {
			description := strExtract(data, "## Description", "##", 1)
			examples := strExtract(data, "## Examples", "##", 1)

			parsedDesc := strings.Replace(description, "`", `"`, -1)
			parsedDesc = strings.TrimSpace(parsedDesc)
			parsedExam := strings.Replace(examples, "`", `"`, -1)
			parsedExam = strings.TrimSpace(parsedExam)

			// add a break if examples
			if parsedExam != "" {
				parsedExam = fmt.Sprint("\n\n", parsedExam)
			}

			return fmt.Sprint(parsedDesc, parsedExam)
		})

	fmt.Fprint(w, "}\n")
}

func mapGenerator(w io.Writer, callback func(data string) string) {
	// check if dir exits
	d, err := os.Open(filepath.Clean(cmdDir))
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
			// data from file
			data, err := ioutil.ReadFile(filepath.Clean(path.Join(cmdDir, fi.Name())))
			if err != nil {
				log.Fatal(err)
			}

			ymlData := make(map[string]string)
			err = yaml.Unmarshal([]byte(data), &ymlData)
			if err != nil {
				log.Fatal(err)
			}

			item := callback(string(data))

			fmt.Fprint(w, "  \"", ymlData["sidebar_current"], "\": `", item, "`,")
			fmt.Fprint(w, "\n")
		}
	}
}

func strExtract(sExper, sAdelim, sCdelim string, nOccur int) string {
	aExper := strings.Split(sExper, sAdelim)

	if len(aExper) <= nOccur {
		return ""
	}

	sMember := aExper[nOccur]
	aExper = strings.Split(sMember, sCdelim)

	if len(aExper) == 1 {
		return sMember
	}

	return aExper[0]
}
