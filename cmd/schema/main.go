package main

import (
	"log"

	"github.com/xeipuuv/gojsonschema"
)

func main() {
	schemaLoader := gojsonschema.NewReferenceLoader("https://raw.githubusercontent.com/mavryk-network/mavryk-snapshot-metadata-schema/main/mavryk-snapshot-metadata.schema.json")
	documentLoader := gojsonschema.NewReferenceLoader("http://localhost:8080/mavryk-snapshots.json")

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		log.Printf("The document is valid\n")
	} else {
		log.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			log.Printf("- %s\n", desc)
		}
	}
}
