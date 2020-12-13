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
	detector = startup.MustGetDetector()
	registry = startup.MustGetRegistry()

	// args
	commit = flag.Bool("commit", false, "commits the change")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	if *commit {
		log.Println("Warning: the change will be removed without notifying a handler")
	}

	change, err := detector.NextChange(ctx)
	if err != nil {
		log.Fatal(err)
	}

	changeEvent, err := operation.HandleChange(ctx, registry, change)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewEncoder(os.Stdout).Encode(changeEvent); err != nil {
		log.Fatal(err)
	}

	if *commit {
		change.Commit()
	} else {
		change.Close()
	}
}
