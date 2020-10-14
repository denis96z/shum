package service

import (
	"fmt"
	"os"

	"github.com/takama/daemon"
)

const (
	srvName  = "shum"
	srvDescr = "service that allows running shell command via HTTP API"
)

type Service struct {
	daemon.Daemon
}

func NewService() *Service {
	d, err := daemon.New(srvName, srvDescr, daemon.SystemDaemon)
	if err != nil {
		panic(err)
	}
	return &Service{
		Daemon: d,
	}
}

func (srv *Service) Manage() (string, error) {
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return srv.Install()
		case "remove":
			return srv.Remove()
		case "start":
			return srv.Start()
		case "stop":
			return srv.Stop()
		case "status":
			return srv.Status()
		}
	}

	usage := fmt.Sprintf(
		"Usage: %s install | remove | start | stop | status", srvName,
	)

	return usage, nil
}
