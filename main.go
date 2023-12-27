// main.go

package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/michaelahli/octopus/src/services"
	"github.com/michaelahli/octopus/svcutils/mainpkg"
	"github.com/spf13/viper"
)

func main() {
	cfg, err := mainpkg.ServiceConfig("env/config")
	if err != nil {
		log.Fatalf("initialize config: %s\n", err.Error())
	}

	serve(cfg)
}

func serve(cfg *viper.Viper) {
	port := cfg.GetString("server.port")

	svc := services.New(cfg)

	http.HandleFunc("/book", svc.HandleBooks)
	http.HandleFunc("/", svc.CommonHandler)

	log.Printf("Server is listening on port %s...\n", port)
	port = strings.Join([]string{"", port}, ":")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting the server: %s\n", err.Error())
	}
}
