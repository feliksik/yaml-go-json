package main

import (
	"os"
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"

	"net/url"
	"net/http"

	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"github.com/imdario/mergo"


)

// parse YAML, assuming:
// - it can be converted to JSON (i.e. dicts use string-keys)
func parseYaml_jsonable(txt []byte) (map[string]interface{}, error) {
	structure := make(map[string]interface{})
	err := yaml.Unmarshal(txt, structure)
	if err!=nil {
		return nil, err
	}
	r, e := transformData(structure)
	return r.(map[string]interface{}), e
}


func getSourceData(loc string) ([]byte, error){
	pUrl, err := url.Parse(loc)
	if err!=nil{
		log.Fatal("Could not parse uri: "+loc)
	}
	switch pUrl.Scheme {
	case "":
		fallthrough
	case "file":
		return ioutil.ReadFile(loc)
	case "stdin":
		return ioutil.ReadAll(os.Stdin)
	case "http":
		fallthrough
	case "https":
		resp, err := http.Get(loc)
		if err!=nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	default:
		log.Fatalf("Unsupported source scheme: %s", pUrl.Scheme)
	}

	// could buffer this...
	return []byte{}, nil
}

func parseSource(loc string) map[string]interface{} {

	data, err := getSourceData(loc)
	if err!=nil {
		log.Fatal("Cannot read file: "+loc)
	}

	var structure map[string]interface{}
	structure, err = parseYaml_jsonable(data)
	if err!=nil {
		log.Fatalf("Cannot parse data from source %s: %s", loc, err)
	}
	return structure
}

func runApp(c *cli.Context) {
	// get ansible_ssh_host

	sources := c.StringSlice("source")

	mergedData := make(map[string]interface{})
	
	if len(sources)==0 { // use stdin if no source is specified
    	sources = []string{"stdin://"}
	}
	
	for _,source := range sources {
		s := parseSource(source)
		mergo.Merge(&mergedData, s)
	}
	// apply static variables
	var result []byte
	var err error

	switch c.String("format"){
	case "yaml":
		result, err = yaml.Marshal(mergedData)
	case "json-pretty":
		result, err = json.MarshalIndent(mergedData, "", "  ")
	default: // we parsed already, just print json even if something unknown was asked.
		result, err = json.Marshal(mergedData)
	}
	if err != nil{
		log.Fatal("Could not serialize structure")
	}
	fmt.Printf("%s\n", string(result));
	// generate and add derived variables
}

func createCommandlineApp() *cli.App {

	app := cli.NewApp()
	app.Name = "yaml-to-json"
	app.Usage = `Collect yaml data from sources, and merge them into json to stdout.

   The sources are in merged in order; where fields are redundant, the fields in the first
   sources take precedence.

   Lists are *not* merged. `
	app.Version = "0.1"

	app.Flags = []cli.Flag {
		cli.StringSliceFlag{
			Name: "s,source",
			Usage: "A list of sources (can be http://|file|stdin://|). No source means stdin is used.",
			//Value: "blah",
		},
		cli.StringFlag {
			Name: "f,format",
			Value: "json",
			Usage: "Output format: [json|json-pretty|yaml]",
		},
	}

	app.Action = runApp
	return app
}

func main() {
	app := createCommandlineApp()
	app.Run(os.Args)
}

