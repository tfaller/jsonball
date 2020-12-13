package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tfaller/jsonball/internal/operation"
	"github.com/tfaller/jsonball/internal/startup"
)

var (
	registry = startup.MustGetRegistry()

	// args
	docType = flag.String("docType", "", "-docType DOC_TYPE")
	docName = flag.String("docName", "", "-docName DOC_NAME")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	doc, err := operation.GetDocumentContent(ctx, registry, *docType, *docName)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(os.Stdout).Encode(doc)
	if err != nil {
		log.Fatal(err)
	}
}
