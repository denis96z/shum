package main

import (
	"log"

	"github.com/maildealru/shum/pkg/shum/service"
)

func main() {
	srv := service.NewService()
	if s, err := srv.Manage(); err == nil {
		log.Print(s)
	} else {
		log.Fatalf("%s\n%s", s, err.Error())
	}
}
