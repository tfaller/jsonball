package main

import (
	"context"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tfaller/jsonball/internal/startup"
)

var (
	registry = startup.MustGetRegistry()

	// args
	docType = flag.String("type", "", "the document type to register")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	err := registry.RegisterDocumentType(ctx, *docType)
	if err != nil {
		log.Fatal(err)
	}
}
