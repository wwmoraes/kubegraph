package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	apiExtensionsApiServerScheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/client-go/kubernetes/scheme"
	aggregatorScheme "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/scheme"

	"github.com/goccy/go-graphviz"
	kubegraph "github.com/wwmoraes/kubegraph/internal/kubegraph"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln(errors.New("usage: kubegraph <file>"))
	}

	_ = aggregatorScheme.AddToScheme(scheme.Scheme)
	_ = apiExtensionsApiServerScheme.AddToScheme(scheme.Scheme)
	decode := scheme.Codecs.UniversalDeserializer().Decode

	log.Println("reading file...")
	fileBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
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
			log.Fatal(err)
		}
	}
	fileString = cleanFileString.String()
	// ioutil.WriteFile("clean.yaml", []byte(fileString), 0644)

	log.Println("splitting documents...")
	documents := strings.Split(fileString, "---")

	log.Println("initializing kubegraph instance...")
	kubegraphInstance, err := kubegraph.New()
	if err != nil {
		log.Fatal(err)
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

		_, err = kubegraphInstance.Transform(obj)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("connecting nodes...")
	kubegraphInstance.ConnectNodes()

	log.Println("generating graph...")
	if err := kubegraphInstance.Render("graph.dot", graphviz.XDOT); err != nil {
		log.Fatal(err)
	}
	if err := kubegraphInstance.Render("graph.svg", graphviz.SVG); err != nil {
		log.Fatal(err)
	}
}
