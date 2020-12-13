package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tfaller/jsonball/event"
	"github.com/tfaller/jsonball/internal/operation"
	"github.com/tfaller/jsonball/internal/startup"
)

var (
	detector = startup.MustGetDetector()
)

func main() {
	flag.Parse()
	ctx := context.Background()

	var listen event.ListenOnChange
	err := json.NewDecoder(os.Stdin).Decode(&listen)
	if err != nil {
		log.Fatal(err)
	}

	err = operation.Listen(ctx, detector, listen)
	if err != nil {
		log.Fatal(err)
	}
}
