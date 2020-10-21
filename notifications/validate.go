package main

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Missing filename argument. Usage: validate file.json\n")
		return
	}

	fmt.Print(Validate(os.Args[1]))
}

func Validate(document string) string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + cwd + "/batchNotification.json")
	documentLoader := gojsonschema.NewReferenceLoader("file://" + cwd + "/" + document)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		return "The document is valid\n"
	} else {
		var strBuilder = strings.Builder{}

		strBuilder.WriteString("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			strBuilder.WriteString(fmt.Sprintf("- %s\n", desc))
		}
		return strBuilder.String()
	}
}
