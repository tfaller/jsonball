package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tfaller/jsonball/internal/startup"
)

var (
	registry = startup.MustGetRegistry()

	// args
	docType    = flag.String("docType", "", "-docType DOC_TYPE")
	startToken = flag.String("startToken", "", "-startToken TOKEN")
	maxDocs    = flag.Uint("maxDocs", 100, "-docName 100")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	if *maxDocs > 65535 {
		log.Fatalf("%v is too large as maxDocs", maxDocs)
	}

	docList, err := registry.ListDocuments(ctx, *docType, *startToken, uint16(*maxDocs))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("NextToken: %v\n", docList.NextToken)

	for _, doc := range docList.Documents {
		fmt.Println(doc)
	}
}
