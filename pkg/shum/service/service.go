package service

import (
	"fmt"
	"strings"

	"github.com/maildealru/shum/pkg/shum/consts"
	"github.com/maildealru/shum/pkg/shum/errs"

	"github.com/takama/daemon"
)

type Service struct {
	daemon.Daemon
}

func NewService() *Service {
	d, err := daemon.New(
		consts.Name, consts.Description, daemon.SystemDaemon,
	)
	if err != nil {
		panic(err)
	}
	return &Service{
		Daemon: d,
	}
}

func (srv *Service) Manage(args []string) (bool, string, error) {
	if len(args) > 1 {
		if len(args) == 2 {
			command := args[1]
			switch command {
			case "install":
				s, err := srv.Install()
				return true, s, err
			case "remove":
				s, err := srv.Remove()
				return true, s, err
			case "start":
				s, err := srv.Start()
				return true, s, err
			case "stop":
				s, err := srv.Stop()
				return true, s, err
			case "status":
				s, err := srv.Status()
				return true, s, err
			}
			//NOTE: only management commands will have args[1]
			//      that does not have flag/option prefix
			if !strings.HasPrefix(command, "-") {
				usage := fmt.Sprintf(
					"Usage: %s install | remove | start | stop | status [OPTIONS]", consts.Name,
				)
				return true, usage, errs.Errorf("unknown command %q", command)
			}
		}
	}
	return false, "", nil
}
