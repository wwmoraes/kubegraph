package loader

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/wwmoraes/kubegraph/internal/kubegraph"
)

// FromYAML creates a KubeGraph instance from a YAML file contents
func FromYAML(fileName string) (*kubegraph.KubeGraph, error) {
	decode := getDecoder()

	log.Println("reading file...")
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// normalize line breaks
	log.Println("normalizing linebreaks...")
	fileString := string(fileBytes[:])
	fileString = strings.ReplaceAll(fileString, "\r\n", "\n")
	fileString = strings.ReplaceAll(fileString, "\r", "\n")

	// removes all comments from yaml and json
	log.Println("removing comments and empty lines...")
	commentLineMatcher, err := regexp.Compile("^[ ]*((#|//).*)?$")
	if err != nil {
		return nil, err
	}
	fileStringLines := strings.Split(fileString, "\n")
	var cleanFileString strings.Builder
	for _, line := range fileStringLines {
		if commentLineMatcher.MatchString(line) {
			continue
		}
		if line == "\n" || line == "" {
			continue
		}

		_, err := cleanFileString.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			return nil, err
		}
	}
	fileString = cleanFileString.String()

	log.Println("splitting documents...")
	documents := strings.Split(fileString, "---")

	log.Println("initializing kubegraph instance...")
	instance, err := kubegraph.New()
	if err != nil {
		return nil, err
	}

	for _, document := range documents {
		if document == "\n" || document == "" {
			continue
		}

		obj, _, err := decode([]byte(document), nil, nil)
		if err != nil {
			log.Printf("unable to decode document: %++v\n", err)
			continue
		}

		_, err = instance.Transform(obj)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("connecting nodes...")
	instance.ConnectNodes()

	return instance, nil
}
