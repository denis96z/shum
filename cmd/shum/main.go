package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/maildealru/shum/pkg/shum/api"
	"github.com/maildealru/shum/pkg/shum/conf"
	"github.com/maildealru/shum/pkg/shum/errs"
	"github.com/maildealru/shum/pkg/shum/service"
)

func main() {
	c := conf.NewConfig()
	if err := c.TryLoad(os.Args); err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	srv := service.NewService()
	if ok, s, err := srv.Manage(os.Args); err == nil {
		if ok {
			log.Print(s)
			return
		}
	} else {
		log.Fatalf("%s\n%s", s, err.Error())
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	apiSrv := api.NewAPI(c)
	go func() {
		if err := apiSrv.Start(); err != nil {
			if errs.Cause(err) != http.ErrServerClosed {
				log.Fatalf("Failed to start API service\n%s", err.Error())
			}
		}
	}()

	<-ch

	if err := apiSrv.Stop(); err != nil {
		log.Fatalf("Failed to stop API service\n%s", err.Error())
	}
}
