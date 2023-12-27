package main

import (
	"log"

	svcpostgres "github.com/michaelahli/octopus/svcutils/storage/postgres"

	_ "github.com/lib/pq"
)

func main() {
	if err := svcpostgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
