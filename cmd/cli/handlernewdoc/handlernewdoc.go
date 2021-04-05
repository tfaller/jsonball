package main

import (
	"context"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/operation"
	"github.com/tfaller/jsonball/internal/startup"
)

var (
	registry = startup.MustGetRegistry()
	detector = startup.MustGetDetector()

	// args
	docType  = flag.String("docType", "", "docType which should be watched for new documents")
	handler  = flag.String("handler", "", "handler which should trigger")
	existing = flag.Bool("existing", false, "whether existing docs should be treated as new")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	err := operation.HandlerNewDoc(ctx, registry, detector, &event.HandlerNewDoc{
		Handler:  *handler,
		Type:     *docType,
		Existing: *existing,
	})
	if err != nil {
		log.Fatal(err)
	}
}
