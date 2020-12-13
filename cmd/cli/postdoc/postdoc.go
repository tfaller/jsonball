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
	detector = startup.MustGetDetector()

	// args
	quietDirectAccess = flag.Bool("quietDirectAccess", false, "to disable direct access warning")
	docType           = flag.String("docType", "", "-docType DOC_TYPE")
	docName           = flag.String("docName", "", "-docName DOC_NAME")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	if !*quietDirectAccess {
		log.Println("Warning: A document should not be directly inserted or changed")
	}

	var doc interface{}
	err := json.NewDecoder(os.Stdin).Decode(&doc)
	if err != nil {
		log.Fatal(err)
	}

	err = operation.PutDocument(ctx, registry, detector, *docType, *docName, doc)
	if err != nil {
		log.Fatal(err)
	}
}
