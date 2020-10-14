package api

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func (srv *API) HandleCmd(ctx *gin.Context) {
	name := ctx.Param("name")

	command, ok := srv.config.Shell.Commands[name]
	if !ok {
		WriteFailWithBody(ctx, makeCMDUnknownBody(name))
		return
	}

	args := make([]string, len(srv.config.Shell.Args)+1)
	copy(args, srv.config.Shell.Args)
	args[len(srv.config.Shell.Args)] = command.Command

	cmd := exec.Command(srv.config.Shell.Bin, args...)

	if command.Async {
		go func() {
			if err := cmd.Run(); err != nil {
				/*WRITE LOG*/
			}
		}()

		WriteOKWithBody(ctx, makeCMDAsyncOKBody())
		return
	}

	var bOut, bErr bytes.Buffer
	cmd.Stdout = &bOut
	cmd.Stderr = &bErr

	if err := cmd.Run(); err != nil {
		WriteErrWithBody(ctx, err)
		return
	}

	if command.RevealOutput {
		WriteOKWithBody(
			ctx, cmdOKBody{
				stdout: bOut,
				stderr: bErr,
			},
		)
	} else {
		WriteOK(ctx)
	}
}

type cmdOKBody struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (v cmdOKBody) MarshalJSON() ([]byte, error) {
	b := []byte(fmt.Sprintf(
		`{"status":"ok","stdout":%q,"stderr":%q}`,
		v.stdout.String(), v.stderr.String(),
	))
	return b, nil
}

type cmdAsyncOKBody struct{}

func makeCMDAsyncOKBody() cmdAsyncOKBody {
	return cmdAsyncOKBody{}
}

func (v cmdAsyncOKBody) MarshalJSON() ([]byte, error) {
	return []byte(`{"status":"ok","async":true}`), nil
}

type cmdUnknownBody struct {
	name string
}

func makeCMDUnknownBody(name string) cmdUnknownBody {
	return cmdUnknownBody{name: name}
}

func (v cmdUnknownBody) MarshalJSON() ([]byte, error) {
	return []byte(`{"status":"fail","message":"unknown command"}`), nil
}
